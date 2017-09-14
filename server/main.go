package main

import "github.com/adamluo159/admin-react/server/yada"

func main() {
	y := yada.New("./config.json")
	y.Run()
}
