/*
Schema:

    id: PK auto increment
    type: worker type
    ctx1: context e,g rid
    ctx2: context e,g uri
    k1: key e,g 
    k2: key
    k3: key
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
        type VARCHAR(30) NOT NULL, 
        ctx VARCHAR(100),
        k1 VARCHAR(100),
        k2 VARCHAR(100),
        k3 VARCHAR(100),
        num FLOAT, 
        ctime DATETIME DEFAULT CURRENT_TIMESTAMP)`
)
