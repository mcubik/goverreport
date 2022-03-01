package report

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mattn/go-zglob"
	"github.com/mcubik/goverreport/config"
	"golang.org/x/tools/cover"
)

// Coverage summary for a file or module
type Summary struct {
	Name                                       string
	Blocks, Stmts, MissingBlocks, MissingStmts int
	BlockCoverage, StmtCoverage                float64
}

// Report of the coverage results
type Report struct {
	Total Summary   // Global coverage
	Files []Summary // Coverage by file
}

// Generates a coverage report given the coverage profile file, and the following configurations:
// exclusions: packages to be excluded (if a package is excluded, all its subpackages are excluded as well)
// sortBy: the order in which the files will be sorted in the report (see sortResults)
// order: the direction of the the sorting
func GenerateReport(coverprofile string, configuration config.Configuration, sortBy, order string, packages bool) (Report, error) {
	profiles, err := cover.ParseProfiles(coverprofile)
	if err != nil {
		return Report{}, fmt.Errorf("Invalid coverprofile: '%s'", err)
	}
	total := &accumulator{name: "Total"}
	files := make(map[string]*accumulator)
	for _, profile := range profiles {
		fileName := normalizeName(profile.FileName, configuration.Root, packages)
		if isExcluded(fileName, configuration.Exclusions) {
			continue
		}
		fileCover, ok := files[fileName]
		if !ok {
			// Create new accumulator
			fileCover = &accumulator{name: fileName}
			files[fileName] = fileCover
		}
		total.addAll(profile.Blocks)
		fileCover.addAll(profile.Blocks)
	}
	return makeReport(total, files, sortBy, order)
}

// Removes root dir part if configured to do so
func normalizeName(fileName string, root string, packages bool) string {
	if packages {
		fileName = filepath.Dir(fileName)
	}

	if root == "" {
		return fileName
	}
	if packages {
		return "." + strings.Replace(fileName, root, "", -1)
	}
	return strings.Replace(fileName, root, "", -1)
}

func isExcluded(fileName string, exclusions []string) bool {
	for _, exclusion := range exclusions {
		if ok, _ := zglob.Match(exclusion, fileName); ok {
			return true
		}
	}
	return false
}

// Creates a Report struct from the coverage sumarization results
func makeReport(total *accumulator, files map[string]*accumulator, sortBy, order string) (Report, error) {
	fileReports := make([]Summary, 0, len(files))
	for _, fileCover := range files {
		fileReports = append(fileReports, fileCover.results())
	}
	if err := sortResults(fileReports, sortBy, order); err != nil {
		return Report{}, err
	}
	return Report{
		Total: total.results(),
		Files: fileReports}, nil
}

// Accumulates the coverage of a file and returns a summary
type accumulator struct {
	name                                       string
	blocks, stmts, coveredBlocks, coveredStmts int
}

// Accumulates a profile block
func (a *accumulator) add(block cover.ProfileBlock) {
	a.blocks++
	a.stmts += block.NumStmt
	if block.Count > 0 {
		a.coveredBlocks++
		a.coveredStmts += block.NumStmt
	}
}

func (a *accumulator) addAll(blocks []cover.ProfileBlock) {
	for _, block := range blocks {
		a.add(block)
	}
}

// Creates a summary with the accumulated values
func (a *accumulator) results() Summary {
	return Summary{
		Name:          a.name,
		Blocks:        a.blocks,
		Stmts:         a.stmts,
		MissingBlocks: a.blocks - a.coveredBlocks,
		MissingStmts:  a.stmts - a.coveredStmts,
		BlockCoverage: float64(a.coveredBlocks) / float64(a.blocks) * 100,
		StmtCoverage:  float64(a.coveredStmts) / float64(a.stmts) * 100}
}

// Sorts the individual coverage reports by a given column
// (block --block coverage--, stmt --stmt coverage--, missing-blocks or missing-stmts)
// and a sorting direction (asc or desc)
func sortResults(reports []Summary, mode string, order string) error {
	var reverse bool
	var cmp func(i, j int) bool
	switch order {
	case "asc":
		reverse = false
	case "desc":
		reverse = true
	default:
		return errors.New("Order must be either asc or desc")
	}
	switch mode {
	case "filename", "package":
		cmp = func(i, j int) bool {
			return reports[i].Name < reports[j].Name
		}
	case "block":
		cmp = func(i, j int) bool {
			return reports[i].BlockCoverage < reports[j].BlockCoverage
		}
	case "stmt":
		cmp = func(i, j int) bool {
			return reports[i].StmtCoverage < reports[j].StmtCoverage
		}
	case "missing-blocks":
		cmp = func(i, j int) bool {
			return reports[i].MissingBlocks < reports[j].MissingBlocks
		}
	case "missing-stmts":
		cmp = func(i, j int) bool {
			return reports[i].MissingStmts < reports[j].MissingStmts
		}
	default:
		return errors.New("Invalid sort colum, must be one of filename, package, block, stmt, missing-blocks or missing-stmts")
	}
	sort.Slice(reports, func(i, j int) bool {
		if reverse {
			return !cmp(i, j)
		} else {
			return cmp(i, j)
		}
	})
	return nil
}
