package report

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// PrintTable prints the report to the terminal
func PrintTable(r Report, w io.Writer, packages bool) {
	// Create table with ASCII border style for compatibility with tests
	table := tablewriter.NewTable(w, 
		tablewriter.WithSymbols(tw.NewSymbols(tw.StyleASCII)),
		tablewriter.WithHeaderAutoFormat(tw.Off), // Disable auto-formatting to preserve case
	)

	// Determine header label based on packages flag
	item := "File"
	if packages {
		item = "Package"
	}

	// Set headers to match all columns from makeRow
	table.Header(item, "Blocks", "Missing", "Stmts", "Missing", "Block cover %", "Stmt cover %")

	// Add rows for all files
	for _, s := range r.Files {
		table.Append(makeRow(s))
	}

	// Add footer with totals
	table.Footer(makeRow(r.Total))

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
