package lib

import (
	"fmt"
	"os"
)

type CmdOptions struct {
	Date []string `long:"date" required:"false" description:"Set a single date filter. Use the format YYYY[-MM[-DD]]. Can be applied multiple times."`
}

type AppConfig struct {
	BaseOutputPath string
	Secrets        *Secrets

	DateFilter []string
}

func CreateAppConfig(args []string, opts CmdOptions) AppConfig {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Output path missing.\n")
		os.Exit(1)
	}
	outPath := args[0]
	err := os.MkdirAll(outPath, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	conf := AppConfig{
		BaseOutputPath: outPath,
		Secrets:        LoadSecrets(),
		DateFilter:     opts.Date,
	}
	if err := conf.Secrets.EnsureUserSecrets(); err != nil {
		panic(err)
	}

	return conf
}
