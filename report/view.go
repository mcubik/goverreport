package report

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
)

// Prints the report to the terminal
func PrintTable(report Report, writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{
		"File", "Blocks", "Missing", "Stmts", "Missing",
		"Block cover %", "Stmt cover %"})
	for _, fileCoverage := range report.Files {
		table.Append(makeRow(fileCoverage))
	}
    table.SetAutoFormatHeaders(false)
	table.SetFooter(makeRow(report.Total))
	table.Render()
}

// Converts a Summary to a slice of string so that it
// can be printed in the table
func makeRow(c Summary) []string {
	return []string{
		c.Name,
		fmt.Sprintf("%d", c.Blocks),
		fmt.Sprintf("%d", c.MissingBlocks),
		fmt.Sprintf("%d", c.Stmts),
		fmt.Sprintf("%d", c.MissingStmts),
		fmt.Sprintf("%.2f", c.BlockCoverage),
		fmt.Sprintf("%.2f", c.StmtCoverage)}
}
