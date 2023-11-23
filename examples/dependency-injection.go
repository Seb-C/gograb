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
