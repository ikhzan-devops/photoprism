package report

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// Credentials returns a text-formatted table with credentials.
func Credentials(idName, idValue, secretName, secretValue string) string {
	result := &bytes.Buffer{}
	table := tablewriter.NewWriter(result)

	// Set values.
	rows := make([][]string, 2)
	rows[0] = []string{idName, secretName}
	rows[1] = []string{idValue, secretValue}

	_ = table.Bulk(rows)
	_ = table.Render()

	return result.String()
}
