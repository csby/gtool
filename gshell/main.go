package main

import (
	"fmt"
	"github.com/kardianos/service"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defer log.Close()
	defer log.Info("shell exit", fmt.Sprintln(""))

	info := &Info{}
	err := info.LoadFromFile(filepath.Join(svcDir, "info.json"))
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	if len(info.Exec) < 1 {
		log.Error("可执行程序为空(info.json->exec='')")
		os.Exit(-1)
	}
	log.Info("svc exec: ", info.Exec)
	log.Info("svc args: ", info.Args)
	server.shell.Exec = info.Exec
	server.shell.Args = info.Args
	server.shell.Prepares = info.Prepares
	c := len(info.Prepares)
	for i := 0; i < c; i++ {
		p := info.Prepares[i]
		if p == nil {
			continue
		}
		log.Info(fmt.Sprintf("prepare %d exec: ", i+1), p.Exec)
		log.Info(fmt.Sprintf("prepare %d args: ", i+1), p.Args)
	}
	log.Info("log level: ", info.Log.Level)
	log.Info("log max size: ", info.Log.MaxSize)

	if service.Interactive() == false {
		if len(info.Name) < 1 {
			log.Error("服务名称为空(info.json->name='')")
			os.Exit(-11)
		}
		svcName := info.ServiceName()

		cfg := &service.Config{
			Name:        svcName,
			DisplayName: svcName,
		}
		svc, e := service.New(server, cfg)
		if e != nil {
			log.Error("创建服务失败: ", e)
			os.Exit(-12)
		}
		server.service = svc
	}

	log.maxSize = info.Log.MaxSize
	log.levels = make([]string, 0)
	levels := strings.Split(strings.ToUpper(info.Log.Level), "|")
	c = len(levels)
	for i := 0; i < c; i++ {
		level := strings.TrimSpace(levels[i])
		if len(level) > 0 {
			log.levels = append(log.levels, level)
		}
	}

	server.Run()
}
