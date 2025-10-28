package main

import (
	"strings"

	"golang.org/x/net/html"
)

type Table struct {
	Node      *html.Node
	HasHeader bool       // If true, Rows[0] contains the header row
	Rows      [][]string // All rows including header (if HasHeader is true)
}

// ExtractHtmlTable extracts a table structure from an HTML table node
// Handles three cases:
//   1. <table><thead>...</thead><tbody>...</tbody></table>
//   2. <table><tbody>...</tbody></table> (no thead)
//   3. <table><tr>...</tr>...</table> (direct tr children, no thead/tbody)
func ExtractHtmlTable(n *html.Node) Table {
	if n == nil || n.Data != "table" {
		return Table{Node: n, Rows: [][]string{}}
	}

	table := Table{
		Node:      n,
		HasHeader: false,
		Rows:      [][]string{},
	}

	var hasTheadOrTbody bool
	var theadRows [][]string
	var tbodyRows [][]string
	var directRows [][]string

	// First pass: check structure and collect rows
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode {
			continue
		}

		switch child.Data {
		case "thead":
			hasTheadOrTbody = true
			table.HasHeader = true
			theadRows = extractRowsFromSection(child)

		case "tbody":
			hasTheadOrTbody = true
			tbodyRows = extractRowsFromSection(child)

		case "tr":
			// Direct tr child (no thead/tbody wrapper)
			directRows = append(directRows, extractRowCells(child))
		}
	}

	// Assemble final rows based on structure
	if hasTheadOrTbody {
		// Case 1 or 2: has thead and/or tbody
		table.Rows = append(table.Rows, theadRows...)
		table.Rows = append(table.Rows, tbodyRows...)
	} else {
		// Case 3: direct tr children
		table.Rows = directRows

		// Heuristic: if first row has all th elements, treat it as header
		if len(directRows) > 0 && isHeaderRow(n.FirstChild) {
			table.HasHeader = true
		}
	}

	return table
}

// extractRowsFromSection extracts all rows from a thead or tbody element
func extractRowsFromSection(section *html.Node) [][]string {
	var rows [][]string
	for row := section.FirstChild; row != nil; row = row.NextSibling {
		if row.Type == html.ElementNode && row.Data == "tr" {
			cells := extractRowCells(row)
			if len(cells) > 0 {
				rows = append(rows, cells)
			}
		}
	}
	return rows
}

// extractRowCells extracts all cell values from a tr element
// Handles both th and td elements
func extractRowCells(row *html.Node) []string {
	var cells []string
	for cell := row.FirstChild; cell != nil; cell = cell.NextSibling {
		if cell.Type == html.ElementNode && (cell.Data == "th" || cell.Data == "td") {
			cellText := strings.TrimSpace(getTextContent(cell))
			cells = append(cells, cellText)
		}
	}
	return cells
}

// isHeaderRow checks if a tr element contains primarily th elements
func isHeaderRow(tr *html.Node) bool {
	// Find the first tr element
	for tr != nil && (tr.Type != html.ElementNode || tr.Data != "tr") {
		tr = tr.NextSibling
	}
	if tr == nil {
		return false
	}

	var thCount, tdCount int
	for cell := tr.FirstChild; cell != nil; cell = cell.NextSibling {
		if cell.Type == html.ElementNode {
			if cell.Data == "th" {
				thCount++
			} else if cell.Data == "td" {
				tdCount++
			}
		}
	}

	// If all cells are th, or majority are th, consider it a header
	return thCount > 0 && thCount >= tdCount
}
