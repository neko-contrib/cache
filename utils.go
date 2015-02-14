package cache

import (
	"fmt"
)

func Increment(val interface{}, n int64) (interface{}, error) {
	switch val.(type) {
	case int:
		val = val.(int) + int(n)
	case int32:
		val = val.(int32) + int32(n)
	case int64:
		val = val.(int64) + n
	case uint:
		val = val.(uint) + uint(n)
	case uint32:
		val = val.(uint32) + uint32(n)
	case uint64:
		val = val.(uint64) + uint64(n)
	default:
		return val, fmt.Errorf("the item is not an integer")
	}
	return val, nil
}

func Decrement(val interface{}, n int64) (interface{}, error) {
	switch val.(type) {
	case int:
		val = val.(int) - int(n)
	case int32:
		val = val.(int32) - int32(n)
	case int64:
		val = val.(int64) - n
	case uint:
		if val.(uint) >= uint(n) {
			val = val.(uint) - uint(n)
		} else {
			return val, fmt.Errorf("the item is less than 0")
		}
	case uint32:
		if val.(uint32) >= uint32(n) {
			val = val.(uint32) - uint32(n)
		} else {
			return val, fmt.Errorf("the item is less than 0")
		}
	case uint64:
		if val.(uint64) >= uint64(n) {
			val = val.(uint64) - uint64(n)
		} else {
			return val, fmt.Errorf("the item is less than 0")
		}
	default:
		return val, fmt.Errorf("the item is not an integer")
	}
	return val, nil
}
