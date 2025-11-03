package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDSN_Parse(t *testing.T) {
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
				Driver:   MySQL,
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
				Driver:   MySQL,
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
				Driver:   MySQL,
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
				Driver: SQLite3,
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
				Driver: SQLite3,
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

func TestDSN_HostAndPort(t *testing.T) {
	tests := []struct {
		name string
		in   string
		host string
		port int
	}{
		{
			name: "MySQLTCP",
			in:   "user:secret@tcp(localhost:3307)/photoprism?parseTime=true",
			host: "localhost",
			port: 3307,
		},
		{
			name: "MySQLIPv6",
			in:   "user:secret@tcp([2001:db8::1]:3307)/photoprism",
			host: "2001:db8::1",
			port: 3307,
		},
		{
			name: "MySQLDefaultPort",
			in:   "user:secret@tcp(mysql.local)/photoprism",
			host: "mysql.local",
			port: 3306,
		},
		{
			name: "PostgresURL",
			in:   "postgres://user:secret@localhost/mydb",
			host: "localhost",
			port: 5432,
		},
		{
			name: "PostgresKeyValue",
			in:   "user=alice password=secret host=/var/run/postgresql port=6432 dbname=app",
			host: "/var/run/postgresql",
			port: 6432,
		},
		{
			name: "PostgresPortOnly",
			in:   "user=alice password=secret port=5433 dbname=app",
			host: "",
			port: 5433,
		},
		{
			name: "SQLite",
			in:   "file:/data/index.db",
			host: "",
			port: 0,
		},
		{
			name: "InvalidPortFallback",
			in:   "user:secret@tcp(localhost:abc)/photoprism",
			host: "localhost",
			port: 3306,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDSN(tt.in)
			assert.Equal(t, tt.host, d.Host())
			assert.Equal(t, tt.port, d.Port())
		})
	}
}

func TestDSN_MaskPassword(t *testing.T) {
	d := NewDSN("user:secret@tcp(localhost:3306)/db")
	assert.Equal(t, "user:***@tcp(localhost:3306)/db", d.MaskPassword())

	p := NewDSN("user=alice password=s3cr3t dbname=app")
	assert.Equal(t, "user=alice password=*** dbname=app", p.MaskPassword())

	noPass := NewDSN("user@tcp(localhost:3306)/db")
	assert.Equal(t, "user@tcp(localhost:3306)/db", noPass.MaskPassword())
}

func TestDSN_ParsePostgres(t *testing.T) {
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
			d := DSN{DSN: tt.in}
			ok := d.parsePostgres()

			assert.Equal(t, tt.in, d.String())

			if ok != tt.ok {
				t.Fatalf("parsePostgres(%q) ok=%v, want %v", tt.in, ok, tt.ok)
			}

			if ok && d != tt.want {
				t.Fatalf("parsePostgres(%q) = %#v, want %#v", tt.in, d, tt.want)
			}
		})
	}
}

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
