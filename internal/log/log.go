package log

import (
	"context"
	"log"
	"os"
	"path/filepath"

	icontext "github.com/YunchengHua/hotels-data-merge/internal/context"
)

// TO-DO: Be able to config the log path, log level
const (
	defaultLogPath = "./log"
	defaultLogFile = "application.log"
)

func MustInit() {
	// Create log directory if it doesn't exist
	err := os.MkdirAll(defaultLogPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new log file
	logFile, err := os.OpenFile(
		filepath.Join(defaultLogPath, defaultLogFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening or creating log file: %v", err)
	}

	// Set the output of the standard logger to the log file
	log.SetOutput(logFile)
}

func Error(ctx context.Context, msg string, v ...any) {
	log.Printf("Error|Trace:"+icontext.GetTraceID(ctx)+" : "+msg+"\n", v...)
}

func Warn(ctx context.Context, msg string, v ...any) {
	log.Printf("Warn|Trace:"+icontext.GetTraceID(ctx)+" : "+msg+"\n", v...)
}

func Info(ctx context.Context, msg string, v ...any) {
	log.Printf("Info|Trace:"+icontext.GetTraceID(ctx)+" : "+msg+"\n", v...)
}

func Debug(ctx context.Context, msg string, v ...any) {
	log.Printf("Debug|Trace:"+icontext.GetTraceID(ctx)+" : "+msg+"\n", v...)
}
