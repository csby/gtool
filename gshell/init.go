package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	moduleName    = "gshell"
	moduleVersion = "1.0.1.2"
)

var (
	svcDir = ""
	server = &program{}
	log    = &LogWriter{}
)

// args['exec', 'svc-folder', 'log-folder']
func init() {
	args := os.Args
	path, _ := filepath.Abs(args[0])

	argc := len(args)
	for i := 1; i < argc; i++ {
		arg := args[i]
		if strings.ToLower(arg) == "-v" {
			ver := &struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			}{
				moduleName,
				moduleVersion,
			}
			vd, ve := json.Marshal(ver)
			if ve == nil {
				if arg == "-V" {
					fmt.Println(base64.StdEncoding.EncodeToString(vd))
				} else {
					fmt.Println(string(vd))
				}

			}
			os.Exit(0)
		}
	}

	if argc > 1 {
		svcDir = args[1]
	} else {
		svcDir = filepath.Dir(path)
	}
	if argc > 2 {
		log.folder = args[2]
	}

	if len(svcDir) > 0 {
		os.Chdir(svcDir)
	}
	curDir, _ := os.Getwd()

	log.Info("shell run at: ", path)
	log.Info("shell version: ", moduleVersion)
	log.Info("log folder: ", log.folder)
	log.Info("cur folder: ", curDir)
	log.Info("svc folder: ", svcDir)

	server.shell.Directory = svcDir
}
