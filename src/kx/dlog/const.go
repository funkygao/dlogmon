package dlog

import "log"

const (
    LZOP_CMD          = "lzop"
    LZOP_OPTION       = "-dcf"
    EOL               = '\n'
    DLOG_BASE_DIR     = "/kx/dlog/"
    SAMPLER_HOST      = "100.123"
    FLAG_TIMESPAN_SEP = "-"
)

const (
   LOG_OPTIONS_DEBUG = log.Ldate | log.Lshortfile | log.Ltime | log.Lmicroseconds
   LOG_OPTIONS = log.LstdFlags
   LOG_PREFIX_DEBUG = "debug "
   LOG_PREFIX = ""
)
