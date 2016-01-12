package main

import (
	"testing"
)

func TestDereferencePath(t *testing.T) {
	url, err := dereferencePath("https://www.iopipe.com/")
	if err != nil {
		t.Error(err)
	}
	if url.Scheme != "https" {
		t.Error("Expected https")
	}
}

func TestDereferencePathPlainString(t *testing.T) {
	url, err := dereferencePath("plainstring")
	if err != nil {
		t.Error(err)
	}
	if url.Scheme != "https" {
		t.Error("Expected https")
	}
}

func TestDereferencePathHTTP(t *testing.T) {
	url, err := dereferencePath("http://address")
	if err != nil {
		t.Error(err)
	}
	if url.Scheme != "http" {
		t.Error("Expected http")
	}
}

func TestDereferenceObj(t *testing.T) {
	url, err := dereferencePath("http://127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	obj := dereferenceObj(url)
	if url != obj.path {
		t.Error("URL and url.path did not match")
	}
}
