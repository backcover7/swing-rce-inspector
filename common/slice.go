package common

import (
	"fmt"
)

func RemoveDup(originals []string) ([]string, error) {
	temp := map[string]struct{}{}
	result := make([]string, 0, len(originals))
	for _, item := range originals {
		key := fmt.Sprint(item)
		if _, ok := temp[key]; !ok {
			temp[key] = struct{}{}
			result = append(result, item)
		}
	}
	return result, nil
}
