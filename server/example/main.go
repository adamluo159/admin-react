package main

import "github.com/adamluo159/admin-react/server"

func main() {
	y := yada.New("./config.json")
	y.Run()
}
