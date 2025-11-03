package config

import "testing"

func TestMaskDatabaseDSN(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{name: "Empty", in: "", out: ""},
		{name: "NoPassword", in: "user@tcp(localhost:3306)/db", out: "user@tcp(localhost:3306)/db"},
		{name: "WithPassword", in: "user:secret@tcp(localhost:3306)/db", out: "user:***@tcp(localhost:3306)/db"},
		{name: "URIStyle", in: "mysql://user:secret@localhost:3306/db?parseTime=true", out: "mysql://user:***@localhost:3306/db?parseTime=true"},
		{name: "FallbackWhenNoNeedle", in: "user:secret@tcp(localhost:3306)/db?password=secret", out: "user:***@tcp(localhost:3306)/db?password=secret"},
		{name: "UnixSocket", in: "user:secret@unix(/var/run/mysql.sock)/db", out: "user:***@unix(/var/run/mysql.sock)/db"},
		{name: "NoPasswordQuery", in: "user@tcp(localhost:3306)/db?password=secret", out: "user@tcp(localhost:3306)/db?password=secret"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskDatabaseDSN(tt.in); got != tt.out {
				t.Fatalf("MaskDatabaseDSN(%q) = %q, want %q", tt.in, got, tt.out)
			}
		})
	}
}

func TestDSNParse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want DSN
	}{
		{
			name: "ClassicTCP",
			in:   "user:secret@tcp(localhost:3306)/photoprism?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true",
			want: DSN{
				User:     "user",
				Password: "secret",
				Net:      "tcp",
				Server:   "localhost:3306",
				Name:     "photoprism",
				Params:   "charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true",
			},
		},
		{
			name: "URIStyle",
			in:   "mysql://user:secret@localhost:3306/photoprism?parseTime=true",
			want: DSN{
				Driver:   "mysql",
				User:     "user",
				Password: "secret",
				Server:   "localhost:3306",
				Name:     "photoprism",
				Params:   "parseTime=true",
			},
		},
		{
			name: "UnixSocket",
			in:   "user:secret@unix(/var/run/mysql.sock)/photoprism",
			want: DSN{
				User:     "user",
				Password: "secret",
				Net:      "unix",
				Server:   "/var/run/mysql.sock",
				Name:     "photoprism",
			},
		},
		{
			name: "FileDSN",
			in:   "file:/data/index.db?_busy_timeout=5000",
			want: DSN{
				Server: "file:/data",
				Name:   "index.db",
				Params: "_busy_timeout=5000",
			},
		},
		{
			name: "EmptyInput",
			in:   "",
			want: DSN{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDSN(tt.in)
			if got != tt.want {
				t.Fatalf("NewDSN(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}
