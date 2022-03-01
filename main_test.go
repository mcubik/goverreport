package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/mcubik/goverreport/config"
	"github.com/mcubik/goverreport/report"
	"github.com/mcubik/goverreport/testdata"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	assert := assert.New(t)
	conf, err := loadConfig(".goverreport.yml")
	assert.NoError(err)
	assert.Equal(conf, config.Configuration{
		Root:       "github.com/mcubik/goverreport",
		Exclusions: []string{"test", "vendor"},
		Threshold:  80,
		Metric:     "stmt"})
}

func TestEmptyConfig(t *testing.T) {
	assert := assert.New(t)
	conf, err := loadConfig(testdata.Filename("emptyconfig.yml"))
	assert.NoError(err)
	assert.Equal(conf, config.Configuration{
		Exclusions: []string{},
		Threshold:  0,
		Metric:     ""})
}

func TestEmptyConfigWhenFileMissing(t *testing.T) {
	assert := assert.New(t)
	conf, err := loadConfig("xxxxxx.yml")
	assert.NoError(err)
	assert.Equal(conf, config.Configuration{
		Exclusions: []string{},
		Threshold:  0,
		Metric:     ""})
}

func TestThreshold(t *testing.T) {
	assert := assert.New(t)
	summary := report.Summary{BlockCoverage: 79.9, StmtCoverage: 82.3}
	passed, err := checkThreshold(80, summary, "block")
	assert.NoError(err)
	assert.False(passed)

	passed, err = checkThreshold(79, summary, "block")
	assert.NoError(err)
	assert.True(passed)

	passed, err = checkThreshold(79.9, summary, "block")
	assert.NoError(err)
	assert.True(passed)

	passed, err = checkThreshold(82, summary, "stmt")
	assert.NoError(err)
	assert.True(passed)

	passed, err = checkThreshold(83, summary, "stmt")
	assert.NoError(err)
	assert.False(passed)
}

func TestNoThreshold(t *testing.T) {
	assert := assert.New(t)
	summary := report.Summary{BlockCoverage: 79.9, StmtCoverage: 82.3}
	passed, err := checkThreshold(0, summary, "block")
	assert.NoError(err)
	assert.True(passed)
}

func TestInvalidMetric(t *testing.T) {
	assert := assert.New(t)
	summary := report.Summary{BlockCoverage: 79.9, StmtCoverage: 82.3}
	_, err := checkThreshold(80, summary, "xxxx")
	assert.Error(err)
}

func TestRun(t *testing.T) {
	assert := assert.New(t)
	args := arguments{
		coverprofile: testdata.Filename("sample_coverage.out"),
		threshold:    82,
		metric:       "block",
		sortBy:       "filename",
		order:        "asc"}
	buf := bytes.Buffer{}
	passed, err := run(config.Configuration{}, args, &buf)
	assert.NoError(err)
	assert.False(passed)
	assert.Contains(buf.String(), "Total", "Table generated")
}

func TestRunAboveThreshold(t *testing.T) {
	assert := assert.New(t)
	args := arguments{
		coverprofile: testdata.Filename("sample_coverage.out"),
		threshold:    75,
		metric:       "block",
		sortBy:       "filename",
		order:        "asc"}
	buf := bytes.Buffer{}
	passed, err := run(config.Configuration{}, args, &buf)
	assert.NoError(err)
	assert.True(passed)
}

func TestRunFailInvalidArugment(t *testing.T) {
	assert := assert.New(t)
	_, err := run(config.Configuration{}, arguments{
		coverprofile: testdata.Filename("sample_coverage.out"),
		threshold:    80,
		metric:       "xxx",
		sortBy:       "filename",
		order:        "asc"},
		new(bytes.Buffer))
	assert.Error(err)
}

func TestTakesConfigurationIfNotOverridenByCommandLineArgs(t *testing.T) {
	assert := assert.New(t)
	config := config.Configuration{Threshold: 80, Metric: "stmt"}
	args := arguments{
		coverprofile:    testdata.Filename("sample_coverage.out"),
		threshold:       0,
		metric:          "block",
		sortBy:          "filename",
		order:           "asc",
		metricDefaulted: true}
	buf := bytes.Buffer{}
	passed, err := run(config, args, &buf)
	assert.NoError(err)
	assert.True(passed) // Passes stmt coverage
}

func TestCommandLineArgsOverridesConfiguration(t *testing.T) {
	assert := assert.New(t)
	config := config.Configuration{Threshold: 80, Metric: "block"}
	args := arguments{
		coverprofile: testdata.Filename("sample_coverage.out"),
		threshold:    0,
		metric:       "stmt",
		sortBy:       "filename",
		order:        "asc"}
	buf := bytes.Buffer{}
	passed, err := run(config, args, &buf)
	assert.NoError(err)
	assert.True(passed) // Passes stmt coverage
}

func TestMetricArgument(t *testing.T) {
	os.Args[1] = ""
	assert := assert.New(t)
	// Metric default value
	parseArguments()
	assert.Equal("block", args.metric)
	assert.True(args.metricDefaulted)

	// Metric set
	os.Args[1] = "-metric=stmt"
	parseArguments()
	assert.Equal("stmt", args.metric)
	assert.False(args.metricDefaulted)
}

func TestRunPackages(t *testing.T) {
	assert := assert.New(t)
	args := arguments{
		coverprofile: testdata.Filename("sample_coverage.out"),
		packages:     true,
		sortBy:       "package",
		order:        "asc"}
	buf := bytes.Buffer{}
	passed, err := run(config.Configuration{}, args, &buf)
	assert.NoError(err)
	assert.True(passed)
	assert.Contains(buf.String(), "Package", "Column title is package")
	assert.Contains(buf.String(), "| github.com/mcubik/goverreport ", "Package .")
	assert.Contains(buf.String(), "| github.com/mcubik/goverreport/report |", "Package ./report")
}

func TestRunPackagesWithRoot(t *testing.T) {
	assert := assert.New(t)
	config := config.Configuration{
		Root: "github.com/mcubik/goverreport",
	}
	args := arguments{
		coverprofile: testdata.Filename("sample_coverage.out"),
		packages:     true,
		sortBy:       "package",
		order:        "asc"}
	buf := bytes.Buffer{}
	passed, err := run(config, args, &buf)
	assert.NoError(err)
	assert.True(passed)
	assert.Contains(buf.String(), "Package", "Column title is package")
	assert.Contains(buf.String(), "| . ", "Package .")
	assert.Contains(buf.String(), "| ./report |", "Package ./report")
}
