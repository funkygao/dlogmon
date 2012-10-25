package db

import (
    "database/sql"
    // _ "github.com/mattn/go-sqlite3"
    "kx/mr"
    //"kx/util"
)

var (
    db                                 *sql.DB
    dbfile, dbengine, sql_create_table string
)

func Initialize(engine, file string) {
    /*
    dbfile = file
    dbengine = engine

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

    if _, err := db.Exec(SQL_CREATE_TABLE); err != nil {
        panic(err)
    }
    */
}

func ImportResult(name string, r mr.ReduceResult) {
    r.DumpToSql()
}
