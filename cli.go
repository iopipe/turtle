package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"fmt"
	"log"
	"os"
	"path"

	"encoding/json"
	"io/ioutil"
	"net/url"
)

var debug bool = false
var name string

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
	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Create pipescript in local index.",
			Action: cmdCreate,
		},
		{
			Name:  "remove",
			Usage: "Remove pipescript from local index.",
			Action: func(c *cli.Context) {
				logrus.Debug("Deleting ", c.Args().First())
			},
		},
		{
			Name:   "exec",
			Usage:  "Pipe from <src> to stdout",
			Action: cmdExec,
		},
		{
			Name:   "export",
			Usage:  "Export will write a Javascript library for the given pipeline",
			Action: cmdExport,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "name",
					Usage:       "NPM package name",
					Destination: &name,
				},
			},
		},
		{
			Name:   "import",
			Usage:  "Import will bring a filter in from a javascript file or STDIN for -",
			Action: cmdImport,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "name",
					Usage:       "Tag to name",
					Destination: &name,
				},
			},
		},
		{
			Name:   "list",
			Usage:  "List local and subscribed pipes",
			Action: cmdList,
		},
		{
			Name:  "login",
			Usage: "Login to an endpoint",
			Action: func(c *cli.Context) {
				logrus.Debug("Logging in to ", c.Args().First())
			},
		},
		{
			Name:  "publish",
			Usage: "Publish a pipescript or pipeline",
			Action: func(c *cli.Context) {
				logrus.Debug("Pushing ", c.Args().First())
			},
		},
		{
			Name:  "subscribe",
			Usage: "Add pipescript or pipeline to local index",
			Action: func(c *cli.Context) {
				logrus.Debug("Pulling ", c.Args().First())
			},
		},
	}
	app.Run(os.Args)
}

func cmdList(c *cli.Context) {
	list, err := listScripts()
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range list {
		println(name)
	}
}

func cmdLogin(c *cli.Context) {
}

func cmdPublish(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("No pipeline specified for subscription.")
	}
	pipeline := c.Args()[0]

	publishPipeline(pipeline)
}

func cmdSubscribe(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("No pipeline specified for subscription.")
	}
	pipeline := c.Args()[0]

	subscribePipeline(pipeline)
}

func cmdCreate(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("No filters specified.")
	}
	pipeline := c.Args()

	plid, err := createPipeline(pipeline)
	if err != nil {
		log.Fatal(err)
	}

	if name != "" {
		err = tagPipeline(plid)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func cmdImport(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("No pipeline specified for import.")
	}
	file := c.Args()[0]
	id, err := importScript(file)
	if err != nil {
		log.Fatal(err)
	}
	if name != "" {
		if err = tagFilter(id, name); err != nil {
			log.Fatal(err)
		}
	}
	println(id)
}

func cmdExport(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("No pipeline specified for export.")
	}
	pipeline := c.Args()[0]
	if name == "" {
		name = pipeline
	}

	exportScript(pipeline, name)
}

func execFilter(lastObj string, toObjType string) (msg string, err error) {
	var obj objectInterface
	if err = json.Unmarshal([]byte(lastObj), &obj); err != nil {
		return "", err
	}
	lastObjType := obj.ClassID

	logrus.Info("pipe[lastObjType/toObjType]: %s/%s\n", lastObjType, toObjType)

	filter, err := getFilter(path.Join(lastObjType, toObjType))
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
			// Search
			/*msg, err = execFilter(lastObj, arg)
			if err != nil {
				log.Fatal(err)
				return
			}*/
			filter, err := getFilter(arg)
			if err != nil {
				log.Fatal("Filter not found.")
				return
			}
			// Search
			if i == 0 {
				msg, err = filter("")
			} else {
				msg, err = filter(lastObj)
			}
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		logrus.Debug(fmt.Sprintf("pipe[%i][raw]: %s\n", i, msg))

		if i == len(c.Args()) {
			println(msg)
			log.Print(msg)
			return
		}
		lastObj = msg
	}
}
