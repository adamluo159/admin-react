package main

import (
	"fmt"

	"github.com/adamluo159/admin-react/server/db"
	"github.com/adamluo159/admin-react/server/machine"
	"github.com/labstack/echo"
)

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		fmt.Println("aaaaaa", c.Request().RequestURI, c.Request().Method, err)
		return err
	}
}

func main() {

	db.Connect()

	e := echo.New()
	e.Use(ServerHeader)
	machine.Register(e)

	e.Static("/", "../client/")
	e.File("/", "../client/index.html")
	e.Logger.Fatal(e.Start(":1323"))
}
