package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadJson[T any](filePath string, collection *T) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, collection); err != nil {
		return err
	}

	return nil
}

func WriteJson[T any](filePath string, collection *T) error {
	outJson, err := json.Marshal(collection)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, outJson, 0666); err != nil {
		return err
	}

	return nil
}
