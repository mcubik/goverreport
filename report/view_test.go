package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintTable_FileMode(t *testing.T) {
	// Create a sample report with multiple files
	report := Report{
		Files: []Summary{
			{
				Name:          "/main.go",
				Blocks:        30,
				MissingBlocks: 10,
				Stmts:         44,
				MissingStmts:  15,
				BlockCoverage: 66.67,
				StmtCoverage:  65.91,
			},
			{
				Name:          "/report/report.go",
				Blocks:        47,
				MissingBlocks: 5,
				Stmts:         60,
				MissingStmts:  5,
				BlockCoverage: 89.36,
				StmtCoverage:  91.67,
			},
		},
		Total: Summary{
			Name:          "Total",
			Blocks:        77,
			MissingBlocks: 15,
			Stmts:         104,
			MissingStmts:  20,
			BlockCoverage: 80.52,
			StmtCoverage:  80.77,
		},
	}

	var buf bytes.Buffer
	require.NoError(t, PrintTable(report, &buf, false)) // false = file mode (default)

	output := buf.String()

	// Verify header contains "File" not "Package"
	assert.Contains(t, output, "File", "Header should contain 'File' in file mode")
	assert.NotContains(t, output, "Package", "Header should not contain 'Package' in file mode")

	// Verify column headers
	assert.Contains(t, output, "Blocks", "Should contain 'Blocks' column")
	assert.Contains(t, output, "Missing", "Should contain 'Missing' column")
	assert.Contains(t, output, "Stmts", "Should contain 'Stmts' column")
	assert.Contains(t, output, "Block cover %", "Should contain 'Block cover %' column")
	assert.Contains(t, output, "Stmt cover %", "Should contain 'Stmt cover %' column")

	// Verify individual files are shown
	assert.Contains(t, output, "/main.go", "Should display main.go file")
	assert.Contains(t, output, "/report/report.go", "Should display report.go file")

	// Verify file data is shown
	assert.Contains(t, output, "30", "Should show blocks count for main.go")
	assert.Contains(t, output, "66.67", "Should show block coverage for main.go")
	assert.Contains(t, output, "65.91", "Should show stmt coverage for main.go")

	// Verify total is shown
	assert.Contains(t, output, "Total", "Should display Total row")
	assert.Contains(t, output, "77", "Should show total blocks")
	assert.Contains(t, output, "80.52", "Should show total block coverage")
	assert.Contains(t, output, "80.77", "Should show total stmt coverage")

	// Verify it's a proper table with borders
	assert.Contains(t, output, "+---", "Should have ASCII table borders")
	assert.Contains(t, output, "|", "Should have table column separators")
}

func TestPrintTable_PackageMode(t *testing.T) {
	// Create a sample report with packages
	report := Report{
		Files: []Summary{
			{
				Name:          ".",
				Blocks:        30,
				MissingBlocks: 10,
				Stmts:         44,
				MissingStmts:  15,
				BlockCoverage: 66.67,
				StmtCoverage:  65.91,
			},
			{
				Name:          "./report",
				Blocks:        51,
				MissingBlocks: 5,
				Stmts:         67,
				MissingStmts:  5,
				BlockCoverage: 90.20,
				StmtCoverage:  92.54,
			},
		},
		Total: Summary{
			Name:          "Total",
			Blocks:        81,
			MissingBlocks: 15,
			Stmts:         111,
			MissingStmts:  20,
			BlockCoverage: 81.48,
			StmtCoverage:  81.98,
		},
	}

	var buf bytes.Buffer
	require.NoError(t, PrintTable(report, &buf, true)) // true = package mode

	output := buf.String()

	// Verify header contains "Package" not "File"
	assert.Contains(t, output, "Package", "Header should contain 'Package' in package mode")

	// Note: In package mode, the first occurrence of "File" should not be in the header
	lines := strings.Split(output, "\n")
	headerFound := false
	for _, line := range lines {
		if strings.Contains(line, "Package") && strings.Contains(line, "Blocks") {
			headerFound = true
			assert.NotContains(t, line, "| File ", "Header line should not contain '| File ' in package mode")
			break
		}
	}
	assert.True(t, headerFound, "Should find header line with Package")

	// Verify column headers
	assert.Contains(t, output, "Blocks", "Should contain 'Blocks' column")
	assert.Contains(t, output, "Missing", "Should contain 'Missing' column")
	assert.Contains(t, output, "Stmts", "Should contain 'Stmts' column")
	assert.Contains(t, output, "Block cover %", "Should contain 'Block cover %' column")
	assert.Contains(t, output, "Stmt cover %", "Should contain 'Stmt cover %' column")

	// Verify packages are shown
	assert.Contains(t, output, ".", "Should display root package")
	assert.Contains(t, output, "./report", "Should display report package")

	// Verify package data is shown
	assert.Contains(t, output, "30", "Should show blocks count for root package")
	assert.Contains(t, output, "51", "Should show blocks count for report package")
	assert.Contains(t, output, "66.67", "Should show block coverage for root package")
	assert.Contains(t, output, "90.20", "Should show block coverage for report package")

	// Verify total is shown
	assert.Contains(t, output, "Total", "Should display Total row")
	assert.Contains(t, output, "81", "Should show total blocks")
	assert.Contains(t, output, "81.48", "Should show total block coverage")
	assert.Contains(t, output, "81.98", "Should show total stmt coverage")

	// Verify it's a proper table with borders
	assert.Contains(t, output, "+---", "Should have ASCII table borders")
	assert.Contains(t, output, "|", "Should have table column separators")
}

