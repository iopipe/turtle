package main

import (
	"log"
	"net/url"
)

// Create an ObjPath from a string
func dereferencePath(reqPath string) *url.URL {
	path, err := url.Parse(reqPath)
	if err != nil {
		log.Fatal(err)
	}
	if path.Scheme == "" {
		path.Scheme = "https"
	}
	return path
}

// Download resource at path &
// validate resource matches declared object type definition.
func dereferenceObj(path *url.URL) *Object {
	obj := new(Object)
	obj.path = path
	return obj
}
