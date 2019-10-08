// common functions used across multiple files

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func getFileFullPath(filename string) (string, error) {
	if !getDebugMode() {
		homePath := os.Getenv("HOME")
		if homePath == "" {
			return "", errors.New("Use Linux and set your $HOME variable you filthy casual")
		}
		return fmt.Sprintf("%s/.config/gapbot/%s", homePath, filename), nil
	}
	return fmt.Sprintf("./%s", filename), nil
}

func loadJSON(filename string, v interface{}) error {
	filePath, err := getFileFullPath(filename)
	if err != nil {
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func writeJSON(filename string, v interface{}) error {
	filePath, err := getFileFullPath(filename)
	if err != nil {
		return err
	}
	marshalledJSON, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, marshalledJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}

func itemInSlice(item string, slice []string) bool {
	for _, v := range slice {
		if item == v {
			return true
		}
	}
	return false
}
