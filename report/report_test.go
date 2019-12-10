package report

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var results []Summary
var cover1 Summary
var cover2 Summary

func init() {
	cover1 = Summary{
		Name:          "file1",
		BlockCoverage: 0.5, StmtCoverage: 0.6,
		MissingBlocks: 10, MissingStmts: 8}

	cover2 = Summary{
		Name:          "file2",
		BlockCoverage: 0.6, StmtCoverage: 0.5,
		MissingBlocks: 8, MissingStmts: 10}

	results = []Summary{cover1, cover2}
}

func TestSortByFileName(t *testing.T) {
	assert.NoError(t, sortResults(results, "filename", "asc"))
	assert.Equal(t, results, []Summary{cover1, cover2})
}

func TestSortByBlockCoverage(t *testing.T) {
	assert.NoError(t, sortResults(results, "block", "desc"))
	assert.Equal(t, results, []Summary{cover2, cover1})

}

func TestSortByStmtCoverage(t *testing.T) {
	assert.NoError(t, sortResults(results, "stmt", "desc"))
	assert.Equal(t, results, []Summary{cover1, cover2})
}

func TestSortByMissingBlocks(t *testing.T) {
	assert.NoError(t, sortResults(results, "missing-blocks", "asc"))
	assert.Equal(t, results, []Summary{cover2, cover1})
}

func TestSortByMissingStmts(t *testing.T) {
	assert.NoError(t, sortResults(results, "missing-stmts", "asc"))
	assert.Equal(t, results, []Summary{cover1, cover2})
}

func TestInvalidParameters(t *testing.T) {
	assert.Error(t, sortResults(results, "xxx", "asc"))
	assert.Error(t, sortResults(results, "block", "yyy"))
}

func TestReport(t *testing.T) {
	assert := assert.New(t)
	report, err := GenerateReport("../sample_coverage.out", "", []string{}, "block", "desc", false)
	assert.NoError(err)
	assert.InDelta(81.4, report.Total.BlockCoverage, 0.1)
	assert.InDelta(81.9, report.Total.StmtCoverage, 0.1)
	assert.Equal(111, report.Total.Stmts)
	assert.Equal(81, report.Total.Blocks)
}

func TestInvalidCoverProfile(t *testing.T) {
	_, err := GenerateReport("../xxx.out", "", []string{}, "block", "desc", false)
	assert.Error(t, err)
}
