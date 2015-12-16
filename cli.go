package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"fmt"
	"log"
	"os"

	"encoding/json"
	"io/ioutil"
	"net/url"
)

var debug bool = false

func main() {
	//var debug bool = false
	app := cli.NewApp()
	app.Name = "iopipe"
	app.Usage = "cross-API interoperability & data manager"
	app.Action = func(c *cli.Context) {
		logrus.Debug("object object")
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "enable debugging output",
			Destination: &debug,
		},
	}
	app.Action = func(c *cli.Context) {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
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
				logrus.Debug("Deleting ", c.Args().First())
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
				logrus.Debug("Creating ", c.Args().First())
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
	logrus.Debug("Fetching ", c.Args().First())
	fromPath := dereferencePath(c.Args().First())
	fromObj := dereferenceObj(fromPath)

	msg := fromObj.read()

	logrus.Debug("Sending to ", c.Args().Get(1))
	logrus.Debug("Content: ", msg)

	destPath := dereferencePath(c.Args().Get(1))
	destObj := dereferenceObj(destPath)

	response := destObj.update(msg)

	logrus.Debug("Recipient response: ", response)
}

// Handle the 'fetch' CLI command.
func cmdFetch(c *cli.Context) {
	logrus.Debug("Fetching ", c.Args().First())
	fromPath := dereferencePath(c.Args().First())
	fromObj := dereferenceObj(fromPath)

	msg := fromObj.read()
	logrus.Debug("Raw msg:")
	logrus.Debug(msg)
}

func execFilter(lastObj string, toObjType string) (msg string, err error) {
	var obj objectInterface
	if err = json.Unmarshal([]byte(lastObj), &obj); err != nil {
		return "", err
	}
	lastObjType := obj.ClassID

	logrus.Info("pipe[lastObjType/toObjType]: %s/%s\n", lastObjType, toObjType)

	filter, err := findFilter(lastObjType, toObjType)
	if err != nil {
		return "", err
	}
	return filter(lastObj)
}

// Handle the 'exec' CLI command.
func cmdExec(c *cli.Context) {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	//var err error
	var lastObj string

	for i := 0; i < len(c.Args()); i++ {
		arg := c.Args()[i]
		logrus.Debug("Deconstructing ", arg)
		var msg string

		path, err := url.Parse(arg)
		if err != nil {
			log.Fatal(err)
		}

		if path.Scheme == "http" || path.Scheme == "https" {
			argPath := dereferencePath(arg)
			argObj := dereferenceObj(argPath)

			logrus.Info("pipe[argPath]: " + argPath)
			// If first argument, then we must GET,
			// note that this case follows the '-' so all
			// shell input will pipe into the POST.
			if i == 0 {
				msg = argObj.read()
			} else {
				msg = argObj.update(lastObj)
			}
		} else if arg == "-" {
			logrus.Debug("From STDIN")
			bytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
				return
			}
			if i == 0 {
				msg = string(bytes[:])
			} else {
				// If not the first argument, then expect pipescript
				filter, err := makeFilter(string(bytes[:]))
				if err != nil {
					log.Fatal(err)
					return
				}
				if msg, err = filter(lastObj); err != nil {
					log.Fatal(err)
					return
				}
			}
		} else {
			logrus.Debug("via default")
			msg, err = execFilter(lastObj, arg)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		logrus.Debug(fmt.Sprintf("pipe[%i][raw]: %s\n", i, msg))

		if i == len(c.Args()) {
			logrus.Debug("output: " + msg)
			return
		}
		lastObj = msg
	}
}

// Handle the 'create' CLI command.
func cmdCreate(c *cli.Context) {
	logrus.Debug("Creating object ", c.Args().First())
}
