package db

import (
    "database/sql"
    "kx/util"
    "kx/mr"
    _ "github.com/mattn/go-sqlite3"
)

var (
    db *sql.DB
    dbfile, dbengine, sql_create_table string
)

func Initialize(engine, file, sql_create string) {
    dbfile = file
    dbengine = engine
    sql_create_table = sql_create

    if util.FileExists(dbfile) {
        return
    }

    // create the table
    var err error
    db, err = sql.Open(dbengine, dbfile)
    if err != nil {
        panic(err)
    }

    defer db.Close()

    if _, err := db.Exec(sql_create_table); err != nil {
        panic(err)
    }
}

func ImportResult(name string, r mr.ReduceResult) {
    r.Println()
}
