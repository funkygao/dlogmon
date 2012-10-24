/*
Schema:

    id: PK auto increment
    k1:
    k2:
    type: worker type
    n:
    min:
    max:
    sum:
    mean:
    std:

*/
package db

const (
    SQL_CREATE_TABLE = `CREATE TABLE IF NOT EXISTS dlog(
        id INTEGER PRIMARY KEY,
        type VARCHAR(30), 
        k1 VARCHAR(100),
        k2 VARCHAR(100),
        k3 VARCHAR(100),
        num FLOAT, 
        ctime DATETIME DEFAULT CURRENT_TIMESTAMP)`
)
