package lib

import (
	"fmt"
	"os"
)

type CmdOptions struct {
	Date               []string `long:"date" value-name:"YYYY[-MM[-DD]]" description:"Set a single date filter, in the format YYYY-MM-DD. Partial date possible.\nExample: --date=2023 will fetch all photos from the year 2023."`
	ForceOverride      bool     `long:"force" short:"f" description:"Overrides local files in any case. Default is skip if a file exists locally."`
	ForceNewerOverride bool     `long:"force-newer" short:"n" description:"Override local files only if the remote file is newer, skip otherwise"`
	NrOfThreads        int      `long:"threads" short:"t" default:"5" value-name:"nr" description:"Number of download threads to use"`
}

type AppConfig struct {
	BaseOutputPath string
	Secrets        *Secrets

	DateFilter         []string
	ForceOverride      bool
	ForceNewerOverride bool
	NrOfThreads        int
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
		BaseOutputPath:     outPath,
		Secrets:            LoadSecrets(),
		DateFilter:         opts.Date,
		ForceOverride:      opts.ForceOverride,
		ForceNewerOverride: opts.ForceNewerOverride,
		NrOfThreads:        opts.NrOfThreads,
	}
	if conf.NrOfThreads < 1 {
		conf.NrOfThreads = 1
	}
	if err := conf.Secrets.EnsureUserSecrets(); err != nil {
		panic(err)
	}

	return conf
}
