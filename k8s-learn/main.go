package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func main() {
	ReadConfig()
}

type RService struct {
	Name     string   `yaml:"name"`
	Hostname string   `yaml:"hostname"`
	Path     string   `yaml:"path"`
	Rewrite  Rewrite  `yaml:"rewrite"`
	Targets  []Target `yaml:"targets"`
}

type Rewrite struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Old string `yaml:"old"`
	New string `yaml:"new"`
}

type Target struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func ReadConfig() (*RService, error) {
	file := "config.yml"

	res, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf(" 180 err: %+v", err)
		return nil, err
	}

	fmt.Printf("====== info: %+v", string(res))

	svc := &RService{}

	err = yaml.Unmarshal(res, svc)
	if err != nil {
		fmt.Printf("189 err: %+v", err)
		return nil, err
	}
	fmt.Printf("\n========================= svc: %+v ====================\n", *svc)
	return svc, nil
}
