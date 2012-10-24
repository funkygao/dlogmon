package dlog

import "log"

const (
    LZOP_CMD          = "lzop"
    LZOP_OPTION       = "-dcf"
    EOL               = '\n'
    DLOG_BASE_DIR     = "/kx/dlog/"
    SAMPLER_HOST      = "100.123"
    FLAG_TIMESPAN_SEP = "-"
    KEYTYPE_SEP       = ":"
)

const (
    LOG_OPTIONS_DEBUG = log.Ldate | log.Lshortfile | log.Ltime | log.Lmicroseconds
    LOG_OPTIONS       = log.LstdFlags
    LOG_PREFIX_DEBUG  = "debug "
    LOG_PREFIX        = ""
)

const (
    availableMemory        = 24 << 30  // 24 GB
    avgMemoryPerWorker     = 800 << 10 // 800 KB
    MAX_CONCURRENT_WORKERS = availableMemory / avgMemoryPerWorker
)

const (
    VarDir   = "var"
    DbEngine = "sqlite3"
    DbFile   = VarDir + "/dlogmon.db"
)
