package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

const (
    engine = "sqlite3"
    dbfile = "var/dlogmon.db"
)

var (
    agent *DbAgent
)

type DbAgent struct {
    *sql.DB
}

// Factory of DbAgent
func New() *DbAgent {
    if agent != nil {
        return agent
    }

    agent = new(DbAgent)
    var err error
    if agent.DB, err = sql.Open(engine, dbfile); err != nil {
        panic(err)
    }

    return agent
}

