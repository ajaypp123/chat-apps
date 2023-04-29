package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
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

var log_mapper map[string]*Logger

func NewLogger(ctx context.Context, filename string, level Level) error {
	filename, err := getLogFilepath(filename)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	lg := &Logger{
		level:    level,
		filename: filename,
		file:     file,
	}

	if log_mapper == nil {
		log_mapper = make(map[string]*Logger)
	}

	index, _ := ctx.Value("index").(string)
	log_mapper[index] = lg
	return nil
}

func logf(ctx context.Context, filename, line string, level Level, format string, args ...interface{}) {
	index, _ := ctx.Value("index").(string)
	l := log_mapper[index]

	if l.level < level {
		return
	}

	now := time.Now()
	timeFormat := now.Format("2006/01/02 15:04:05")
	timestamp := now.Format(time.RFC3339)
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("%s [%s] %s:%s:%s %s: %s", timeFormat, timestamp, index, filename, line, level, message)

	select {
	case <-ctx.Done():
		return
	default:
	}

	log.Println(logLine) // print on console
	l.file.WriteString(logLine + "\n")
}

func Info(ctx context.Context, format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	logf(ctx, filename, strconv.Itoa(line), INFO, format, args...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	logf(ctx, filename, strconv.Itoa(line), WARN, format, args...)
}

func Debug(ctx context.Context, format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	logf(ctx, filename, strconv.Itoa(line), DEBUG, format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	logf(ctx, filename, strconv.Itoa(line), ERROR, format, args...)
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
