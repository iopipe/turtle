package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/robertkrimen/otto"

	"errors"
	"log"
	"os"
	"path"
	"strings"

	"crypto/sha256"
	"io/ioutil"
	"net/http"
	"os/user"
)

func listScripts() {
	open()
	read()
	fancyOutput()
}

func exportScript(pipeline string, name string) {
	/*
			// Make directory...

			var cfgf = io.FileWriter //>>iopipe-cfg.json:
			var scriptf = io.Filewriter //>index.js
			var packagef = io.Filewriter //>package.json
			cfgf.write(`{"auth": {}}`)

			scriptf.write(`var iopipe = require("iopipe")
		        var config = require("./iopipe-cfg.json")
			iopipe.load_config(config)
			exports.run = function() {
		          iopipe.exec()
		        }`)

			packagef.write(`
			{
			  "name": "iopipe",
			  "private": true,
			  "version": "0.0.1",
			  "description": "iopipe sdk",
			  "author": "Eric Windisch",
			  "dependencies": {
			    "read-stream": "",
			    "request": ""
			  },
			  "main": "./iopipe.js"
			}
			`)
	*/
}

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

func fetchFilter(filterPath string) ([]byte, error) {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	path := path.Join(FILTER_BASE, filterPath)
	res, err = http.Get(path)
	if err != nil {
		return nil, err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	/* Verify digest */
	chksum := sha256.Sum256(body[:])
	if filterPath != string(chksum[:]) {
		return nil, errors.New("Checksum failure")
	}

	return body, nil
}

func writeCache(body []byte) {
	/* Verify digest */
	chksum := sha256.Sum256(body[:])
	diskPath := getCachePath(chksum)

	/* Write cache */
	if err = ioutil.WriteFile(diskPath, script, 0600); err != nil {
		return nil, err
	}
}

func getCachePath(name string) string {
	reqPathParts := strings.Split(name, "/")

	myuser, err := user.Current()
	if err != nil {
		return nil, err
	}
	pathParts := []string{myuser.HomeDir, ".iopipe", "filter_cache"}
	pathParts = append(pathParts, reqPathParts...)
	return path.Join(pathParts...)
}

func readFilterCache(name string) ([]byte, error) {
	diskPath := getCachePath(name)

	/* Do we have this cached? */
	if _, err := os.Stat(diskPath); err != nil {
		return "", err
	}
	script, err = ioutil.ReadFile(diskPath)
	return script[:], nil
}

func importScript(file string) {
	if file == "-" {
		fH := os.STDIN
	}

	chksum := sha256.Sum256(body[:])
	writeFilter(chksum, body[:])
}

func getFilter(filterPath string) (func(input string) (string, error), error) {
	var script []byte
	var err error

	diskPath := getCachePath(filterPath)

	/* Do we have this cached? */
	if script, err := readFilterCache(filterPath); err != nil {
		return makeFilter(string(script[:]))
	} else {
		return "", err
	}

	/* If not, fetch */
	if script, err = fetchFilter(filterPath); err != nil {
		return nil, err
	}
	if err = writeCache(script); err != nil {
		return nil, err
	}

	return makeFilter(string(script[:]))
}

func getPipeline(filterPath string) (func(input string) (string, error), error) {
	var script []byte
	var err error
	reqPathParts := strings.Split(filterPath, "/")
	myuser, err := user.Current()
	if err != nil {
		return nil, err
	}
	pathParts := []string{myuser.HomeDir, ".iopipe", "filter_cache"}
	pathParts = append(pathParts, reqPathParts...)
	diskPath := path.Join(pathParts...)
	//myuser.HomeDir, ".iopipe", "filter_cache", pathParts...)

	/* Do we have this cached? */
	if _, err := os.Stat(diskPath); err == nil {
		script, err = ioutil.ReadFile(diskPath)
		return makeFilter(string(script[:]))
	}
	/* If not, fetch */
	if script, err = fetchFilter(filterPath); err != nil {
		return nil, err
	}
	/* Write cache */
	if err = ioutil.WriteFile(diskPath, script, 0600); err != nil {
		return nil, err
	}
	return makeFilter(string(script[:]))
}

func findFilters(fromObjType string, toObjType string) (string, error) {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	path := path.Join(FILTER_BASE, fromObjType, toObjType)
	res, err = http.Get(path)
	if err != nil {
		return "", err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body.Close()
	response := string(body[:])
	return response, nil
}

func findPipelines(fromObjType string, toObjType string) (string, error) {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	path := path.Join(FILTER_BASE, fromObjType, toObjType)
	res, err = http.Get(path)
	if err != nil {
		return "", err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body.Close()
	response := string(body[:])
	return response, nil
}

func publishPipeline(pipeline string) {
}

func subscribePipeline(pipeline string) {
}
