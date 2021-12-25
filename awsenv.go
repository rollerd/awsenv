package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
)

func main() {
	// Create path to aws credentials file (default location)
	homeDir := os.Getenv("HOME")
	awsCredFile := fmt.Sprintf("%s/.aws/credentials", homeDir)

	// Create path to file that holds the name of the new default environment (useful for tmux status bar info)
	envFile := fmt.Sprintf("%s/.awsenv", homeDir)

	newEnv := getEnv()

	updateConfig(awsCredFile, newEnv)

	updateEnvFile(newEnv, envFile)
}

// Get target config from commandline arg
func getEnv() (env string) {
	if len(os.Args) < 2 {
		fmt.Println("awsenv requires one argument: <env to use>")
		os.Exit(1)
	}
	env = os.Args[1]
	return
}

// Create backup of existing cred file
func backupCreds(awsCredFile string) {
	awsCredFileBackup := awsCredFile + ".bkup"

	input, err := ioutil.ReadFile(awsCredFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(awsCredFileBackup, input, 0644)
	if err != nil {
		fmt.Println("Error creating", awsCredFileBackup)
		fmt.Println(err)
		os.Exit(1)
	}
}

// Create new version of config file with the newEnv as [default]
func updateConfig(awsCredFile, newEnv string) {
	cfg, err := ini.Load(awsCredFile)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	if sectionExists(newEnv, cfg) {
	} else {
		fmt.Printf("Error: Could not find config section: %s in credentials file\n", newEnv)
		os.Exit(1)
	}

	if sectionExists("default", cfg) {
	} else {
		_, err := cfg.NewSection("default")
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}

	// remove default section if it exists
	cfg.DeleteSection("default")

	for _, key := range cfg.Section(newEnv).KeyStrings() {
		_, err := cfg.Section("default").NewKey(key, cfg.Section(newEnv).Key(key).String())
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}

	backupCreds(awsCredFile)

	fmt.Printf("Setting environment to: %s\n", newEnv)
	err = cfg.SaveTo(awsCredFile)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}

// Check if config file [section] exists
func sectionExists(sectionName string, cfg *ini.File) bool {
	for _, v := range cfg.SectionStrings() {
		if sectionName == v {
			return true
		}
	}
	return false
}

// Write the new target env to a file. Useful for things like tmux status bar
func updateEnvFile(newEnv, envFile string) {
	data := []byte(newEnv)
	err := os.WriteFile(envFile, data, 0644)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
