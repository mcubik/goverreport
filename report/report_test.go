package report

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReport(t *testing.T) {
	assert := assert.New(t)
	report, err := GenerateReport("../sample_coverage.out", "", []string{}, "block", "desc")
	assert.NoError(err)
	assert.InDelta(75.6, report.Total.BlockCoverage, 0.1)
	assert.InDelta(80, report.Total.StmtCoverage, 0.1)
	assert.Equal(1801, report.Total.Stmts)
	assert.Equal(1097, report.Total.Blocks)
}
