package main

import (
        "github.com/robertkrimen/otto"

        //"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const FILTER_BASE string = "http://192.241.174.50/filters/"

type filterTuple struct {
	fromObjType string
	toObjType   string
}

func findFilter(fromObjType string, toObjType string) (func(input string) (string, error), error) {
	path := FILTER_BASE + fromObjType + "/" + toObjType;
        res, err := http.Get(path)
        if err != nil {
                log.Fatal(err)
        }
        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
                log.Fatal(err)
        }
        res.Body.Close()
        script := string(body[:])

	return func(input string) (string, error) {
		vm := otto.New()
		vm.Set("input", input)
		println("Executing script: " + script)
		val, err := vm.Run(script)
		if err != nil {
			return "", err
		}
		return val.ToString()
	}, nil

}

