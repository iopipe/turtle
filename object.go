package main

import (
	"log"
	"strings"

	"io/ioutil"
	"net/http"
	"net/url"
)

/*******************************************************
 Object Mapper
*******************************************************/
type objectInterface struct {
	ClassID    string                 `json:"classid"`
	Properties map[string]interface{} `json:"properties"`
}

type Object struct {
	path *url.URL
}

func (object *Object) read() string {
	path := object.path.String()
	object.path.String()
	res, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	return string(body[:])
}

func (object *Object) update(content string) string {
	path := object.path.String()
	reader := strings.NewReader(content)

	res, err := http.Post(path, "application/json", reader)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body[:])
}
