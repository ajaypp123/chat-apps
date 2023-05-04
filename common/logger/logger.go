package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/ajaypp123/chat-apps/common/appcontext"
)

type Level int

const (
	INFO Level = iota
	ERROR
	WARN
	DEBUG
)

func (l Level) String() string {
	switch l {
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case DEBUG:
		return "DEBUG"
	case ERROR:
		return "ERROR"
	default:
		return ""
	}
}

type Logger struct {
	level    Level
	filename string
	file     *os.File
}

var logMapper map[string]*Logger

func NewLogger(ctx *appcontext.AppContext, filename string, level Level) error {
	index := ctx.GetValue("index").(string)
	if _, ok := logMapper[index]; ok {
		return nil
	}

	filename, err := getLogFilepath(filename)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	lg := &Logger{
		level:    level,
		filename: filename,
		file:     file,
	}

	if logMapper == nil {
		logMapper = make(map[string]*Logger)
	}

	logMapper[index] = lg
	return nil
}

func logf(ctx *appcontext.AppContext, filename, line string, level Level, message string) {
	index, _ := ctx.GetValue("index").(string)
	l := logMapper[index]
	fmt.Println(l)

	if l.level < level {
		return
	}

	now := time.Now()
	timeFormat := now.Format("2006/01/02 15:04:05")
	timestamp := now.Format(time.RFC3339)
	logLine := fmt.Sprintf("%s [%s] %s:%s:%s %s: %s", timeFormat, timestamp, index, filename, line, level, message)

	select {
	case <-ctx.Done():
		return
	default:
	}

	// convert to json
	log.Println(logLine)
	if _, err := l.file.WriteString(logLine + "\n"); err != nil {
		log.Println(err)
	}
}

func Info(ctx *appcontext.AppContext, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	logf(ctx, filename, strconv.Itoa(line), INFO, fmt.Sprint(args...))
}

func Warn(ctx *appcontext.AppContext, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	logf(ctx, filename, strconv.Itoa(line), WARN, fmt.Sprint(args...))
}

func Debug(ctx *appcontext.AppContext, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	logf(ctx, filename, strconv.Itoa(line), DEBUG, fmt.Sprint(args...))
}

func Error(ctx *appcontext.AppContext, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	logf(ctx, filename, strconv.Itoa(line), ERROR, fmt.Sprint(args...))
}

func Close(ctx *appcontext.AppContext) {
	index, _ := ctx.GetValue("index").(string)
	l := logMapper[index]
	l.file.Close()
}

func getLogFilepath(filename string) (string, error) {
	// Get the absolute path of the executable directory
	path, err := os.Getwd() // filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	// Create the logs directory if it doesn't exist
	logsDir := filepath.Join(path, "logs")
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err := os.Mkdir(logsDir, 0755); err != nil {
			return "", err
		}
	}

	// Create the log file if it doesn't exist
	logFilepath := filepath.Join(logsDir, filename)
	if _, err := os.Stat(logFilepath); os.IsNotExist(err) {
		if _, err := os.Create(logFilepath); err != nil {
			return "", err
		}
	}

	return logFilepath, nil
}
