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
    LINES_PER_FILE    = 600000 // deprecated
    PROGRESS_LINES_STEP = 50000 // TODO calculated instead of const
    PROGRESS_CHAN_BUF = 1000
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
    DbShard  = "0" // 可以任意扩充总的容量
    VarDir   = "var"
    DbEngine = "sqlite3"
    DbFile   = VarDir + "/dlogmon" + DbShard + ".db"
)