func TestPrintTable_EmptyReport(t *testing.T) {
	// Test with an empty report
	report := Report{
		Files: []Summary{},
		Total: Summary{
			Name:          "Total",
			Blocks:        0,
			MissingBlocks: 0,
			Stmts:         0,
			MissingStmts:  0,
			BlockCoverage: 0.0,
			StmtCoverage:  0.0,
		},
	}

	var buf bytes.Buffer
	require.NoError(t, PrintTable(report, &buf, false))

	output := buf.String()

	// Should still have proper table structure
	assert.Contains(t, output, "File", "Should have File header even with no files")
	assert.Contains(t, output, "Total", "Should display Total row")
	assert.Contains(t, output, "|", "Should have table structure")
}

func TestPrintTable_SingleFile(t *testing.T) {
	// Test with a single file
	report := Report{
		Files: []Summary{
			{
				Name:          "/main.go",
				Blocks:        10,
				MissingBlocks: 2,
				Stmts:         15,
				MissingStmts:  3,
				BlockCoverage: 80.00,
				StmtCoverage:  80.00,
			},
		},
		Total: Summary{
			Name:          "Total",
			Blocks:        10,
			MissingBlocks: 2,
			Stmts:         15,
			MissingStmts:  3,
			BlockCoverage: 80.00,
			StmtCoverage:  80.00,
		},
	}

	var buf bytes.Buffer
	require.NoError(t, PrintTable(report, &buf, false))

	output := buf.String()

	// Verify single file is shown
	assert.Contains(t, output, "/main.go", "Should display the single file")
	assert.Contains(t, output, "10", "Should show blocks count")
	assert.Contains(t, output, "80.00", "Should show coverage")

	// Verify total matches (should be same as single file)
	totalCount := strings.Count(output, "80.00")
	assert.GreaterOrEqual(t, totalCount, 2, "Should show coverage values for both file and total")
}

func TestMakeRow(t *testing.T) {
	summary := Summary{
		Name:          "/test.go",
		Blocks:        50,
		MissingBlocks: 5,
		Stmts:         75,
		MissingStmts:  10,
		BlockCoverage: 90.00,
		StmtCoverage:  86.67,
	}

	row := makeRow(summary)

	// Verify row has correct number of columns
	assert.Equal(t, 7, len(row), "Row should have 7 columns")

	// Verify row contents
	assert.Equal(t, "/test.go", row[0], "First column should be name")
	assert.Equal(t, "50", row[1], "Second column should be blocks")
	assert.Equal(t, "5", row[2], "Third column should be missing blocks")
	assert.Equal(t, "75", row[3], "Fourth column should be stmts")
	assert.Equal(t, "10", row[4], "Fifth column should be missing stmts")
	assert.Equal(t, "90.00", row[5], "Sixth column should be block coverage")
	assert.Equal(t, "86.67", row[6], "Seventh column should be stmt coverage")
}
