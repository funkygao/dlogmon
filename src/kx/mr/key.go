package mr

import (
    "strings"
)

func NewKey(keys ...string) Key {
    k := strings.Join(keys, KEY_SEP)
    this := Key{key: k, keySize: len(keys)}
    return this
}

func (this Key) Keys() []string {
    return strings.Split(this.key, KEY_SEP)
}

func (this Key) Less(that Key) bool {
    return this.key < that.key
}

func NewGroupKey(group string, keys ...string) GroupKey {
    k := strings.Join(keys, KEY_SEP)
    this := GroupKey{group:group, Key:Key{key: k, keySize: len(keys)}}
    return this
}

func (this GroupKey) Less(that GroupKey) bool {
    if this.group != that.group {
        return this.group < that.group
    }
    return this.key < that.key
}

func (this GroupKey) KeySize() int {
    return this.keySize
}

func (this GroupKey) Group() string {
    return this.group
}

func (this GroupKey) Keys() []string {
    return strings.Split(this.key, KEY_SEP)
}
