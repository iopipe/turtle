package main

import (
	"net/url"
)

func dereferencePath(reqPath string) (*url.URL, error) {
	path, err := url.Parse(reqPath)
	if err != nil {
		return nil, err
	}
	if path.Scheme == "" {
		path.Scheme = "https"
	}
	return path, nil
}

// Download resource at path &
// validate resource matches declared object type definition.
func dereferenceObj(path *url.URL) *Object {
	obj := new(Object)
	obj.path = path
	return obj
}
