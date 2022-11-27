package main

import (
	"bufio"
	"fmt"
	"mserverlauncher/app"
	"mserverlauncher/app/color"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func main() {

	clear()

	fmt.Printf(color.Blue + "Bienvenue dans le Minecraft Server Launcher\n\n" + color.Reset)

	config, err := app.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	if config.RootPath == "" {
		panic("root_path is not defined in config.json")
	}

	printServers(app.LoadServers(config))

	userPrompt(config)

	fmt.Printf("Goodbye !😊\n")
}

func userPrompt(config app.Config) {
	printHelp()
	scanner := bufio.NewScanner(os.Stdin)

mainloop:
	for {
		fmt.Print("m-server : ")

		//text := readMainCommand()
		scanner.Scan()
		text := scanner.Text()
		parts := strings.Split(text, " ")

		switch parts[0] {
		case "help":
			printHelp()
		case "list":
			printServers(app.LoadServers(config))
		case "start":
			startServer(parts, text, config)

			clear()
			printServers(app.LoadServers(config))
			printHelp()

		case "props":
			if len(parts) < 2 {
				fmt.Printf("'%s' n'est pas une commande valide, il faut donner un numero de serveur aussi (par exemple 1 ou 2 etc ...)\n", text)
				break
			}

			server, err := app.LoadServer(parts[1], config)
			if err != nil {
				fmt.Printf("'%s' n'est pas une commande valide, il faut donner un numero de serveur existant aussi (par exemple 1 ou 2 etc ...)\n", text)
				break
			}
			fmt.Println(server.Properties)

		case "new":
			newServer(config)
			clear()
			printServers(app.LoadServers(config))
			printHelp()

		case "edit":
			if len(parts) < 2 {
				fmt.Printf("'%s' n'est pas une commande valide, il faut donner un numero de serveur aussi (par exemple 1 ou 2 etc ...)\n", text)
				break
			}

			server, err := app.LoadServer(parts[1], config)
			if err != nil {
				fmt.Printf("'%s' n'est pas une commande valide, il faut donner un numero de serveur existant aussi (par exemple 1 ou 2 etc ...)\n", text)
				break
			}

			cmd := exec.Command("code")
			cmd.Args = append(cmd.Args, path.Join(server.Path, "server.properties"))
			//cmd.Stdout = os.Stdout
			cmd.Run()
			cmd.Wait()

		case "q", "quit", "Q", "QUIT", "exit", "stop":
			break mainloop
		case "":
		default:
			fmt.Printf("'%s' n'est pas une command valide, tape 'help' pour plus d'information\n", text)
		}

	}
}

func startServer(parts []string, text string, config app.Config) {
	if len(parts) < 2 {
		fmt.Printf("'%s' n'est pas une commande valide, il faut donner un numero de serveur aussi (par exemple 1 ou 2 etc ...)\n", text)
		return
	}

	server, err := app.LoadServer(parts[1], config)
	if err != nil {
		fmt.Printf("'%s' n'est pas une commande valide, %s (par exemple 1 ou 2 etc) ...)\n", text, err.Error())
		return
	}

	printIP()
	server.StartServer()
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printIP() {

	fmt.Print("SERVER IP : ")

	cmd := exec.Command("ifconfig eth0 | grep 'inet ' | cut -d: -f2 |awk '{print $2}'")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println()
}

func newServer(config app.Config) {
	fmt.Printf(color.Blue + "CREATION DE NOUVEAU SERVEUR !\n" + color.Reset)
	name := promptChoice("Entrez le Nom du serveur : ")
	sname := sanitizeFileName(name)
	if sname == "" {
		fmt.Println("The name is empty or invalid")
		return
	}

	versions := app.LoadVersions(config)
	if len(versions) == 0 {
		fmt.Printf("Error: Can not find any .jar files in the versions folders, please download minecraft.server.jar and put it under the root_path/versions/ folder")
		return
	}

	srvInfo := app.Server{
		Path:       path.Join(config.RootPath, config.ServersPath, sname),
		Version:    promptChoice("Choisissez la version du serveur:\n", versions...),
		Properties: app.ParseServerPropertiesString(app.DefaultServerPropertiesContent),
	}

	srvInfo.SetName(name)
	srvInfo.SetGameMode(promptChoice("Choisissez un mode de jeux:\n", "survival", "creative", "adventure"))
	if srvInfo.GameMode() != "creative" {
		srvInfo.SetDifficulty(promptChoice("Choisissez une difficultée:\n", "peaceful", "easy", "normal", "hard", "hardcore"))
	}

	srvInfo.SetSeed(promptChoice("Entrez une graine (seed): "))

	err := app.CreateServer(config, srvInfo)
	if err != nil {
		fmt.Printf("Echec de la création du serveur : %s\n", err)
	}
	fmt.Println(color.Blue + "Le serveur a été crée avec succès. Utilise la commande start pour le démmarer" + color.Reset)
}

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "&", "_")
	name = strings.ReplaceAll(name, ",", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, ";", "_")
	name = strings.ReplaceAll(name, "$", "_")
	return name
}

