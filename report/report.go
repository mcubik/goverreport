package report

import (
	"fmt"
	"golang.org/x/tools/cover"
	"sort"
	"strings"
)

type CoverInfo struct {
	Name                                       string
	Blocks, Stmts, MissingBlocks, MissingStmts int
	BlockCoverage, StmtCoverage                float64
}

type CoverReport struct {
	Total CoverInfo
	Files []CoverInfo
}

func GenerateReport(coverprofile string, root string, exclusions []string, sortBy, order string) (CoverReport, error) {
	profiles, err := cover.ParseProfiles(coverprofile)
	if err != nil {
		return CoverReport{}, fmt.Errorf("Invalid coverprofile: %s", err)
	}
	global := &coverAcc{name: "Total"}
	coverByFile := make(map[string]*coverAcc)
	for _, profile := range profiles {
		var fileName string
		if root == "" {
			fileName = profile.FileName
		} else {
			fileName = strings.Replace(profile.FileName, root+"/", "", -1)
		}
		skip := false
		for _, exclusion := range exclusions {
			if strings.HasPrefix(fileName, exclusion) {
				skip = true
			}
		}
		if skip {
			continue
		}
		fileCover, ok := coverByFile[fileName]
		if !ok {
			fileCover = &coverAcc{name: fileName}
			coverByFile[fileName] = fileCover
		}
		for _, block := range profile.Blocks {
			global.addBlock(block)
			fileCover.addBlock(block)
		}
	}
	return makeReport(global, coverByFile, sortBy, order), nil
}

func makeReport(global *coverAcc, coverByFile map[string]*coverAcc, sortBy, order string) CoverReport {
	fileReports := make([]CoverInfo, 0, len(coverByFile))
	for _, fileCover := range coverByFile {
		fileReports = append(fileReports, fileCover.report())
	}
	sortResults(fileReports, sortBy, order)
	return CoverReport{
		Total: global.report(),
		Files: fileReports}
}

type coverAcc struct {
	name                                       string
	blocks, stmts, coveredBlocks, coveredStmts int
}

func (c *coverAcc) report() CoverInfo {
	return CoverInfo{
		Name:          c.name,
		Blocks:        c.blocks,
		Stmts:         c.stmts,
		MissingBlocks: c.blocks - c.coveredBlocks,
		MissingStmts:  c.stmts - c.coveredStmts,
		BlockCoverage: float64(c.coveredBlocks) / float64(c.blocks) * 100,
		StmtCoverage:  float64(c.coveredStmts) / float64(c.stmts) * 100}
}

func (c *coverAcc) addBlock(block cover.ProfileBlock) {
	c.blocks++
	c.stmts += block.NumStmt
	if block.Count > 0 {
		c.coveredBlocks++
		c.coveredStmts += block.NumStmt
	}
}

func sortResults(reports []CoverInfo, mode string, order string) {
	applyOrder := func(b bool) bool { return b }
	if order == "desc" {
		applyOrder = func(b bool) bool { return !b }
	}
	if mode == "filename" {
		sort.Slice(reports, func(i, j int) bool {
			return applyOrder(reports[i].Name < reports[j].Name)
		})
	} else if mode == "block" {
		sort.Slice(reports, func(i, j int) bool {
			return applyOrder(reports[i].BlockCoverage < reports[j].BlockCoverage)
		})
	} else if mode == "stmt" {
		sort.Slice(reports, func(i, j int) bool {
			return applyOrder(reports[j].StmtCoverage < reports[j].StmtCoverage)
		})
	} else if mode == "missing-blocks" {
		sort.Slice(reports, func(i, j int) bool {
			return applyOrder(reports[i].MissingBlocks < reports[j].MissingBlocks)
		})
	} else if mode == "missing-blocks" {
		sort.Slice(reports, func(i, j int) bool {
			return applyOrder(reports[i].MissingStmts < reports[j].MissingStmts)
		})
	} else {
		panic(fmt.Sprintf("Sort by %s is not implemented", mode))
	}
}
