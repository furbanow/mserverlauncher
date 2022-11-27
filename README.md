# MINECRAFT SERVER LAUNCHER

## WHAT
A simple command line tool to easly create and manage and configure Minecraft Servers 
Designed to be used by children

## Limitation
- Works on Linux only but can be run with windows using WSL
- You need to create this arboresence folder in your home dir and 
games/minecraft_server
    /versions/server.jar  <= download from official minecraft site and put it here
    /servers/ <= the servers instances will be created here

## HOW
This simple tool is written in golang

````
go build
go install
mserverlauncher

Bienvenue dans le Minecraft Server Launcher

Liste des serveurs disponibles:
┎──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┒
│    Numéro │                            Nom │           Monde │       Seed │       Port │   Gamemode │ Difficulté │
├──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│         1 │         Leiwa Minecraft Server │           world │      leiwa │      55555 │   survival │       hard │
│         2 │                   La Familiale │           world │            │      55556 │   survival │     normal │
│         3 │         La Machinerie Creative │           world │ machinerie │      55557 │   creative │       hard │
┖──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┚

Aide:
    help                   : affiche l'aide
    list                   : affiche la liste des serveurs
    start X                : démmarre le serveur numéro X (par example: start 1)
    props X                : affiche les propriétés du server X (par example: props 1)
    edit  X                : edite la configuration du serveur (par example: edit 1)
    new                    : créee un nouveau server
    quit (ou stop ou exit) : quitte minecraft s/stoperver launcher

````



## TODO
- download automaticcaly latest minecraft server jar file
- backup saves management
- manage remote servers




