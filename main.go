package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/mcubik/goverreport/report"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

// Command arguments
type arguments struct {
	coverprofile, metric, sortBy, order string
	threshold                           float64
}

var args arguments

const configFile = "goverreport.yml"

// Configuration
type configuration struct {
	Root       string   `yaml:"root"`
	Exclusions []string `yaml:"exclusions"`
	Threshold  float64  `yaml:"threshold,omitempty"`
	Metric     string   `yaml:"thresholdType,omitempty"`
}

// Parser arguments
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

	rep, err := report.GenerateReport(args.coverprofile, config.Root, config.Exclusions, args.sortBy, args.order)
	if err != nil {
		return false, err
	}
	report.PrintTable(rep, writer)
	passed, err := checkThreshold(args.threshold, rep.Total, args.metric)
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
