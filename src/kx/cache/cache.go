package cache

import (
    "sync"
    "time"
)

const (
    SIG_STOP = true
    SIG_RESTART = false
)

type RetrievalFunc func() (interface{}, error)

type CachedValue struct {
    value interface{}
    retrievalFunc RetrievalFunc
    ttl int64
    mutex sync.Mutex
    ticker *time.Ticker
    chSignal chan bool
}

func NewCachedValue(r RetrievalFunc, ttl int64) *CachedValue {
    this := new(CachedValue)
    this.ttl = ttl
    this.retrievalFunc = r
    this.chSignal = make(chan bool)
    this.ticker = time.NewTicker(time.Duration(ttl))

    go this.houseKeep()

    return this
}

func (this *CachedValue) Value() (v interface{}, err error) {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    // cache hit
    if this.value != nil {
        return this.value, nil
    }

    // missed
    if this.value, err = this.retrievalFunc(); err != nil {
        this.value = nil
        return this.value, err
    }

    // retrieval is ok, restart ticker and return value
    this.ticker.Stop()

    this.ticker = time.NewTicker(time.Duration(this.ttl))
    this.chSignal <- SIG_RESTART

    return this.value, nil
}

func (this *CachedValue) Clear() {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    this.value = nil
}

func (this *CachedValue) Stop() {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    this.value, this.retrievalFunc = nil, nil

    this.ticker.Stop()

    this.chSignal <- SIG_STOP
}

func (this *CachedValue) houseKeep() {
    for {
        select {
        case <- this.ticker.C:
            // ttl expired
            this.Clear()
        case stop:= <- this.chSignal:
            if stop {
                // leave the endless loop
                return
            }
        }
    }
}
