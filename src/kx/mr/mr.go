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

func (this KeyValues) Keys() interface{} {
    key := this.getOneKey()
    if key == nil {
        return nil
    }

    if _, ok := key.(string); ok {
        keys := make([]string, len(this))
        var i int
        for k, _ := range this {
            keys[i] = k.(string)
            i ++
        }

        return keys
    } else if _, ok := key.([2]string); ok {
        keys := make([][2]string, len(this))
        var i int
        for k, _ := range this {
            keys[i] = k.([2]string)
            i ++
        }

        return keys
    } else if _, ok := key.([3]string); ok {
        keys := make([][3]string, len(this))
        var i int
        for k, _ := range this {
            keys[i] = k.([3]string)
            i ++
        }

        return keys
    } else if _, ok := key.([4]string); ok {
        keys := make([][4]string, len(this))
        var i int
        for k, _ := range this {
            keys[i] = k.([4]string)
            i ++
        }

        return keys
    } else if _, ok := key.([5]string); ok {
        keys := make([][5]string, len(this))
        var i int
        for k, _ := range this {
            keys[i] = k.([5]string)
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
