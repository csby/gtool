package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var stdWriter io.Writer = os.Stdout

type LogWriter struct {
	sync.Mutex

	folder  string
	year    int
	month   time.Month
	day     int
	file    *os.File
	levels  []string
	curSize int
	maxSize int
}

func (s *LogWriter) Write(p []byte) (n int, err error) {
	if len(p) < 1 {
		return 0, nil
	}

	if len(s.levels) > 0 {
		msg := strings.ToUpper(string(p))
		if !s.isLevelMatch(msg) {
			return 0, nil
		}
	}

	writer := s.getLogger()
	if writer != nil {
		s.Lock()
		defer s.Unlock()

		n, err = writer.Write(p)
		if err == nil {
			s.curSize += n
		}

		if s.maxSize > 0 {
			if s.curSize > s.maxSize {
				s.curSize = 0
				s.year = 0
				s.month = 0
				s.day = 0
			}
		}

		return
	}

	return 0, nil
}

func (s *LogWriter) Close() error {
	s.Lock()
	defer s.Unlock()

	if s.file != nil {
		err := s.file.Close()
		s.file = nil
		return err
	}

	return nil
}

func (s *LogWriter) Info(v ...interface{}) string {
	return s.output(" INFO", fmt.Sprint(v...))
}

func (s *LogWriter) Error(v ...interface{}) string {
	return s.output("ERROR", fmt.Sprint(v...))
}

func (s *LogWriter) Output(v ...interface{}) string {
	msg := fmt.Sprint(v...)
	s.Write([]byte(fmt.Sprintln(msg)))
	return msg
}

func (s *LogWriter) output(l string, m string) string {
	now := time.Now()
	str := fmt.Sprintf("%s %d --- %s", l, os.Getpid(), m)
	msg := fmt.Sprintf("%s %s", now.Format("2006-01-02 15:04:05.000"), str)

	s.Write([]byte(fmt.Sprintln(msg)))

	return str
}

func (s *LogWriter) getLogger() io.Writer {
	if s.folder != "" {
		now := time.Now()
		if s.year != now.Year() || s.month != now.Month() || s.day != now.Day() || s.file == nil {
			s.Lock()
			defer s.Unlock()

			s.year = now.Year()
			s.month = now.Month()
			s.day = now.Day()
			if s.file != nil {
				s.file.Close()
				s.file = nil
			}

			os.MkdirAll(s.folder, 0777)
			filePath := s.getValidFilePath(s.folder, now)
			file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return nil
			}
			s.file = file
		}
	}

	if s.file != nil {
		return s.file
	}

	return stdWriter
}

func (s *LogWriter) getValidFilePath(folder string, t time.Time) string {
	fileName := fmt.Sprintf("%s.log", t.Format("2006-01-02"))
	filePath := filepath.Join(folder, fileName)

	index := 1
	for s.isExist(filePath) {
		fileName = fmt.Sprintf("%s_%d.log", t.Format("2006-01-02"), index)
		filePath = filepath.Join(folder, fileName)
		index++
	}

	return filePath
}

func (s *LogWriter) isExist(name string) bool {
	_, err := os.Stat(name)

	return err == nil || os.IsExist(err)
}

func (s *LogWriter) isLevelMatch(msg string) bool {
	c := len(s.levels)
	if c < 1 {
		return true
	}

	for i := 0; i < c; i++ {
		level := s.levels[i]
		if strings.Index(msg, level) > -1 {
			return true
		}
	}

	return false
}
