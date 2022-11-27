package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type Config struct {
	RootPath     string `json:"root_path"`
	ServersPath  string `json:"servers_path"`
	VersionsPath string `json:"versions_path"`
}

func configFilePath() string {
	configFile := ".mserverlauncher-config.json"
	usr, _ := user.Current()
	usrDir := usr.HomeDir

	configFilePath := path.Join(usrDir, configFile)
	return configFilePath
}

func mserverlauncherPath() string {
	configFile := "games/minecraft_server/"
	usr, _ := user.Current()
	usrDir := usr.HomeDir

	configFilePath := path.Join(usrDir, configFile)
	return configFilePath
}

func LoadConfig() (Config, error) {

	configFilePath := configFilePath()

	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return CreateConfig()
	}

	var config Config
	buf, err := os.ReadFile(configFilePath)
	if err != nil {

		return Config{}, err
	}
	err = json.Unmarshal(buf, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func CreateConfig() (Config, error) {

	fmt.Printf("Creating config file at %s\n", mserverlauncherPath())
	c := Config{
		RootPath:     mserverlauncherPath(),
		ServersPath:  "servers",
		VersionsPath: "versions"}

	buf, _ := json.MarshalIndent(c, "", "\t")

	return c, ioutil.WriteFile(configFilePath(), buf, 0644)

	//return c, nil

	// _, err := os.Stat("~/games")
	// if os.IsNotExist(err) {
	// 	err = os.Mkdir("~/games", 0755)
	// 	return c, err
	// }
	// _, err = os.Stat("~/games/minecraft_server")
	// if os.IsNotExist(err) {
	// 	err = os.Mkdir("~/games/minecraft_server", 0755)
	// 	return c, err
	// }

}
