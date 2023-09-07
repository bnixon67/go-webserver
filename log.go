/*
Copyright 2023 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strings"
)

const (
	logOpenFileFlag = os.O_APPEND | os.O_CREATE | os.O_WRONLY

	ownerReadWrite  = 0o600
	logOpenFileMode = ownerReadWrite
)

// LogLevelMap maps log level names to values.
// Lowercase names are used to allow case insensitive lookups.
var LogLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// LogLevels returns a comma-separated string of available log level.
func LogLevels() string {
	keys := make([]string, 0, len(LogLevelMap))

	// get all keys
	for key := range LogLevelMap {
		keys = append(keys, key)
	}

	// sort by key value, i.e., slog.Level
	sort.Slice(keys, func(i, j int) bool {
		return LogLevelMap[keys[i]] < LogLevelMap[keys[j]]
	})

	return strings.Join(keys, ", ")
}

// LogLevel returns the corresponding slog.Level value.
// The lookup is case insensitive.
func LogLevel(logLevelStr string) (slog.Level, error) {
	// assume LogLevelMap keys are lowercase
	logLevelStr = strings.ToLower(logLevelStr)

	logLevel, ok := LogLevelMap[logLevelStr]
	if !ok {
		return slog.LevelInfo, fmt.Errorf("invalid loglevel: %s", logLevelStr)
	}

	return logLevel, nil
}

var validLogTypes = []string{"json", "text"}

// InitLog initializes logging for the application.
func InitLog(name, handlerType string, level slog.Level, addSource bool) error {
	// configure log writter
	var w io.Writer = os.Stderr // default to Stderr is logFileName empty
	if name != "" {
		file, err := os.OpenFile(name, logOpenFileFlag, logOpenFileMode)
		if err != nil {
			return fmt.Errorf("InitLog: %w", err)
		}
		// do not defer file.Close() since the file must remain open
		w = file
	}

	// configure logger
	opts := &slog.HandlerOptions{
		AddSource: addSource,
		Level:     level,
	}

	// configure handler
	var handler slog.Handler
	switch handlerType {
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	case "text":
		handler = slog.NewTextHandler(w, opts)
	default:
		handler = slog.NewTextHandler(w, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("InitLog",
		slog.Group("log",
			slog.String("Filename", name),
			slog.String("Type", handlerType),
			slog.String("Level", level.String()),
			slog.Bool("AddSource", addSource),
		),
	)

	return nil
}
