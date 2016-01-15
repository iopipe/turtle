package main

import (
	"github.com/Sirupsen/logrus"

	"fmt"
	"os"

	"io/ioutil"
	"net/url"
)

func Exec(args ...string) (string, error) {
	//var err error
	var lastObj string
	var msg string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		path, err := url.Parse(arg)
		if err != nil {
			return "", err
		}

		logrus.Info("pipe[arg]: " + arg)

		if path.Scheme == "http" || path.Scheme == "https" {
			argPath, err := dereferencePath(arg)
			if err != nil {
				return "", err
			}
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
				return "", err
			}
			if i == 0 {
				// If first argument, assume is generic bytes input.
				msg = string(script[:])
			} else {
				// If not the first argument, then expect pipescript
				filter, err := makeFilter(string(script[:]))
				if err != nil {
					return "", err
				}
				if msg, err = filter(lastObj); err != nil {
					return "", err
				}
			}
		} else {
			filter, err := getFilter(arg)
			if err != nil {
				return "", err
			}
			// Search
			if i == 0 {
				msg, err = filter("")
			} else {
				msg, err = filter(lastObj)
			}
			if err != nil {
				return "", err
			}
		}
		logrus.Debug(fmt.Sprintf("pipe[%i][raw]: %s\n", i, msg))

		lastObj = msg
	}
	return msg, nil
}
