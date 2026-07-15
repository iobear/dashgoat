/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/labstack/gommon/log"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
var logFile *os.File

// setupLogger configures the global slog logger.
// If path is non-empty, logs are written to the file (append mode).
// Otherwise stdout is used. Format is "txt" (default) or "json".
func setupLogger(path, format string) {
	var handler slog.Handler
	if path != "" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logger.Error("setupLogger", "err", err)
			os.Exit(1)
		}
		logFile = f
		if format == "json" {
			handler = slog.NewJSONHandler(f, nil)
		} else {
			handler = slog.NewTextHandler(f, nil)
		}
	} else {
		if format == "json" {
			handler = slog.NewJSONHandler(os.Stdout, nil)
		} else {
			handler = slog.NewTextHandler(os.Stdout, nil)
		}
	}
	logger = slog.New(handler)
}

// slogAdapter implements echo.Logger by forwarding to the global slog logger.

type slogAdapter struct {
	prefix string
	level  log.Lvl
}

func (a *slogAdapter) Output() io.Writer     { return os.Stdout }
func (a *slogAdapter) SetOutput(w io.Writer) {}
func (a *slogAdapter) Prefix() string        { return a.prefix }
func (a *slogAdapter) SetPrefix(p string)    { a.prefix = p }
func (a *slogAdapter) Level() log.Lvl        { return a.level }
func (a *slogAdapter) SetLevel(v log.Lvl)    { a.level = v }
func (a *slogAdapter) SetHeader(h string)    {}

func (a *slogAdapter) msg(args ...interface{}) string {
	s := fmt.Sprint(args...)
	if a.prefix != "" {
		s = a.prefix + ": " + s
	}
	return s
}

func (a *slogAdapter) Print(args ...interface{})       { logger.Info(a.msg(args...)) }
func (a *slogAdapter) Printf(format string, args ...interface{}) { logger.Info(fmt.Sprintf(format, args...)) }
func (a *slogAdapter) Printj(j log.JSON)               { logger.Info(a.msg(j)) }

func (a *slogAdapter) Debug(args ...interface{})       { logger.Debug(a.msg(args...)) }
func (a *slogAdapter) Debugf(format string, args ...interface{}) { logger.Debug(fmt.Sprintf(format, args...)) }
func (a *slogAdapter) Debugj(j log.JSON)               { logger.Debug(a.msg(j)) }

func (a *slogAdapter) Info(args ...interface{})        { logger.Info(a.msg(args...)) }
func (a *slogAdapter) Infof(format string, args ...interface{})  { logger.Info(fmt.Sprintf(format, args...)) }
func (a *slogAdapter) Infoj(j log.JSON)                { logger.Info(a.msg(j)) }

func (a *slogAdapter) Warn(args ...interface{})        { logger.Warn(a.msg(args...)) }
func (a *slogAdapter) Warnf(format string, args ...interface{})   { logger.Warn(fmt.Sprintf(format, args...)) }
func (a *slogAdapter) Warnj(j log.JSON)                { logger.Warn(a.msg(j)) }

func (a *slogAdapter) Error(args ...interface{})       { logger.Error(a.msg(args...)) }
func (a *slogAdapter) Errorf(format string, args ...interface{}) { logger.Error(fmt.Sprintf(format, args...)) }
func (a *slogAdapter) Errorj(j log.JSON)               { logger.Error(a.msg(j)) }

func (a *slogAdapter) Fatal(args ...interface{})       { logger.Error(a.msg(args...)); os.Exit(1) }
func (a *slogAdapter) Fatalf(format string, args ...interface{}) { logger.Error(fmt.Sprintf(format, args...)); os.Exit(1) }
func (a *slogAdapter) Fatalj(j log.JSON)               { logger.Error(a.msg(j)); os.Exit(1) }

func (a *slogAdapter) Panic(args ...interface{})       { logger.Error(a.msg(args...)); panic(a.msg(args...)) }
func (a *slogAdapter) Panicf(format string, args ...interface{}) { logger.Error(fmt.Sprintf(format, args...)); panic(fmt.Sprintf(format, args...)) }
func (a *slogAdapter) Panicj(j log.JSON)               { logger.Error(a.msg(j)); panic(a.msg(j)) }
