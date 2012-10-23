package db

import (
    "database/sql"
    "kx/dlog"
    _ "github.com/mattn/go-sqlite3"
)

var (
    agent *sql.DB
)

func init() {
    if dlog.FileExists(dlog.DbFile) {
        return
    }

    // create the table
    db, err := sql.Open(dlog.DbEngine, dlog.DbFile)
    if err != nil {
        panic(err)
    }

    defer db.Close()

    if _, err := db.Exec(dlog.SQL_CREATE_TABLE); err != nil {
        panic(err)
    }
}

func ImportResult(name string, r dlog.ReduceResult) {
}
