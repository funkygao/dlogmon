package db

const (
    SQL_CREATE_TABLE = `create table dlog(
        id INTEGER PRIMARY KEY,
        name varchar(20), 
        key varchar(100), 
        num float, 
        ctime datetime)`
)
