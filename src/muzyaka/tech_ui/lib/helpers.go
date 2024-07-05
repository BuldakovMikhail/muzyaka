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

func SaveFile(fileName string, content []byte) error {
	err := os.WriteFile(fileName, content, 0666)
	if err != nil {
		return err
	}

	return nil
}

func IsPNGFormat(fileName string) bool {
	splittedName := strings.Split(fileName, ".")
	ext := splittedName[len(splittedName)-1]
	return ext == "png"
}

func IsMP3Format(fileName string) bool {
	splittedName := strings.Split(fileName, ".")
	ext := splittedName[len(splittedName)-1]
	return ext == "mp3"
}
