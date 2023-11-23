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

func TestGetBlockContent(t *testing.T) {
	t.Run("Matches the expected block", func(t *testing.T) {
		source := `
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
			//go:generate gograb ...
			//gograb:end
			) *Foo {
				return &Foo{
				}
			}
		`

		result, err := getBlockContent([]byte(source), "dependencies")

		assert.Nil(t, err)
		assertSourceEquivalent(t, `
			server *http.Server
			database *sql.DB
			logger *slog.Logger
			someValue int
		`, result)
	})
	t.Run("Error if not found", func(t *testing.T) {
		source := `
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

func TestReplaceTargetBlockContent(t *testing.T) {
	t.Run("With a single block", func(t *testing.T) {
		source := `
			func NewFoo(
			//go:generate gograb ...
			some
			existing
			content
			//gograb:end
			) *Foo {
				return &Foo{}
			}
		`

		result, err := replaceTargetBlockContent([]byte(source), "new content", 3)

		assert.Nil(t, err)
		assertSourceEquivalent(t, `
			func NewFoo(
			//go:generate gograb ...
			new content
			//gograb:end
			) *Foo {
				return &Foo{}
			}
		`, result)
	})
	t.Run("With multiple blocks", func(t *testing.T) {
		source := `
			//go:generate gograb ...
			should not be replaced
			//gograb:end
			func NewFoo(
			//go:generate gograb ...
			some existing content
			//gograb:end
			) *Foo {
				return &Foo{}
				//go:generate gograb ...
				should also not be replaced
				//gograb:end
			}
		`

		result, err := replaceTargetBlockContent([]byte(source), "new content", 6)

		assert.Nil(t, err)
		assertSourceEquivalent(t, `
			//go:generate gograb ...
			should not be replaced
			//gograb:end
			func NewFoo(
			//go:generate gograb ...
			new content
			//gograb:end
			) *Foo {
				return &Foo{}
				//go:generate gograb ...
				should also not be replaced
				//gograb:end
			}
		`, result)
	})
	t.Run("Error if missing end block", func(t *testing.T) {
		source := `
			func NewFoo(
			//go:generate gograb ...
			some existing content
			) *Foo {
				return &Foo{}
			}
		`

		_, err := replaceTargetBlockContent([]byte(source), "new content", 3)

		assert.NotNil(t, err)
	})
}
