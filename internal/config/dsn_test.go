package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		{name: "Postgres", in: "user=alice password=s3cr3t dbname=app", out: "user=alice password=*** dbname=app"},
		{name: "PostgresQuoted", in: "user=alice password=\"s ec ret\" dbname=app", out: "user=alice password=\"***\" dbname=app"},
		{name: "PostgresSingleQuoted", in: "password='secret' user=alice dbname=app", out: "password='***' user=alice dbname=app"},
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
				DSN:      "user:secret@tcp(localhost:3306)/photoprism?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true",
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
				DSN:      "mysql://user:secret@localhost:3306/photoprism?parseTime=true",
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
				DSN:      "user:secret@unix(/var/run/mysql.sock)/photoprism",
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
				DSN:    "file:/data/index.db?_busy_timeout=5000",
				Server: "file:/data",
				Name:   "index.db",
				Params: "_busy_timeout=5000",
			},
		},
		{
			name: "SQLite",
			in:   "/index.db?_busy_timeout=5000",
			want: DSN{
				DSN:    "/index.db?_busy_timeout=5000",
				Server: "",
				Name:   "index.db",
				Params: "_busy_timeout=5000",
			},
		},
		{
			name: "PostgresKeyValue",
			in:   "user=alice password=s3cr3t dbname=app host=db.internal port=5432 connect_timeout=5 sslmode=require",
			want: DSN{
				DSN:      "user=alice password=s3cr3t dbname=app host=db.internal port=5432 connect_timeout=5 sslmode=require",
				Driver:   Postgres,
				User:     "alice",
				Password: "s3cr3t",
				Server:   "db.internal:5432",
				Name:     "app",
				Params:   "connect_timeout=5 sslmode=require",
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
			assert.Equal(t, tt.in, got.String())
			if got != tt.want {
				t.Fatalf("NewDSN(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestDSNParsePostgres(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want DSN
		ok   bool
	}{
		{
			name: "Basic",
			in:   "user=alice password=s3cr3t dbname=app",
			want: DSN{
				DSN:      "user=alice password=s3cr3t dbname=app",
				Driver:   Postgres,
				User:     "alice",
				Password: "s3cr3t",
				Name:     "app",
			},
			ok: true,
		},
		{
			name: "WithHostPortAndParams",
			in:   "user=alice password=s3cr3t dbname=app host=db.internal port=5432 connect_timeout=5 sslmode=require",
			want: DSN{
				DSN:      "user=alice password=s3cr3t dbname=app host=db.internal port=5432 connect_timeout=5 sslmode=require",
				Driver:   Postgres,
				User:     "alice",
				Password: "s3cr3t",
				Server:   "db.internal:5432",
				Name:     "app",
				Params:   "connect_timeout=5 sslmode=require",
			},
			ok: true,
		},
		{
			name: "QuotedValues",
			in:   `user="alice" password="s ec ret" dbname="app" host=db.internal`,
			want: DSN{
				DSN:      `user="alice" password="s ec ret" dbname="app" host=db.internal`,
				Driver:   Postgres,
				User:     "alice",
				Password: "s ec ret",
				Server:   "db.internal",
				Name:     "app",
			},
			ok: true,
		},
		{
			name: "MissingDatabase",
			in:   "user=alice host=db.internal",
			want: DSN{DSN: "user=alice host=db.internal"},
			ok:   false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var got = NewDSN(tt.in)
			ok := got.parsePostgres()

			assert.Equal(t, tt.in, got.String())

			if ok != tt.ok {
				t.Fatalf("parsePostgres(%q) ok=%v, want %v", tt.in, ok, tt.ok)
			}

			if ok && got != tt.want {
				t.Fatalf("parsePostgres(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}
