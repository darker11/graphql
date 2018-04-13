package loghelper

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	log "gitlab.ucloudadmin.com/wu/logrus"
)

var logFilePath string

func registerRotate() {
	c := make(chan os.Signal)
	for {
		signal.Notify(c, syscall.SIGHUP)
		s := <-c
		if s == syscall.SIGHUP {
			reloadLogFile()
		}
	}
}

func reloadLogFile() {
	folder := path.Dir(logFilePath)
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.WithError(err).Fatal("Fail to mk folder")
	}
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.WithError(err).Fatal("Fail to open log file")
	}
	log.SetOutput(f)
}

func initLogger() {
	log.Info("init logger now!!!")
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fail to find current working directory %v", err)
	}
	logFileName := os.Getenv("UMQ_LOG_FILE_NAME")
	if len(logFileName) == 0 {
		logFileName = "app"
	}
	logFilePath = path.Join(dir, "logs", logFileName+".log")
	if os.Getenv("U_DEPLOY_STAGE") == "production" {
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
		})
	}
	reloadLogFile()
	go registerRotate()
}

func registerControl() {
	http.HandleFunc("/control/loglevel", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "PUT" {
			err := req.ParseForm()
			if err != nil {
				http.Error(w, "only accept form data", 400)
				return
			}
			v := req.PostForm.Get("level")
			if v == "" {
				http.Error(w, "missing level field", 400)
				return
			}
			err = setLevel(v)
			if err != nil {
				http.Error(w, err.Error(), 400)
			}
			vv := getLevelDesc()
			w.Header().Set("Content-type", "plain/text")
			fmt.Fprintf(w, vv)
		} else if req.Method == "GET" {
			w.Header().Set("Content-type", "plain/text")
			fmt.Fprintf(w, getLevelDesc())
		}
	})
}

func getLevelDesc() string {
	l := log.GetLevel()
	switch l {
	case log.DebugLevel:
		return "debug"
	case log.InfoLevel:
		return "info"
	case log.WarnLevel:
		return "warn"
	case log.ErrorLevel:
		return "error"
	default:
		return "unknown"
	}
}

func setLevel(level string) error {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		return errors.New("Unknown level")
	}
	return nil
}
