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

func getFilter(filterPath string) (func(input string) (string, error), error) {
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

func findFilters(fromObjType string, toObjType string) string {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	path := path.Join(FILTER_BASE, fromObjType, toObjType)
	res, err = http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	response := string(body[:])
	return response
}

func findPipelines(fromObjType string, toObjType string) string {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	path := path.Join(FILTER_BASE, fromObjType, toObjType)
	res, err = http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	response := string(body[:])
	return response
}
