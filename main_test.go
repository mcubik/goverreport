package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReport(t *testing.T) {
	assert := assert.New(t)
	report := generateReport("sample_coverage.out")
	assert.Equal(79.7, report.BlockCoverage)
	assert.Equal(82.5, report.StmtCoverage)

}
