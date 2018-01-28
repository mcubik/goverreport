package main

import (
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/tools/cover"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var (
	coverprofile, thresholdType, sortBy, order string
	threshold                                  int
)

type config struct {
	Root       string   `yaml:"root"`
	Exclusions []string `yaml:"exclusions"`
}

func init() {
	flag.StringVar(&coverprofile, "coverprofile", "coverage.out", "Write a coverage profile to the file after all tests have passed")
	flag.IntVar(&threshold, "threshold", 0, "Return an error if the coverage is below a threshold")
	flag.StringVar(&thresholdType, "type", "block", "Use a specific metric for the threshold: block, line")
	flag.StringVar(&sortBy, "sort", "filename", "Sort: filename, block, stmt, missing-blocks, missing-stmts")
	flag.StringVar(&order, "order", "asc", "Sort order: asc, desc")
}

func main() {
	flag.Parse()
	fmt.Println(sortBy)
	config := config{}
	data, err := ioutil.ReadFile("coverreport.yml")
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Sprintf("Error reading configuration file: %s", err))
		}
		fmt.Println("No config file")
	} else {
		if err := yaml.Unmarshal(data, &config); err != nil {
			panic(fmt.Sprintf("Error loading settings: %v", err))
		}
		fmt.Println(config)
	}

	profiles, err := cover.ParseProfiles(coverprofile)
	if err != nil {
		fmt.Println("Invalid coverprofile")
		return
	}
	global := &coverage{name: "Total"}
	coverByFile := make(map[string]*coverage)
	for _, profile := range profiles {
		fileName := strings.Replace(profile.FileName, config.Root+"/", "", -1)
		skip := false
		for _, exclusion := range config.Exclusions {
			if strings.HasPrefix(fileName, exclusion) {
				skip = true
			}
		}
		if skip {
			continue
		}
		fileCover, ok := coverByFile[fileName]
		if !ok {
			fileCover = &coverage{name: fileName}
			coverByFile[fileName] = fileCover
		}
		for _, block := range profile.Blocks {
			global.addBlock(block)
			fileCover.addBlock(block)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"File", "Blocks", "Missing", "Stmts", "Missing",
		"Block cover %", "Stmt cover %"})
	keys := sortedKeys(coverByFile, sortBy, order)
	for _, key := range keys {
		table.Append(coverByFile[key].report())
	}
	table.SetFooter(global.report())
	table.Render()

	if threshold > 0 {
		if thresholdType == "block" {
			if global.blockCoverage() < threshold {
				os.Exit(1)
			}
		}
		if thresholdType == "line" {
			if global.stmtCoverage() < threshold {
				os.Exit(1)
			}
		}
	}
}

func sortedKeys(results map[string]*coverage, mode string, order string) []string {
	keys := make([]string, 0, len(results))
	for key, _ := range results {
		keys = append(keys, key)
	}
	applyOrder := func(b bool) bool { return b }
	if order == "desc" {
		applyOrder = func(b bool) bool { return !b }
	}
	if mode == "filename" {
		sort.Strings(keys)
	} else if mode == "block" {
		sort.Slice(keys, func(i, j int) bool {
			return applyOrder(results[keys[i]].blockCoverage() < results[keys[j]].blockCoverage())
		})
	} else if mode == "stmt" {
		sort.Slice(keys, func(i, j int) bool {
			return applyOrder(results[keys[i]].stmtCoverage() < results[keys[j]].stmtCoverage())
		})
	} else if mode == "missing-blocks" {
		sort.Slice(keys, func(i, j int) bool {
			return applyOrder(results[keys[i]].missingBlocks() < results[keys[j]].missingBlocks())
		})
	} else if mode == "missing-blocks" {
		sort.Slice(keys, func(i, j int) bool {
			return applyOrder(results[keys[i]].missingStmts() < results[keys[j]].missingStmts())
		})
	} else {
		panic(fmt.Sprintf("Sort by %s is not implemented", mode))
	}
	return keys
}

type coverage struct {
	name                                       string
	blocks, stmts, coveredBlocks, coveredStmts int
}

func (c *coverage) blockCoverage() int {
	return int(float64(c.coveredBlocks) / float64(c.blocks) * 100)
}

func (c *coverage) stmtCoverage() int {
	return int(float64(c.coveredStmts) / float64(c.stmts) * 100)
}

func (c *coverage) missingBlocks() int {
	return c.blocks - c.coveredBlocks
}

func (c *coverage) missingStmts() int {
	return c.stmts - c.coveredStmts
}

func (c *coverage) addBlock(block cover.ProfileBlock) {
	c.blocks++
	c.stmts += block.NumStmt
	if block.Count > 0 {
		c.coveredBlocks++
		c.coveredStmts += block.NumStmt
	}
}

func (c *coverage) report() []string {
	return []string{
		c.name,
		fmt.Sprintf("%d", c.blocks),
		fmt.Sprintf("%d", c.missingBlocks()),
		fmt.Sprintf("%d", c.stmts),
		fmt.Sprintf("%d", c.missingStmts()),
		fmt.Sprintf("%d", c.blockCoverage()),
		fmt.Sprintf("%d", c.stmtCoverage())}
}
