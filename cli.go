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

		Exec   - fetches an object and creates/updates a
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

// Handle the 'fetch' CLI command.
func cmdFetch(c *cli.Context) {
	logrus.Debug("Fetching ", c.Args().First())
	fromPath := dereferencePath(c.Args().First())
	fromObj := dereferenceObj(fromPath)

	msg := fromObj.read()
	logrus.Debug("Raw msg:")
	println(msg)
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
		var msg string

		path, err := url.Parse(arg)
		if err != nil {
			log.Fatal(err)
		}

		logrus.Info("pipe[arg]: " + arg)

		if path.Scheme == "http" || path.Scheme == "https" {
			argPath := dereferencePath(arg)
			argObj := dereferenceObj(argPath)

			// If first argument, then we must GET,
			// note that this case follows the '-' so all
			// shell input will pipe into the POST.
			if i == 0 {
				msg = argObj.read()
			} else {
				msg = argObj.update(lastObj)
			}
		} else if arg == "-" {
			script, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
				return
			}
			if i == 0 {
				// If first argument, assume is generic bytes input.
				msg = string(script[:])
			} else {
				// If not the first argument, then expect pipescript
				filter, err := makeFilter(string(script[:]))
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
			msg, err = execFilter(lastObj, arg)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		logrus.Debug(fmt.Sprintf("pipe[%i][raw]: %s\n", i, msg))

		if i == len(c.Args()) {
			println(msg)
			return
		}
		lastObj = msg
	}
}

// Handle the 'create' CLI command.
func cmdCreate(c *cli.Context) {
	logrus.Debug("Creating object ", c.Args().First())
}
