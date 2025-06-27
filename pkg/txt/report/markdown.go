package report

import (
	"bytes"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// MarkdownTable returns a text-formatted table with caption, optionally as valid Markdown,
// so the output can be pasted into the docs.
func MarkdownTable(rows [][]string, cols []string, opt Options) string {
	// Escape Markdown.
	if opt.Valid {
		for i := range rows {
			for j := range rows[i] {
				if strings.ContainsRune(rows[i][j], '|') {
					rows[i][j] = strings.ReplaceAll(rows[i][j], "|", "\\|")
				}
				if strings.ContainsRune(rows[i][j], '*') {
					rows[i][j] = strings.ReplaceAll(rows[i][j], "* * *", "\\* \\* \\*")
				}
			}
		}
	}

	result := &bytes.Buffer{}

	var tableRenderer tw.Renderer
	var tableConfig tablewriter.Config

	if opt.Valid {
		tableRenderer = renderer.NewMarkdown()
		tableConfig = tablewriter.Config{
			Header: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignLeft}, Formatting: tw.CellFormatting{AutoFormat: -1}},
			Row: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
		}
	} else {
		tableRenderer = renderer.NewBlueprint()
		tableConfig = tablewriter.Config{
			Header: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignCenter}, Formatting: tw.CellFormatting{AutoFormat: -1}},
			Row: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
		}
	}

	// RenderFormat.
	table := tablewriter.NewTable(result,
		tablewriter.WithRenderer(tableRenderer),
		tablewriter.WithConfig(tableConfig),
	)

	// Set Caption.
	if opt.Caption != "" {
		table.Caption(tw.Caption{Text: opt.Caption})
	}

	table.Header(cols)
	_ = table.Bulk(rows)
	_ = table.Render()

	return result.String()
}
