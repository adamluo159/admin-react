package main

import (
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/adamluo159/admin-react/server/db"
	"github.com/adamluo159/admin-react/server/machine"
	"github.com/adamluo159/admin-react/server/zone"
	"github.com/adamluo159/gameAgent/agentServer"
	"github.com/labstack/echo"
	permissions "github.com/xyproto/permissions2"
)

var (
	perm *permissions.Permissions
)

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		fmt.Println(c.Request().Method, c.Request().RequestURI, err)
		return err
	}
}

func DbPing(e *echo.Echo) {
	for {
		err := db.Session.Ping()
		if err != nil {
			rerr := db.ReConnect()
			if rerr == nil {
				machine.Register(e)
				//zone.Register(e)
			}
		}
		time.Sleep(time.Second * 10)
	}
}

func RegisterPerm(redisHost string, redisPwd string, e *echo.Echo) {
	userstate, err := permissions.NewUserStateWithPassword2(redisHost, redisPwd)
	if err != nil {
		log.Fatal(err)
	}
	perm = permissions.NewPermissions(userstate)
	perm.AddUserPath("/machine")
	perm.AddUserPath("/zone")
	f := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if perm.Rejected(c.Response(), c.Request()) {
				// Deny the request
				//return echo.NewHTTPError(http.StatusForbidden, denyMessage)
				return c.String(http.StatusOK, "verify")
			}
			// Continue the chain of middleware
			return next(c)
		}
	}
	e.Use(f)
}

func main() {

	e := echo.New()
	e.Use(ServerHeader)
	RegisterPerm("192.168.1.252", "", e)
	Login := func(c echo.Context) error {
		perm.UserState().AddUser("bob", "hunter1", "bob@zombo.com")
		perm.UserState().AddUser("adamluo", "adamluo0011", "bo@zombo.com")
		user := c.FormValue("user")
		passwd := c.FormValue("password")

		log.Println("user:", user, "passwd:", passwd, "res=")
		if perm.UserState().CorrectPassword(user, passwd) {
			err := perm.UserState().Login(c.Response().Writer, user)
			if err != nil {
				c.String(http.StatusOK, err.Error())
			} else {
				c.String(http.StatusOK, "admin")
			}
		} else {
			c.String(http.StatusInternalServerError, "Login fail")
		}
		log.Println("wwwwwwww-", perm.UserState().CookieSecret())
		return nil
	}

	db.Connect()
	go DbPing(e)

	e.POST("/login", Login)

	s := agentServer.New(":3300")
	m := machine.Register(e)
	z := zone.Register(e)

	m.InitMgr(s)
	z.InitMgr(m, s)
	s.Init(m)

	go s.Listen()

	e.Static("/", "../client/dist/")
	e.File("/", "../client/dist/index.html")
	e.Logger.Fatal(e.Start(":1323"))

}
