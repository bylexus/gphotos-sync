package lib

import (
	"fmt"
	"os"
)

type AppConfig struct {
	BaseOutputPath string
	Secrets        *Secrets
}

func CreateAppConfig() AppConfig {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Output path missing.\n")
		os.Exit(1)
	}
	outPath := os.Args[1]
	err := os.MkdirAll(outPath, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	conf := AppConfig{
		BaseOutputPath: outPath,
		Secrets:        LoadSecrets(),
	}
	if err := conf.Secrets.EnsureUserSecrets(); err != nil {
		panic(err)
	}

	return conf
}
