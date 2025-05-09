package model

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Element struct {
	Name    string     `json:"name"`
	Recipes [][]string `json:"recipes"`
	Tier    int        `json:"tier"`
}

type Data struct {
	Elements []Element `json:"elements"`
}

func LoadElementsFromFile(filename string) (Data, error) {
	data := Data{}
	filePath := filepath.Join("..", "data", filename)
	file, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}