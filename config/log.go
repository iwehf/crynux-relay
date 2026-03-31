package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

const (
	defaultNodeHealthLogPath     = "data/logs/node_health.log"
	defaultTaskAssignmentLogPath = "data/logs/task_assignment.log"
)

var nodeHealthLogger *logrus.Logger
var taskAssignmentLogger *logrus.Logger

func InitLog(appConfig *AppConfig) error {

	println("Initializing logger...")

	logrus.SetFormatter(&logrus.TextFormatter{})

	switch appConfig.Log.Output {
	case "", "stderr":
		logrus.SetOutput(os.Stderr)
	case "stdout":
		logrus.SetOutput(os.Stdout)
	default:
		logWriter := newLogWriter(appConfig.Log.Output, appConfig.Log.MaxFileSize, appConfig.Log.MaxDays, appConfig.Log.MaxFileNum)
		mw := io.MultiWriter(os.Stdout, logWriter)
		logrus.SetOutput(mw)
	}

	level, err := logrus.ParseLevel(appConfig.Log.Level)

	if err != nil {
		return err
	}

	logrus.SetLevel(level)
	initNodeHealthLogger(appConfig)
	initTaskAssignmentLogger(appConfig)

	return nil
}

func GetNodeHealthLogger() *logrus.Logger {
	return nodeHealthLogger
}

func GetTaskAssignmentLogger() *logrus.Logger {
	return taskAssignmentLogger
}

func initNodeHealthLogger(appConfig *AppConfig) {
	if !isNodeHealthLogEnabled(appConfig) {
		nodeHealthLogger = nil
		return
	}
	nodeHealthLogger = logrus.New()
	nodeHealthLogger.SetFormatter(&logrus.TextFormatter{})
	nodeHealthLogger.SetLevel(logrus.InfoLevel)
	nodeHealthLogger.SetOutput(newLogWriter(getNodeHealthLogPath(appConfig.Log.Output), appConfig.Log.MaxFileSize, appConfig.Log.MaxDays, appConfig.Log.MaxFileNum))
}

func initTaskAssignmentLogger(appConfig *AppConfig) {
	if !appConfig.Log.Features.TaskAssignmentEnabled {
		taskAssignmentLogger = nil
		return
	}
	taskAssignmentLogger = logrus.New()
	taskAssignmentLogger.SetFormatter(&logrus.TextFormatter{})
	taskAssignmentLogger.SetLevel(logrus.InfoLevel)
	taskAssignmentLogger.SetOutput(newLogWriter(getTaskAssignmentLogPath(appConfig.Log.Output), appConfig.Log.MaxFileSize, appConfig.Log.MaxDays, appConfig.Log.MaxFileNum))
}

func isNodeHealthLogEnabled(appConfig *AppConfig) bool {
	if appConfig.Log.Features.NodeHealthEnabled == nil {
		return false
	}
	return *appConfig.Log.Features.NodeHealthEnabled
}

func getNodeHealthLogPath(mainLogOutput string) string {
	if mainLogOutput == "" || mainLogOutput == "stdout" || mainLogOutput == "stderr" {
		return defaultNodeHealthLogPath
	}
	return filepath.Join(filepath.Dir(mainLogOutput), "node_health.log")
}

func getTaskAssignmentLogPath(mainLogOutput string) string {
	if mainLogOutput == "" || mainLogOutput == "stdout" || mainLogOutput == "stderr" {
		return defaultTaskAssignmentLogPath
	}
	return filepath.Join(filepath.Dir(mainLogOutput), "task_assignment.log")
}

func newLogWriter(filename string, maxFileSize, maxDays, maxFileNum int) *lumberjack.Logger {
	logWriter := &lumberjack.Logger{
		Filename: filename,
		Compress: true,
	}

	if maxFileSize == 0 {
		logWriter.MaxSize = 500
	} else {
		logWriter.MaxSize = maxFileSize
	}

	if maxDays == 0 {
		logWriter.MaxAge = 30
	} else {
		logWriter.MaxAge = maxDays
	}

	if maxFileNum == 0 {
		logWriter.MaxBackups = 10
	} else {
		logWriter.MaxBackups = maxFileNum
	}

	return logWriter
}
