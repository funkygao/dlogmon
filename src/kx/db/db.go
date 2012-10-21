package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "os"
)

const (
    vardir = "var"
    engine = "sqlite3"
    dbfile = vardir + "dlogmon.db"
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

    var err error
    var dir *os.File
    if dir, err = os.Open(vardir); err != nil {
        panic("must run on top dir")
    }
    dir.Close()

    agent = new(DbAgent)
    if agent.DB, err = sql.Open(engine, dbfile); err != nil {
        panic(err)
    }

    return agent
}
