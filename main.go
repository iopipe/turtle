package main

import (
	"github.com/codegangsta/cli"
	"github.com/robertkrimen/otto"

	"errors"
	"log"
	"os"
	"strings"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	app := cli.NewApp()
	app.Name = "iopipe"
	app.Usage = "cross-API interoperability & data manager"
	app.Action = func(c *cli.Context) {
		println("object object")
	}

	/*******************************************************
	Commands:

		Copy   - fetches an object and creates/updates a
				 resource at destination
		Fetch  - a 'copy' to STDOUT (i.e. read resource)
		Delete - deletes an object
		Create - Like copy, but will not update existing
				 records (possible flag to 'copy' instead?)
		Update - Like copy, but refuse to create new records
				 (possible flag to copy instead?)
	*******************************************************/
	app.Commands = []cli.Command{
		{
			Name:   "copy",
			Usage:  "Copy or pipe from <src> to <dest>",
			Action: cmdCopy,
		},
		{
			Name:   "convert",
			Usage:  "Pipe from <src> to stdout",
			Action: cmdConvert,
		},
		{
			Name:   "fetch",
			Usage:  "Fetch <src>, output to STDOUT",
			Action: cmdFetch,
		},
		{
			Name:  "delete",
			Usage: "Delete object",
			Action: func(c *cli.Context) {
				println("Deleting ", c.Args().First())
			},
		},
		{
			Name:   "create",
			Usage:  "Create object. Like copy, but only if can be guaranteed as a new object.",
			Action: cmdCreate,
		},
		{
			Name:  "update",
			Usage: "Update an object, only if it already exists.",
			Action: func(c *cli.Context) {
				println("Creating ", c.Args().First())
			},
		},
		{
			Name:   "exec",
			Usage:  "Perform a registered Action.",
			Action: cmdExec,
		},
	}
	app.Run(os.Args)
}

/*******************************************************
	High level API management:

	APIs must expose:
		* Objects
		* Actions

	Objects:

		These are objects which may be part of a
		CRUD operation (to create or modify objects),
		or as input or output for Actions.

	Actions:

		Actions are functions to perform a task,
		accepting and outputting Objects.
*******************************************************/

// Handle the 'copy' CLI command.
func cmdCopy(c *cli.Context) {
	println("Fetching ", c.Args().First())
	fromPath := dereferencePath(c.Args().First())
	fromObj := dereferenceObj(fromPath)

	//obj := from.read()
	//dest.update(obj)
	msg := fromObj.read()

	println("Sending to ", c.Args().Get(1))
	println("Content: ", msg)

	destPath := dereferencePath(c.Args().Get(1))
	destObj := dereferenceObj(destPath)

	response := destObj.update(msg)

	println("Recipient response: ", response)
}

// Handle the 'fetch' CLI command.
func cmdFetch(c *cli.Context) {
	println("Fetching ", c.Args().First())
	fromPath := dereferencePath(c.Args().First())
	fromObj := dereferenceObj(fromPath)

	msg := fromObj.read()
	println("Raw msg:")
	println(msg)
}

type objectInterface struct {
	ClassID    string                 `json:"classid"`
	Properties map[string]interface{} `json:"properties"`
}

// Handle the 'fetch' CLI command.
func cmdConvert(c *cli.Context) {
	var err error

	println("Fetching ", c.Args().First())
	fromPath := dereferencePath(c.Args().First())
	fromObj := dereferenceObj(fromPath)

	msg := fromObj.read()
	println("Raw msg:")
	println(msg)

	toObjType := c.Args().Get(1)
	//toObjType := "com.iopipe.messaging.GenericMessage"
	println("Converting to: " + toObjType)

	var obj objectInterface
	if err = json.Unmarshal([]byte(msg), &obj); err != nil {
		log.Fatal(err)
		return
	}
	fromObjType := obj.ClassID
	println("Converting from: " + fromObjType)
	filter, err := findFilter(fromObjType, toObjType)
	if err != nil {
		log.Fatal(err)
		return
	}
	//filter(msg)
	resp, err := filter(msg)
	if err != nil {
		log.Fatal(err)
		return
	}
	println("Conversion: " + resp)
}

type filterTuple struct {
	fromObjType string
	toObjType   string
}

func findFilter(fromObjType string, toObjType string) (func(input string) (string, error), error) {
	var script string

	fT := filterTuple{fromObjType: fromObjType, toObjType: toObjType}

	//var m map[filterTuple]func(obj objectInterface) (string, error)
	/* Filter definitions */
	switch fT {
	case filterTuple{
		fromObjType: "com.twitter.statusMessage",
		toObjType:   "com.iopipe.messaging.GenericMessage",
	}:
		script = `
			var obj = JSON.parse(input);
			var tweet = obj["properties"];
			var statusMessage = {
				"id":   "/objects/statusMessage/" + tweet["id_str"],
				"user": "/objects/user/" + tweet["user"]["id_str"],
				"text": tweet["text"]
			};
			JSON.stringify(statusMessage);
		`
	case filterTuple{
		fromObjType: "com.twiter.statusRequest",
		toObjType:   "com.iopipe.messaging.GenericMessage",
	}:
		script = `
			var obj = JSON.parse(input);
			var statusRequest = {
				"status": obj["properties"]["text"]
			};
			JSON.stringify(statusRequest);
		`
	default:
		return nil, errors.New("No filter found.")
	}

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

// Handle the 'exec' CLI commmand.
func cmdExec(c *cli.Context) {
	println("Executing action ", c.Args().First())
}

// Handle the 'create' CLI command.
func cmdCreate(c *cli.Context) {
	println("Creating object ", c.Args().First())
}

/*******************************************************
 Object Mapper
*******************************************************/
type MetaObject struct {
	objtype []string
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

type ObjPath struct {
	host   string
	scheme string
	uri    string
	query  []string
}

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
	/*rawObj := json.Decode... into MetaObject
	return rawObj*/
	obj := new(Object)
	obj.path = path
	return obj
}
