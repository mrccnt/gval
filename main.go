package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v3"
)

const (
	exitOK         = 0
	exitValidation = 1
	exitPattern    = 2
)

var version = "0.0.0-dev"

func main() {
	app := &cli.Command{}
	app.Name = "gval"
	app.Version = version
	app.Usage = "CLI Validator"
	app.Description = "CLI validator utilizing go-playground/validator."
	app.ArgsUsage = `<value> <pattern>
	
ARGUMENTS:
	<value>        the value to validate
	<pattern>      validation tag as of go-playground/validator/v10`
	app.Arguments = []cli.Argument{
		&cli.StringArg{
			Name:      "value",
			UsageText: "the value to validate",
		},
		&cli.StringArg{
			Name:      "pattern",
			UsageText: "validation tag as of go-playground/validator/v10",
		},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "int",
			Aliases: []string{"i"},
			Usage:   "convert string to int before validating",
		},
	}
	app.Action = func(_ context.Context, command *cli.Command) error {
		var finalVal any

		sVal := command.StringArg("value")
		sPat := command.StringArg("pattern")

		v := validator.New(validator.WithRequiredStructEnabled())

		if err := v.Var(sPat, "required,gt=0,lte=255"); err != nil {
			return cli.Exit("", exitPattern)
		}

		finalVal = sVal

		if command.Bool("i") {
			var err error
			if finalVal, err = strconv.Atoi(sVal); err != nil {
				return cli.Exit("", exitValidation)
			}
		}

		var err error

		if err = validate(v, finalVal, sPat); err == nil {
			return cli.Exit("", exitOK)
		}

		if _, ok := errors.AsType[validator.ValidationErrors](err); ok {
			return cli.Exit("", exitValidation)
		}

		return cli.Exit("", exitPattern)
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		os.Exit(exitOK)
	}
}

func validate(v *validator.Validate, value any, tag string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invalid validation tag %q: %v", tag, r)
		}
	}()
	return v.Var(value, tag)
}
