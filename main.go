package main

import (
	scrapper "learngo/scarpper"
	"strings"

	"github.com/labstack/echo"
)

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handlescrape(c echo.Context) error {
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	return nil
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handlescrape)
	e.Logger.Fatal(e.Start(":1323"))
}
