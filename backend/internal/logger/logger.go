package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

const (
	debug = "\033[32m"
	info  = "\033[34m"
	warn  = "\033[33m"
	errr  = "\033[31m"
	end   = "\033[0m\n"
)

var (
	logLevelColors = map[slog.Level]string{
		slog.LevelDebug: debug,
		slog.LevelInfo:  info,
		slog.LevelWarn:  warn,
		slog.LevelError: errr,
	}

	Slog *slog.Logger
)

func init() {
	jsonHandler := &CustomJsonHandler{*slog.NewJSONHandler(os.Stdout, nil)}
	Slog = slog.New(jsonHandler)
}

type CustomJsonHandler struct {
	slog.JSONHandler
}

func (h *CustomJsonHandler) Handle(ctx context.Context, r slog.Record) error {

	level := logLevelColors[r.Level]
	logMessage := map[string]any{
		"Time":    r.Time,
		"Level":   r.Level.String(),
		"Message": r.Message,
	}
	logAttr := map[string]any{}

	for attr := range r.Attrs {
		fmt.Println(attr.Key, attr.Value)
		logAttr[attr.Key] = attr.Value.String()
	}
	fmt.Println("logAttr", logAttr)

	jsonLog, err := json.Marshal(logMessage)
	if err != nil {
		return err
	}

	jsonLogAttr, err := json.Marshal(logAttr)
	if err != nil {
		return err
	}

	fmt.Printf("%s%s%s%s%s", level, r.Level.String()+": ", jsonLog, jsonLogAttr, end)

	return nil
}

func NewCustomJsonHandler() *CustomJsonHandler {
	return &CustomJsonHandler{}
}
