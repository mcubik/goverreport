package main

import (
	"bytes"
	"github.com/mcubik/goverreport/report"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	assert := assert.New(t)
	conf, err := loadConfig(".goverreport.yml")
	assert.NoError(err)
	assert.Equal(conf, configuration{
		Root:       "github.com/mcubik/goverreport",
		Exclusions: []string{"test", "vendor"},
		Threshold:  80,
		Metric:     "stmt"})
}

func TestEmptyConfig(t *testing.T) {
	assert := assert.New(t)
	conf, err := loadConfig("emptyconfig.yml")
	assert.NoError(err)
	assert.Equal(conf, configuration{
		Root:       "",
		Exclusions: []string{},
		Threshold:  0,
		Metric:     ""})
}

func TestEmptyConfigWhenFileMissing(t *testing.T) {
	assert := assert.New(t)
	conf, err := loadConfig("xxxxxx.yml")
	assert.NoError(err)
	assert.Equal(conf, configuration{
		Root:       "",
		Exclusions: []string{},
		Threshold:  0,
		Metric:     ""})
}

func TestThreshold(t *testing.T) {
	assert := assert.New(t)
	summary := report.Summary{BlockCoverage: 79.9, StmtCoverage: 82.3}
	res, err := checkThreshold(80, summary, "block")
	assert.NoError(err)
	assert.False(res)

	res, err = checkThreshold(79, summary, "block")
	assert.NoError(err)
	assert.True(res)

	res, err = checkThreshold(79.9, summary, "block")
	assert.NoError(err)
	assert.True(res)

	res, err = checkThreshold(82, summary, "stmt")
	assert.NoError(err)
	assert.True(res)

}

func TestNoThreshold(t *testing.T) {
	assert := assert.New(t)
	summary := report.Summary{BlockCoverage: 79.9, StmtCoverage: 82.3}
	res, err := checkThreshold(0, summary, "block")
	assert.NoError(err)
	assert.True(res)
}

func TestInvalidTresholdType(t *testing.T) {
	assert := assert.New(t)
	summary := report.Summary{BlockCoverage: 79.9, StmtCoverage: 82.3}
	_, err := checkThreshold(80, summary, "xxxx")
	assert.Error(err)
}

func TestRun(t *testing.T) {
	assert := assert.New(t)
	args := arguments{
		coverprofile: "sample_coverage.out",
		threshold:    80,
		metric:       "block",
		sortBy:       "filename",
		order:        "asc"}
	buf := bytes.Buffer{}
	passed, err := run(configuration{}, args, &buf)
	assert.NoError(err)
	assert.False(passed)
	assert.Contains(buf.String(), "TOTAL", "Table generated")
}

func TestRunAboveThreshold(t *testing.T) {
	assert := assert.New(t)
	args := arguments{
		coverprofile: "sample_coverage.out",
		threshold:    75,
		metric:       "block",
		sortBy:       "filename",
		order:        "asc"}
	buf := bytes.Buffer{}
	passed, err := run(configuration{}, args, &buf)
	assert.NoError(err)
	assert.True(passed)
}

func TestRunFailInvalidArugment(t *testing.T) {
	assert := assert.New(t)
	_, err := run(configuration{}, arguments{
		coverprofile: "sample_coverage.out",
		threshold:    80,
		metric:       "xxx",
		sortBy:       "filename",
		order:        "asc"},
		new(bytes.Buffer))
	assert.Error(err)
}

func TestTakesConfigurationIfNotOverriden(t *testing.T) {
	assert := assert.New(t)
	config := configuration{Threshold: 80, Metric: "block"}
	args := arguments{
		coverprofile:    "sample_coverage.out",
		threshold:       0,
		metric:          "",
		sortBy:          "filename",
		order:           "asc",
		metricDefaulted: true}
	buf := bytes.Buffer{}
	passed, err := run(config, args, &buf)
	assert.NoError(err)
	assert.False(passed)
}
