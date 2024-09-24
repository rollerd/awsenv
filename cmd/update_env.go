package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/ini.v1"
)

func updateEnv(newEnv string) {
	// Create path to aws credentials file (default location)
	homeDir := os.Getenv("HOME")
	awsCredFile := fmt.Sprintf("%s/.aws/credentials", homeDir)
	awsConfigFile := fmt.Sprintf("%s/.aws/config", homeDir)

	// Create path to file that holds the name of the new default environment (useful for tmux status bar info)
	envFile := fmt.Sprintf("%s/.awsenv", homeDir)

	updateConfig(awsCredFile, awsConfigFile, newEnv)

	updateEnvFile(newEnv, envFile)
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
func updateConfig(awsCredFile, awsConfigFile, newEnv string) {
	// load credentials file
	credCfg, err := ini.Load(awsCredFile)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	// load config file
	configCfg, err := ini.Load(awsConfigFile)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	// check if section in credentials
	if sectionExists(newEnv, credCfg) {
	} else {
		fmt.Printf("Error: Could not find section: %s in credentials file\n", newEnv)
		os.Exit(1)
	}

	// check if section in config
	if sectionExists(newEnv, configCfg) {
	} else {
		fmt.Printf("Error: Could not find section: %s in config file\n", newEnv)
		os.Exit(1)
	}

	// create default section if it doesnt exist in credentials
	if sectionExists("default", credCfg) {
	} else {
		_, err := credCfg.NewSection("default")
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}

	// create default section if it doesnt exist in config
	if sectionExists("default", configCfg) {
	} else {
		_, err := configCfg.NewSection("default")
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}

	// remove default section if it exists
	credCfg.DeleteSection("default")
	configCfg.DeleteSection("default")

	// create new default section in credentials
	for _, key := range credCfg.Section(newEnv).KeyStrings() {
		_, err := credCfg.Section("default").NewKey(key, credCfg.Section(newEnv).Key(key).String())
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}

	// create new default section in config
	for _, key := range configCfg.Section(newEnv).KeyStrings() {
		_, err := configCfg.Section("default").NewKey(key, configCfg.Section(newEnv).Key(key).String())
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}

	// backup credentials file
	backupCreds(awsCredFile)

	fmt.Printf("Setting environment to: %s\n", newEnv)
	// write credentials file
	err = credCfg.SaveTo(awsCredFile)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
	// write config file
	err = configCfg.SaveTo(awsConfigFile)
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
