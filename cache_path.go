package main

import (
	"github.com/Sirupsen/logrus"

	"fmt"
	"os"
	"path"

	"crypto/sha256"
	"io/ioutil"
	"os/user"
)

func getCachePath(name string) (string, error) {
	myuser, err := user.Current()
	if err != nil {
		return "", err
	}
	pathParts := []string{myuser.HomeDir, ".iopipe", "filter_cache", name}
	return path.Join(pathParts...), nil
}

func ensureCachePath() error {
	path, err := getCachePath("")
	if err != nil {
		return err
	}
	return os.MkdirAll(path, 0700)
}

func readFilterCache(name string) ([]byte, error) {
	var err error

	diskPath, err := getCachePath(name)
	if err != nil {
		return nil, err
	}

	/* Do we have this cached? */
	if _, err = os.Stat(diskPath); err != nil {
		return nil, err
	}
	script, err := ioutil.ReadFile(diskPath)

	logrus.Debug("Read filter from cache:\n" + string(script[:]))
	return script[:], nil
}

func writeFilterCache(body []byte) (string, error) {
	var err error

	/* Verify digest */
	chksum := sha256.Sum256(body[:])
	id := fmt.Sprintf("%x", chksum)
	diskPath, err := getCachePath(id)
	if err != nil {
		return id, err
	}

	/* Write cache */
	if err = ioutil.WriteFile(diskPath, body, 0600); err != nil {
		return id, err
	}
	return id, nil
}
