package main

import (
	"flag"
	"fmt"
	"github.com/mcubik/goverreport/report"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	coverprofile, thresholdType, sortBy, order string
	threshold                                  float64
)

type configuration struct {
	Root       string   `yaml:"root"`
	Exclusions []string `yaml:"exclusions"`
}

func init() {
	flag.StringVar(&coverprofile, "coverprofile", "coverage.out", "Write a coverage profile to the file after all tests have passed")
	flag.Float64Var(&threshold, "threshold", 0, "Return an error if the coverage is below a threshold")
	flag.StringVar(&thresholdType, "type", "block", "Use a specific metric for the threshold: block, line")
	flag.StringVar(&sortBy, "sort", "filename", "Sort: filename, block, stmt, missing-blocks, missing-stmts")
	flag.StringVar(&order, "order", "asc", "Sort order: asc, desc")
}

func main() {

	flag.Parse()
	config, err := loadConfig()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error loading configuration file: %s", err))
		os.Exit(2)
	}
	report, err := report.GenerateReport(coverprofile, config.Root, config.Exclusions, sortBy, order)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error generating the report: %s", err))
		os.Exit(2)
	}

	printTable(report)
	if !checkThreshold(threshold, report.Total, thresholdType) {
		os.Exit(1)
	}

	os.Exit(0)

}

func printTable(report report.CoverReport) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"File", "Blocks", "Missing", "Stmts", "Missing",
		"Block cover %", "Stmt cover %"})
	for _, fileCoverage := range report.Files {
		table.Append(makeRow(fileCoverage))
	}
	table.SetFooter(makeRow(report.Total))
	table.Render()
}

func checkThreshold(threshold float64, total report.CoverInfo, thresholdType string) bool {
	if threshold > 0 {
		if thresholdType == "block" {
			if total.BlockCoverage < threshold {
				return false
			}
		}
		if thresholdType == "line" {
			if total.StmtCoverage < threshold {
				return false
			}
		}
	}
	return true
}

func loadConfig() (configuration, error) {
	var conf configuration
	data, err := ioutil.ReadFile("coverreport.yml")
	if err != nil {
		if !os.IsNotExist(err) {
			return conf, err
		}
	} else {
		if err := yaml.Unmarshal(data, &conf); err != nil {
			return configuration{}, err
		}
	}
	return conf, nil
}

func makeRow(c report.CoverInfo) []string {
	return []string{
		c.Name,
		fmt.Sprintf("%d", c.Blocks),
		fmt.Sprintf("%d", c.MissingBlocks),
		fmt.Sprintf("%d", c.Stmts),
		fmt.Sprintf("%d", c.MissingStmts),
		fmt.Sprintf("%.2f", c.BlockCoverage),
		fmt.Sprintf("%.2f", c.StmtCoverage)}
}
