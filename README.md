# MINECRAFT SERVER LAUNCHER

## WHAT
A simple command line tool to easly create and manage and configure Minecraft Servers 
Designed to be used by children

## Limitation
Works on Linux only but can be run with windows using WSL
you need to create this arboresence folder in your home dir and 
mserverlauncher
    /versions/server.jar  <= download from official minecraft site
    /servers/

## HOW
This simple tool is written in golang

````
go build
go install
mserverlauncher
````



## TODO
- download automaticcaly latest minecraft server jar file
- backup saves management
- manage remote servers




