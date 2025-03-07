package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/authN"
	"dev_nikki/internal/handlers/api"
	"dev_nikki/internal/handlers/index"
)

func main() {
	app := echo.New()
	app.Static("/static", "./static/dist")

	authHandler := api.AuthHandler{}
	homeHandler := api.HomeHandler{}
	wildCardHandler := index.WildCardHandler{}

	app.POST("/api/login", authHandler.Login)
	app.POST("/api/signup", authHandler.SignUp)
	app.GET("/api/home", homeHandler.Home)
	app.GET("/auth/login", authN.OAuth2)
	app.GET("/auth/callback", authN.OAuth2Callback)
	app.GET("/*", wildCardHandler.FallbackToIndex)

	fmt.Print("🛎️  dev_nikki 🛎️" + "\n")
	fmt.Println("#############################################################")

	app.Logger.Fatal(app.Start(":8080"))
}
