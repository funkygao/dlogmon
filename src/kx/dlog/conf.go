package dlog

import (
    "github.com/kless/goconfig/config"
)

func loadConf(conf string) (*config.Config, error) {
    return config.ReadDefault(conf)
}
