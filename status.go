package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func getStatus() (ret map[string]string) {
	ret = map[string]string{
		"path": file,
		"date": "--/--/----",
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		return
	}
	ret["date"] = fileInfo.ModTime().Format("2006-01-02 15:04:05")
	return
}

func loadNames() []string {
	var data dataset
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return []string{} // probably no file found
	}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data.Names
}
