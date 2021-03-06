package dlog

import "log"

const (
    W_AMF  = "amf"
    W_KXI  = "kxi"
    W_UNI  = "uni"
    W_NOOP = "noop"
    W_FILE = "file"
)

const (
    LZOP_CMD                = "lzop"
    LZOP_OPTION             = "-dcf"
    EOL                     = '\n'
    DEFAULT_DLOG_BASE_DIR   = "/kx/dlog/"
    SAMPLER_HOST            = "100.123"
    FLAG_TIMESPAN_SEP       = "-"
    KEYTYPE_SEP             = ":"
    NIL_KEY                 = ""
    LINES_PER_FILE          = 565000 // deprecated
    PROGRESS_LINES_STEP     = 50000  // TODO calculated instead of const
    PROGRESS_CHAN_BUF       = 1000
    LINE_CHANBUF_PER_WORKER = 100
)

const (
    LOG_OPTIONS_DEBUG = log.Ldate | log.Lshortfile | log.Ltime | log.Lmicroseconds
    LOG_OPTIONS       = log.LstdFlags | log.Lshortfile
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

const (
    TABLE_AMF = "t_amf"
)

const (
    LogEmergency = LogLevel( 1 << iota)
    LogAlert
    LogCritical
    LogError
    LogWarning
    LogNotice
    LogInfo
    LogDebug

    DefaultLogLevel = LogInfo
)
