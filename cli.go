package main

import (
	"github.com/codegangsta/cli"

	"fmt"
	"log"
	"os"

	"encoding/json"
	"io/ioutil"
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
			Name:   "exec",
			Usage:  "Pipe from <src> to stdout",
			Action: cmdExec,
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

// Handle the 'exec' CLI command.
func cmdExec(c *cli.Context) {
	//var err error
	var lastObj string

	for i := 0; i < len(c.Args()); i++ {
		apart := c.Args()[i]
		println("Deconstructing ", apart)
		var msg string

		path, err := url.Parse(apart)
		if err != nil {
			log.Fatal(err)
		}

		if path.Scheme == "http" || path.Scheme == "https" {
			println("From HTTP")
			argPath := dereferencePath(apart)
			argObj := dereferenceObj(argPath)
			// If first argument, then we must GET,
			// note that this case follows the '-' so all
			// shell input will pipe into the POST.
			if (i == 0) {
				msg = argObj.read()
			} else {
				msg = argObj.update(lastObj)
			}
		} else if apart == "-" {
			println("From STDIN")
			bytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
			msg = string(bytes[:])
		} else {
			println("via default")
			toObjType := apart

			var obj objectInterface
			if err = json.Unmarshal([]byte(lastObj), &obj); err != nil {
				log.Fatal(err)
				return
			}
			lastObjType := obj.ClassID

			fmt.Printf("pipe[%i][lastObjType/toObjType]: %s/%s\n", i, lastObjType, toObjType)

			filter, err := findFilter(lastObjType, toObjType)
			if err != nil {
				log.Fatal(err)
				return
			}
			msg, err = filter(lastObj)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		fmt.Printf("pipe[%i][raw]: %s\n", i, msg)

		if i == len(c.Args()) {
			println("output: " + msg)
			return
		}
		lastObj = msg
	}
}

// Handle the 'create' CLI command.
func cmdCreate(c *cli.Context) {
	println("Creating object ", c.Args().First())
}
