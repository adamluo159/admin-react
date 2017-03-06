package main

import (
	"fmt"

	"github.com/adamluo159/admin-react/server/agentServer"
	"github.com/adamluo159/admin-react/server/db"
	"github.com/adamluo159/admin-react/server/machine"
	"github.com/adamluo159/admin-react/server/zone"
	"github.com/labstack/echo"
)

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		fmt.Println(c.Request().Method, c.Request().RequestURI, err)
		return err
	}
}

func main() {

	db.Connect()

	e := echo.New()
	e.Use(ServerHeader)
	machine.Register(e)
	zone.Register(e)

	go agentServer.New(":3300")

	e.Static("/", "../client/")
	e.File("/", "../client/index.html")
	e.Logger.Fatal(e.Start(":1323"))

}
