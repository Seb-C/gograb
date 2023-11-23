package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlockContent(t *testing.T) {
	cleanString := func(input string) string {
		return strings.TrimSpace(strings.ReplaceAll(input, "\t", ""))
	}

	t.Run("Matches the expected block", func(t *testing.T) {
		source := `
			package examples

			import (
				"database/sql"
				"log/slog"
				"net/http"
			)

			type Foo struct {
				//gograb:begin dependencies
				server *http.Server
				database *sql.DB
				logger *slog.Logger
				someValue int
				//gograb:end

				internalProperty string
				otherInternalStuff float64
			}

			func NewFoo(
			//go:generate go run ../gograb.go dependencies "\t([^ ]+) +([^ ]+)" "\t@1 @2,"
			//gograb:end
			) *Foo {
				return &Foo{
				}
			}
		`

		result, err := getBlockContent([]byte(source), "dependencies")

		assert.Nil(t, err)
		assert.Equal(t, cleanString(`
			server *http.Server
			database *sql.DB
			logger *slog.Logger
			someValue int
		`), cleanString(string(result)))
	})
	t.Run("Error if not found", func(t *testing.T) {
		source := `
			package examples

			import (
				"database/sql"
				"log/slog"
				"net/http"
			)

			type Foo struct {
				//gograb:begin not-dependencies
				server *http.Server
				database *sql.DB
				logger *slog.Logger
				someValue int
				//gograb:end

				internalProperty string
				otherInternalStuff float64
			}
		`

		_, err := getBlockContent([]byte(source), "dependencies")

		assert.NotNil(t, err)
	})
	t.Run("Error if multiple blocks", func(t *testing.T) {
		source := `
			package examples

			import (
				"database/sql"
				"log/slog"
				"net/http"
			)

			type Foo struct {
				//gograb:begin dependencies
				server *http.Server
				database *sql.DB
				logger *slog.Logger
				someValue int
				//gograb:end

				//gograb:begin dependencies
				internalProperty string
				otherInternalStuff float64
				//gograb:end
			}
		`

		_, err := getBlockContent([]byte(source), "dependencies")

		assert.NotNil(t, err)
	})
}
