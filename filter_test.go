package main

import (
	"testing"
)

func TestMakeRunFilter(t *testing.T) {
	fun, err := makeFilter(`module.exports = function(x, cxt) { cxt.done("hello") }`)
	if err != nil {
		t.Error(err)
	}
	result, err := fun("")
	if result != "hello" {
		t.Error("You should have had me at hello")
	}
}

func TestMakeRunEchoFilter(t *testing.T) {
	fun, err := makeFilter(`module.exports = function(input, cxt) { cxt.done(input) }`)
	if err != nil {
		t.Error(err)
	}
	result, err := fun("echo")
	if result != "echo" {
		t.Error("Filter did not echo input.")
	}
}

func TestMakeRunFailsInvalidJSFilter(t *testing.T) {
	/* Note: this is not valid ECMAScript */
	fun, err := makeFilter("(╯°□°）╯︵ ┻━┻")
	_, err = fun("")
	if err == nil {
		t.Error("One should not (╯°□°）╯︵ ┻━┻")
	}
}
