package lib

import "os"

func ReadAllFilesFromArray(paths []string) ([][]byte, error) {
	var resultArr [][]byte

	for _, v := range paths {
		data, err := os.ReadFile(v)
		if err != nil {
			return nil, err
		}
		resultArr = append(resultArr, data)
	}

	return resultArr, nil
}
