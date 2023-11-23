package examples

import (
	"database/sql"
	"log/slog"
	"net/http"
)

//go:generate go run ../gograb.go

type Foo struct {
	//gograb:source
	server    *http.Server
	database  *sql.DB
	logger    *slog.Logger
	someValue int
	//gograb:end

	internalProperty   string
	otherInternalStuff float64
}

func NewFoo(
	//gograb:target "\t([^ \n]+) +([^ \n]+)" "\t$1 $2,"
	internalProperty string,
	//gograb:end
) *Foo {
	return &Foo{
		//gograb:target "\t([^ \n]+) +([^ \n]+)" "\t\t$1: $1,"
		internalProperty: internalProperty,
		//gograb:end
	}
}