func promptChoice(text string, choices ...string) string {
	scanner := bufio.NewScanner(os.Stdin)

	printChoices := func() {
		fmt.Print(text)
		for i, c := range choices {
			fmt.Printf("  %d => %s\n", i+1, c)
		}
	}

	printChoices()

	for scanner.Scan() {
		if len(choices) == 0 {
			return scanner.Text()
		}

		if i, err := strconv.Atoi(scanner.Text()); err == nil {
			i = i - 1
			if i >= 0 && i < len(choices) {
				return choices[i]
			} else {
				fmt.Printf("Choix invalide %s, recommencez svp\n", scanner.Text())
			}
		} else {
			fmt.Println(err)
		}
		printChoices()
	}
	return ""
}

func printHelp() {
	fmt.Println("Aide:")
	fmt.Println(color.Green + "    help                   " + color.Reset + ": affiche l'aide")
	fmt.Println(color.Green + "    list                   " + color.Reset + ": affiche la liste des serveurs")
	fmt.Println(color.Green + "    start X                " + color.Reset + ": démmarre le serveur numéro X (par example: start 1)")
	fmt.Println(color.Green + "    props X                " + color.Reset + ": affiche les propriétés du server X (par example: props 1)")
	fmt.Println(color.Green + "    edit  X                " + color.Reset + ": edite la configuration du serveur (par example: edit 1)")
	fmt.Println(color.Green + "    new                    " + color.Reset + ": créee un nouveau server")
	fmt.Println(color.Green + "    quit (ou stop ou exit) " + color.Reset + ": quitte minecraft s/stoperver launcher")
	fmt.Println()
}

func readMainCommand() {

	//reader := bufio.NewReader(os.Stdin)
	// var cmd string =""
	// for rune, _, err := reader.ReadRune() ; err == nil && rune != '\n' ;{
	// 	cmd += rune
	// }

}

func printServers(servers []*app.Server) {

	const format string = "│%10s │ %30s │ %15s │ %10s │ %10s │ %10s │ %10s │\n"
	//tWidth := 89
	tWidth := len(fmt.Sprintf(format, "", "", "", "", "", "", "")) - 17

	fmt.Printf("Liste des serveurs disponibles:\n")
	fmt.Printf("┎%s┒\n", strings.Repeat("─", tWidth-2))
	fmt.Printf(color.Blue+format+color.Reset, "Numéro", "Nom", "Monde", "Seed", "Port", "Gamemode", "Difficulté")
	fmt.Printf("├%s┤\n", strings.Repeat("─", tWidth-2))
	for i, s := range servers {
		_, err := os.ReadFile(path.Join(s.Path, "server.properties"))
		if err == nil {
			fmt.Printf(format,
				fmt.Sprintf("%d", i+1),
				s.Name(),
				s.Properties["level-name"],
				s.Properties["level-seed"],
				s.Properties["server-port"],
				s.GameMode(),
				s.Difficulty())
		}
	}
	fmt.Printf("┖%s┚\n", strings.Repeat("─", tWidth-2))
	fmt.Println()

}
