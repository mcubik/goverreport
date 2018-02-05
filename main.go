package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/mcubik/goverreport/report"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

type arguments struct {
	coverprofile, metric, sortBy, order string
	threshold                           float64
}

var args arguments

const configFile = "goverreport.yml"

type configuration struct {
	Root       string   `yaml:"root"`
	Exclusions []string `yaml:"exclusions"`
	Threshold  float64  `yaml:"threshold,omitempty"`
	Metric     string   `yaml:"thresholdType,omitempty"`
}

func init() {
	flag.StringVar(&args.coverprofile, "coverprofile", "coverage.out", "Coverage output file")
	flag.Float64Var(&args.threshold, "threshold", 0, "Return an error if the coverage is below a threshold")
	flag.StringVar(&args.metric, "metric", "block", "Use a specific metric for the threshold: block, line")
	flag.StringVar(&args.sortBy, "sort", "filename", "Column to sort by: filename, block, stmt, missing-blocks, missing-stmts")
	flag.StringVar(&args.order, "order", "asc", "Sort order: asc, desc")
}

func main() {

	// Parse arguments
	flag.Parse()
	config, err := loadConfig(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	passed, err := run(config, args, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	if !passed {
		os.Exit(1)
	}
}

// Runs the command
func run(config configuration, args arguments, writer io.Writer) (bool, error) {

	// Use config values if arguments aren't set
	if args.metric == "" {
		args.metric = config.Metric
	}
	if args.threshold == 0 {
		args.threshold = config.Threshold
	}

	report, err := report.GenerateReport(args.coverprofile, config.Root, config.Exclusions, args.sortBy, args.order)
	if err != nil {
		return false, err
	}
	printTable(report, writer)
	passed, err := checkThreshold(args.threshold, report.Total, args.metric)
	if err != nil {
		return false, err
	}
	return passed, nil
}

// Loads the report configuration from a yml file
func loadConfig(filename string) (configuration, error) {
	conf := configuration{
		Exclusions: []string{},
		Metric:     "block"}
	data, err := ioutil.ReadFile(filename)
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

// Checks whether the coverage is above a threshold value.
// thresholdType states which value will be used to check the threshold,
// block coverage (block) or statement coverage (stmt).
func checkThreshold(threshold float64, total report.Summary, metric string) (bool, error) {
	if threshold > 0 {
		switch metric {
		case "block":
			if total.BlockCoverage < threshold {
				return false, nil
			}
		case "stmt":
			if total.StmtCoverage < threshold {
				return false, nil
			}
		default:
			return false, errors.New("Invalid threshold type, use block or stmt")
		}
	}
	return true, nil
}

// Prints the report to the terminal
func printTable(report report.Report, writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{
		"File", "Blocks", "Missing", "Stmts", "Missing",
		"Block cover %", "Stmt cover %"})
	for _, fileCoverage := range report.Files {
		table.Append(makeRow(fileCoverage))
	}
	table.SetFooter(makeRow(report.Total))
	table.Render()
}

// Converts a Summary to a slice of string so that it
// can be printed in the table
func makeRow(c report.Summary) []string {
	return []string{
		c.Name,
		fmt.Sprintf("%d", c.Blocks),
		fmt.Sprintf("%d", c.MissingBlocks),
		fmt.Sprintf("%d", c.Stmts),
		fmt.Sprintf("%d", c.MissingStmts),
		fmt.Sprintf("%.2f", c.BlockCoverage),
		fmt.Sprintf("%.2f", c.StmtCoverage)}
}
