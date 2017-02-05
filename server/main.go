package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type machine struct {
	Key      string `json:"key"`
	Hostname string `json:"hostname"`
	IP       string
	OutIP    string `json:"outIP"`
	C        string `json:"type"`
	Edit     bool   `json:"edit"`
}

func main() {
	e := echo.New()
	e.Static("/", "../client/")
	e.File("/", "../client/index.html")
	e.GET("/machines", func(c echo.Context) error {
		fmt.Println("getgetok")
		m := machine{
			Key:      "host0",
			Hostname: "host0.1.1.1",
			IP:       "192.168.1.1",
			OutIP:    "1.1.1.1",
			C:        "login",
			Edit:     false,
		}
		var a []machine
		a = append(a, m)
		b, err := json.Marshal(a)
		if err != nil {
			fmt.Println("json err:", err)
		}

		return c.Blob(http.StatusOK, "application/json", b)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
