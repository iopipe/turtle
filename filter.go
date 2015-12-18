package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/robertkrimen/otto"

	//"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const FILTER_BASE string = "http://192.241.174.50/filters/"
const REQUIREJS_URL string = "http://requirejs.org/docs/release/2.1.22/r.js"

type filterTuple struct {
	fromObjType string
	toObjType   string
}

func makeFilter(script string) (func(input string) (string, error), error) {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	res, err = http.Get(REQUIREJS_URL)
	if err != nil {
		log.Fatal(err)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	rjs := string(body[:])

	return func(input string) (string, error) {
		vm := otto.New()

		logrus.Debug("Adding RequireJS")
		vm.Run(rjs)

		vm.Set("input", input)
		logrus.Debug("Executing script: " + script)
		val, err := vm.Run(script)
		if err != nil {
			return "", err
		}
		return val.ToString()
	}, nil
}

func findFilter(fromObjType string, toObjType string) (func(input string) (string, error), error) {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	path := FILTER_BASE + fromObjType + "/" + toObjType
	res, err = http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	script := string(body[:])

	return makeFilter(script)
}
