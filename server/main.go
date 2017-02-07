package main

import (
	"github.com/adamluo159/admin-react/server/db"
	"github.com/adamluo159/admin-react/server/machine"
	"github.com/labstack/echo"
)

func main() {

	db.Connect()

	e := echo.New()
	machine.Register(e)

	e.Static("/", "../client/")
	e.File("/", "../client/index.html")
	e.Logger.Fatal(e.Start(":1323"))
}
