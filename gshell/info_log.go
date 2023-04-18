package main

type InfoLog struct {
	Level   string `json:"level" note:"日志等级,如ERROR|WARN|INFO|DEBUG，空(默认)表示全部"`
	MaxSize int    `json:"maxSize" note:"文件大小最大值,单位字节, 0(默认)表示不限制"`
}
