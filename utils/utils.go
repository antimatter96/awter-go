package utils

import (
	"fmt"
)

// MapMerge merges the two given maps and returns it
func MapMerge(map1, map2 map[string]interface{}) map[string]interface{} {
	var biggerMap, smallerMap *map[string]interface{}
	biggerMap = &map1
	smallerMap = &map2

	if len(map1) > len(map2) {
		biggerMap = &map2
		smallerMap = &map1
	}
	for k, v := range *smallerMap {
		v2, ok := (*biggerMap)[k]
		if ok {
			if v != v2 {
				panic(fmt.Errorf("Key already present, with different value"))
			}
		}
		(*biggerMap)[k] = v
	}
	return *biggerMap
}
