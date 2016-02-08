package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"log"
	"os"
	"path"

	"encoding/json"
)

var debug bool = false
var name string

func main() {
	ensureCachePath()

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
		/*{
			Name:   "create",
			Usage:  "Create pipescript in local index.",
			Action: cmdCreate,
		},*/
		{
			Name:   "remove",
			Usage:  "Remove filter from local index.",
			Action: cmdRemove,
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
		/*{
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
		},*/
		{
			Name:   "tag",
			Usage:  "Tag a pipescript or pipeline to a name",
			Action: cmdTag,
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

func cmdRemove(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("No filters specified.")
	}
	filter := c.Args()[0]

	err := removeFilter(filter)
	if err != nil {
		log.Fatal(err)
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
	if name == "" {
                log.Fatal("Name not specified. Must pass required --name flag.")
	}
	exportScript(name, c.Args()...)
}

func cmdTag(c *cli.Context) {
	var err error

	id := c.Args()[0]
	name := c.Args()[1]
	if name == "" {
		log.Fatal("Invalid name")
	}
	if err = tagFilter(id, name); err != nil {
		log.Fatal(err)
	}
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
	var msg string
	var err error
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if msg, err = Exec(c.Args()...); err != nil {
		log.Fatal(err)
		return
	}
	println(msg)
}
