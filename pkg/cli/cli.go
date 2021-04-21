// Copyright 2021 SpecializedGeneralist Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"fmt"
	"github.com/SpecializedGeneralist/translator/pkg/configuration"
	"github.com/SpecializedGeneralist/translator/pkg/models"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

// Run is the entry point to the CLI app.
func Run(arguments []string) {
	app := &cli.App{
		HelpName:  "translator",
		Usage:     "Translation service",
		Flags:     flags,
		Action:    runAction,
		Reader:    os.Stdin,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
	}
	err := app.Run(arguments)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "config",
		Aliases:  []string{"c"},
		Usage:    "load configuration from YAML `FILE`",
		Required: true,
	},
}

func runAction(ctx *cli.Context) (err error) {
	config, err := configuration.FromYAMLFile(ctx.String("config"))
	if err != nil {
		return err
	}

	logger := newLogger(zerolog.Level(config.LogLevel))

	defer func() {
		if err != nil {
			logger.Err(err).Send()
		}
	}()

	manager := models.NewManager(config, logger)
	err = manager.LoadModels()
	if err != nil {
		return err
	}

	// TODO: run server ...

	return nil
}

func newLogger(level zerolog.Level) zerolog.Logger {
	w := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	return zerolog.New(w).With().Timestamp().Logger().Level(level)
}
