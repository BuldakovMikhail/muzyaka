package lib

import (
	"os"
	"strings"
)

func ReadAllFilesFromArray(paths []string) ([][]byte, error) {
	var resultArr [][]byte

	for _, v := range paths {
		trimmedVal := strings.TrimRight(v, "\r\n")
		data, err := os.ReadFile(trimmedVal)
		if err != nil {
			return nil, err
		}
		resultArr = append(resultArr, data)
	}

	return resultArr, nil
}
