package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"

	/* object mapper - path dereference */
	//"encoding/json"

	/* object mapper */
	"net/http"

	"strings"

	"io/ioutil"
)

func main() {
	app := cli.NewApp()
	app.Name = "iopipe"
	app.Usage = "API cat - the API rosetta stone"
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
	println(msg)
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
	path *ObjPath
}

func (object *Object) read() string {
	url := object.path.canonicial()
	res, err := http.Get(url)
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
	url := object.path.canonicial()
	reader := strings.NewReader(content)
	req, err := http.NewRequest("PUT", url, reader)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
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

type ObjPath struct {
	host   string
	scheme string
	uri    string
	query  []string
}

/*
Create a canonicial URL from ObjPath
*/
func (path *ObjPath) canonicial() string {
	//return path.scheme + "://" + path.host + "/" + path.uri
	return "https://" + path.host + "/" + path.uri
}

/*
Create an ObjPath from a string
*/
func dereferencePath(reqPath string) *ObjPath {
	path := new(ObjPath)
	parts := strings.SplitN(reqPath, "/", 2)
	path.host = parts[0]
	path.scheme = "https"
	path.uri = parts[1]
	return path
}

/*
Download resource at path &
validate resource matches declared object type definition.
*/
func dereferenceObj(path *ObjPath) *Object {
	/*rawObj := json.Decode... into MetaObject
	return rawObj*/
	obj := new(Object)
	obj.path = path
	return obj
}
