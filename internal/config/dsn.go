package config

import (
	"regexp"
	"strings"
	"unicode"
)

// dsnPattern is a regular expression matching a database DSN string.
var dsnPattern = regexp.MustCompile(
	`^((?P<driver>.*):\/\/)?(?:(?P<user>.*?)(?::(?P<password>.*))?@)?` +
		`(?:(?P<net>[^\(]*)(?:\((?P<server>[^\)]*)\))?)?` +
		`\/(?P<name>.*?)` +
		`(?:\?(?P<params>[^\?]*))?$`)

// dsnPostgresPasswordPattern is a regular expression matching a password in a PostgreSQL-style database DSN string.
var dsnPostgresPasswordPattern = regexp.MustCompile(`(?i)(password\s*=\s*)("[^"]*"|'[^']*'|\S+)`)

// DSN represents parts of a data source name.
type DSN struct {
	DSN      string
	Driver   string
	User     string
	Password string
	Net      string
	Server   string
	Name     string
	Params   string
}

// NewDSN creates a new DSN struct from a string.
func NewDSN(dsn string) DSN {
	d := DSN{DSN: dsn}
	d.parse()
	return d
}

// String returns the original DSN string.
func (d *DSN) String() string {
	return d.DSN
}

// MaskPassword hides the password portion of a DSN while leaving the rest untouched for logging/reporting.
func (d *DSN) MaskPassword() (s string) {
	if d.DSN == "" || d.Password == "" {
		return d.DSN
	}

	s = d.DSN

	// Mask password in regular DSN.
	needle := ":" + d.Password + "@"
	if strings.Contains(s, needle) {
		return strings.Replace(s, needle, ":***@", 1)
	}

	// Mask password in PostgreSQL-style DSN.
	if d.Driver == Postgres || strings.Contains(s, "password=") {
		return dsnPostgresPasswordPattern.ReplaceAllStringFunc(s, func(segment string) string {
			matches := dsnPostgresPasswordPattern.FindStringSubmatch(segment)
			if len(matches) != 3 {
				return segment
			}

			prefix := matches[1]
			value := matches[2]
			unquoted := strings.Trim(value, `'"`)

			if unquoted != d.Password {
				return segment
			}

			switch {
			case strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`):
				return prefix + `"` + "***" + `"`
			case strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`):
				return prefix + `'` + "***" + `'`
			default:
				return prefix + "***"
			}
		})
	}

	// Return DSN with masked password.
	return s
}

// Parse parses a data source name string.
func (d *DSN) parse() {
	// Assume a regular DSN, and if parsing fails, treat it as a PostgreSQL-style DSN.
	if matches := dsnPattern.FindStringSubmatch(d.DSN); len(matches) > 0 {
		names := dsnPattern.SubexpNames()

		for i, match := range matches {
			switch names[i] {
			case "driver":
				d.Driver = match
			case "user":
				d.User = match
			case "password":
				d.Password = match
			case "net":
				d.Net = match
			case "server":
				d.Server = match
			case "name":
				d.Name = match
			case "params":
				d.Params = match
			}
		}

		if d.Net != "" && d.Server == "" {
			d.Server = d.Net
			d.Net = ""
		}
	} else {
		// Parse PostgreSQL-style DSN
		d.parsePostgres()
	}
}

// parsePostgres extracts connection settings from PostgreSQL key/value style DSNs.
func (d *DSN) parsePostgres() bool {
	fields, ok := d.splitKeyValue(d.DSN)

	if !ok {
		return false
	}

	values := make(map[string]string, len(fields))
	order := make([]string, 0, len(fields))

	for _, field := range fields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			return false
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" {
			return false
		}

		// Trim optional surrounding quotes.
		val = strings.Trim(val, `"`)

		values[key] = val
		order = append(order, key)
	}

	name := values["dbname"]

	if name == "" {
		if alt := values["database"]; alt != "" {
			name = alt
		} else {
			return false
		}
	}

	d.Driver = Postgres
	d.User = values["user"]
	d.Password = values["password"]
	d.Name = name

	host := values["host"]
	port := values["port"]

	switch {
	case host != "" && port != "":
		d.Server = host + ":" + port
	case host != "":
		d.Server = host
	case port != "":
		d.Server = ":" + port
	}

	// Remove canonical keys so remaining values become Params.
	delete(values, "user")
	delete(values, "password")
	delete(values, "dbname")
	delete(values, "database")
	delete(values, "host")
	delete(values, "port")

	params := make([]string, 0, len(values))

	for _, key := range order {
		if val, ok := values[key]; ok {
			if strings.Contains(val, " ") {
				val = `"` + val + `"`
			}
			params = append(params, key+"="+val)
		}
	}

	if len(params) > 0 {
		d.Params = strings.Join(params, " ")
	}

	return true
}

// splitKeyValue tokenizes PostgreSQL key/value DSNs, supporting quoted values with spaces.
func (d *DSN) splitKeyValue(input string) ([]string, bool) {
	runes := []rune(strings.TrimSpace(input))

	if len(runes) == 0 || !strings.Contains(input, "=") {
		return nil, false
	}

	var (
		tokens    []string
		current   strings.Builder
		inQuotes  bool
		quoteRune rune
	)

	flush := func() {
		if current.Len() == 0 {
			return
		}
		tokens = append(tokens, current.String())
		current.Reset()
	}

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		switch {
		case inQuotes && r == '\\':
			if i+1 < len(runes) {
				current.WriteRune(runes[i+1])
				i++
			}
		case r == '\'' || r == '"':
			if inQuotes {
				if r == quoteRune {
					inQuotes = false
				} else {
					current.WriteRune(r)
				}
			} else {
				inQuotes = true
				quoteRune = r
			}
		case unicode.IsSpace(r):
			if inQuotes {
				current.WriteRune(r)
			} else {
				flush()
			}
		default:
			current.WriteRune(r)
		}
	}

	if inQuotes {
		return nil, false
	}

	flush()

	if len(tokens) == 0 {
		return nil, false
	}

	for _, token := range tokens {
		if !strings.Contains(token, "=") {
			return nil, false
		}
	}

	return tokens, true
}

// MaskDatabaseDSN hides the password portion of a DSN while leaving the rest untouched for logging/reporting.
func MaskDatabaseDSN(dsn string) string {
	if dsn == "" {
		return ""
	}

	// Parse database DSN.
	d := NewDSN(dsn)

	// Return original DSN if no password was found.
	if d.Password == "" {
		return dsn
	}

	// Return DSN with masked password.
	return d.MaskPassword()
}
