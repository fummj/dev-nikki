package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/api/home"
	"dev_nikki/internal/api/login"
	"dev_nikki/internal/api/signup"
	"dev_nikki/internal/handlers/index"
)

func main() {
	app := echo.New()
	app.Static("/static", "./static/dist")

	wildCardHandler := index.WildCardHandler{}

	app.POST("/api/login", login.Login)
	app.POST("/api/signup", signup.SignUp)
	app.GET("/api/home", home.Home)
	app.GET("/*", wildCardHandler.FallbackToIndex)

	fmt.Print("🛎️  dev_nikki 🛎️" + "\n")
	fmt.Println("#############################################################")

	app.Logger.Fatal(app.Start(":8080"))
}
