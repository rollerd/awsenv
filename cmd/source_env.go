package cmd

import (
	"fmt"
	"os"
	"strings"

	"golang.design/x/clipboard"

	"gopkg.in/ini.v1"
)

func sourceProfile(sourceEnv string) {
	// Create path to aws credentials file (default location)
	homeDir := os.Getenv("HOME")
	awsCredFile := fmt.Sprintf("%s/.aws/credentials", homeDir)

	cfg, err := ini.Load(awsCredFile)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	if sectionExists(sourceEnv, cfg) {
	} else {
		fmt.Printf("Could not find a profile called: %s\n", sourceEnv)
		os.Exit(1)
	}

	var exportString strings.Builder

	for _, key := range cfg.Section(sourceEnv).KeyStrings() {
		newString := fmt.Sprintf("export %s=%s\n", strings.ToUpper(key), cfg.Section(sourceEnv).Key(key).String())
		exportString.WriteString(newString)
	}

	fmt.Printf("Copied %s values to clipboard!\n", sourceEnv)
	copyToClipboard(exportString)
}

func copyToClipboard(exportString strings.Builder) {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	clipboard.Write(clipboard.FmtText, []byte(exportString.String()))
}
