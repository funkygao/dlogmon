package mr

import (
	"fmt"
	T "kx/trace"
)

const (
	PRINT_FMT = "%s %60v %v\n"
)

// Factory
func NewKeyValue() KeyValue {
	return make(KeyValue)
}

// Factory
func NewKeyValues() KeyValues {
	return make(KeyValues)
}

// Self printable
func (this KeyValue) Println() {
	for k, v := range this {
		fmt.Printf(PRINT_FMT, T.CallerFuncName(1), k, v)
	}
}

// Self printable
func (this KeyValue) PrintByOrderedKey(sortedKeys interface{}) {
	for _, k := range sortedKeys.([]string) {
		fmt.Printf(PRINT_FMT, T.CallerFuncName(1), k, this[k])
	}
}

func (this KeyValue) DumpToSql(sortedKeys interface{}) {
	println()
	this.PrintByOrderedKey(sortedKeys)
}

// Self printable
func (this KeyValues) Println() {
	for k, v := range this {
		fmt.Printf(PRINT_FMT, T.CallerFuncName(1), k, v)
	}
}

func (this KeyValues) getOneKey() (key interface{}) {
    for k, _ := range this {
        key = k
        return
    }

    return
}

func (this KeyValues) Empty() bool {
    return len(this) == 0
}

func (this KeyValues) Keys() interface{} {
    key := this.getOneKey()
    if key == nil {
        return nil
    }

    var i int
    switch key.(type) {
    case string:
        keys := make([]string, len(this))
        for k, _ := range this {
            keys[i] = k.(string)
            i ++
        }

        return keys
    case [2]string:
        keys := make([][2]string, len(this))
        for k, _ := range this {
            keys[i] = k.([2]string)
            i ++
        }

        return keys
    case [3]string:
        keys := make([][3]string, len(this))
        for k, _ := range this {
            keys[i] = k.([3]string)
            i ++
        }

        return keys
    case [4]string:
        keys := make([][4]string, len(this))
        for k, _ := range this {
            keys[i] = k.([4]string)
            i ++
        }

        return keys
    case [5]string:
        keys := make([][5]string, len(this))
        for k, _ := range this {
            keys[i] = k.([5]string)
            i ++
        }

        return keys
    case [6]string:
        keys := make([][6]string, len(this))
        for k, _ := range this {
            keys[i] = k.([6]string)
            i ++
        }

        return keys
    }

	return nil
}

func (this KeyValues) Append(key interface{}, val interface{}) {
	if _, ok := this[key]; !ok {
		this[key] = make([]interface{}, 1)
		this[key][0] = val
	} else {
		this[key] = append(this[key], val)
	}
}

func (this KeyValues) AppendSlice(key interface{}, val []interface{}) {
	if _, ok := this[key]; !ok {
		this[key] = make([]interface{}, 1)
		this[key] = val
	} else {
		this[key] = append(this[key], val...)
	}
}
