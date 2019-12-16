package report

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

// Prints the report to the terminal
func PrintTable(report Report, writer io.Writer, packages bool) {
	item := "File"
	if packages {
		item = "Package"
	}
	table := tablewriter.NewWriter(writer)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
	table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	table.SetHeader([]string{
		item, "Blocks", "Missing", "Stmts", "Missing",
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
