package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var assertSourceEquivalent = func(t *testing.T, expect string, gotBytes []byte) {
	t.Helper()

	expect = strings.TrimSpace(strings.ReplaceAll(expect, "\t", ""))
	got := strings.TrimSpace(strings.ReplaceAll(string(gotBytes), "\t", ""))

	assert.Equal(t, expect, got)
}

func TestGetSource(t *testing.T) {
	t.Run("Matches the expected block", func(t *testing.T) {
		code := `
			type Foo struct {
				//gograb:source
				server *http.Server
				database *sql.DB
				logger *slog.Logger
				someValue int
				//gograb:end

				internalProperty string
				otherInternalStuff float64
			}

			func NewFoo(
			//gograb:target foo bar
			//gograb:end
			) *Foo {
				return &Foo{
					//gograb:target foo bar
					//gograb:end
				}
			}
		`

		result, err := getSource([]byte(code))

		assert.Nil(t, err)
		assertSourceEquivalent(t, `
			server *http.Server
			database *sql.DB
			logger *slog.Logger
			someValue int
		`, result)
	})
	t.Run("Error if not found", func(t *testing.T) {
		code := `
			type Foo struct {
				server *http.Server
				database *sql.DB
				logger *slog.Logger
				someValue int
				//gograb:end

				internalProperty string
				otherInternalStuff float64
			}
		`

		_, err := getSource([]byte(code))

		assert.NotNil(t, err)
	})
	t.Run("Error if multiple blocks", func(t *testing.T) {
		code := `
			type Foo struct {
				//gograb:source
				server *http.Server
				database *sql.DB
				logger *slog.Logger
				someValue int
				//gograb:end

				//gograb:source
				internalProperty string
				otherInternalStuff float64
				//gograb:end
			}
		`

		_, err := getSource([]byte(code))

		assert.NotNil(t, err)
	})
}
