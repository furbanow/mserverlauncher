package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"golang.org/x/exp/slices"
)

type Server struct {
	Id         int
	Path       string
	Version    string
	Properties map[string]string
}

func (s Server) SetPort(port int) {
	s.Properties["server-port"] = strconv.Itoa(port)
}

func (s Server) Port() int {
	v, err := strconv.Atoi(s.Properties["server-port"])
	if err != nil {
		return 0
	}

	return v
}

func (s Server) SetName(name string) {
	s.Properties["motd"] = name
}

func (s Server) Name() string {
	return s.Properties["motd"]
}

func (s Server) SetGameMode(mode string) {
	s.Properties["gamemode"] = mode
}

func (s Server) GameMode() string {
	return s.Properties["gamemode"]
}

func (s Server) SetDifficulty(diff string) {

	if diff == "hardcore" {
		s.Properties["difficulty"] = "hard"
		s.Properties["hardcore"] = "true"
	}

	s.Properties["difficulty"] = diff
	s.Properties["hardcore"] = "false"

}

func (s Server) SetSeed(seed string) {
	s.Properties["level-seed"] = seed
}

func (s Server) Difficulty() string {

	if s.Properties["hardcore"] == "true" {
		return "hardcore"
	}

	return s.Properties["difficulty"]
}

func (s Server) StartServer() {

	fmt.Printf("Démmarrage du serveur %v\n", s)

	fmt.Printf("Changing directory to %s\n", s.Path)
	err := os.Chdir(s.Path)
	if err != nil {
		fmt.Printf("Cannot change directory to %s\n", s.Path)
	}

	cmdToRun := "/usr/bin/java"
	args := []string{"-Xmx1024M", "-Xms1024M", "-jar", "server.jar", "nogui"}
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	if process, err := os.StartProcess(cmdToRun, args, procAttr); err != nil {
		fmt.Printf("ERROR Unable to run %s: %s\n", cmdToRun, err.Error())
	} else {
		fmt.Printf("%s running as pid %d\n", cmdToRun, process.Pid)
		process.Wait()
	}

	fmt.Printf("Le serveur %d %s s'est bien arreté\n", s.Id, s.Name())
}

func LoadServers(config Config) []*Server {

	var servers []*Server

	serversPath := path.Join(config.RootPath, config.ServersPath)

	os.Chdir(config.RootPath)
	direntries, _ := os.ReadDir(config.ServersPath)
	for _, d := range direntries {
		if !d.IsDir() {
			continue
		}
		srvPropsPath := path.Join(serversPath, d.Name(), "server.properties")
		props, err := ParseServerPropertiesFile(srvPropsPath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		servers = append(servers, &Server{
			Path:       path.Join(serversPath, d.Name()),
			Properties: props,
		})

	}

	slices.SortFunc(servers, func(srv1, srv2 *Server) bool {
		return srv1.Port() < srv2.Port()
	})

	for i, s := range servers {
		servers[i] = &Server{
			Id:         i + 1,
			Path:       s.Path,
			Properties: s.Properties,
		}
	}

	return servers
}

func LoadServer(id string, config Config) (*Server, error) {

	srvId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		fmt.Printf("'%s' n'est pas un id valide %s\n", id, err.Error())
		return nil, err
	}

	servers := LoadServers(config)

	for _, s := range servers {
		if s.Id == int(srvId) {
			return s, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("server %d not found", srvId))
}

const (
	MinPort = 55555
)

func CreateServer(config Config, server Server) error {
	servers := LoadServers(config)
	if len(servers) > 0 {
		for i := range servers {
			if i == len(servers)-1 || servers[i+1].Port() > servers[i].Port()+1 {
				server.SetPort(servers[i].Port() + 1)
				break
			}
		}
	} else {
		server.SetPort(MinPort)
	}

	fmt.Printf("Creation du nouveau serveur %s version: %s, port : %d ...\n", server.Name(), server.Version, server.Port())
	if err := os.Chdir(config.RootPath); err != nil {
		return err
	}

	if err := os.Mkdir(server.Path, 0777); err != nil {
		return err
	}

	input, err := ioutil.ReadFile(path.Join(config.RootPath, config.VersionsPath, server.Version))
	if err != nil {
		fmt.Println(err)
		return err
	}

	dest := path.Join(server.Path, "server.jar")
	err = ioutil.WriteFile(dest, input, 0644)
	if err != nil {
		return err
	}

	srvPropsPath := path.Join(server.Path, "server.properties")
	err = SaveServerPropertiesFile(server.Properties, srvPropsPath)
	if err != nil {
		return err
	}

	dest = path.Join(server.Path, "eula.txt")
	err = ioutil.WriteFile(dest, []byte(Eula), 0644)
	if err != nil {
		return err
	}

	//TODO find server.jar path and parse version in load

	return nil
}
