package mr

import (
    "strings"
)

func NewKey(keys ...string) Key {
    k := strings.Join(keys, KEY_SEP)
    this := Key{key: k}
    return this
}

func (this Key) Keys() []string {
    return strings.Split(this.key, KEY_SEP)
}

func (this Key) Less(that Key, asc bool) bool {
    if asc {
        return this.key < that.key
    }
    return that.key < this.key
}

func NewGroupKey(group string, keys ...string) GroupKey {
    k := strings.Join(keys, KEY_SEP)
    this := GroupKey{group:group, Key:Key{key: k}}
    return this
}

func (this GroupKey) Less(that GroupKey, asc bool) bool {
    if this.group != that.group {
        if asc {
            return this.group < that.group
        } else {
            return that.group < this.group
        }
    }
    if asc {
        return this.key < that.key
    }
    return that.key < this.key
}

func (this GroupKey) Group() string {
    return this.group
}

func (this GroupKey) Keys() []string {
    return strings.Split(this.key, KEY_SEP)
}
