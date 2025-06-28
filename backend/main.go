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
	app.GET("/api/auth/login", authN.OAuth2)
	app.GET("/auth/callback", authN.OAuth2Callback)
	app.GET("/api/home", homeHandler.Home)
	app.GET("/api/home/:phase", homeHandler.PreHome)
	app.POST("/api/home/:phase", homeHandler.PostPreHome)
	app.PUT("/api/home/updateMarkdown", homeHandler.UpdateMarkdown)
	app.POST("/api/home/create-new-folder", homeHandler.CreateNewFolder)
	app.POST("/api/home/create-new-file", homeHandler.CreateNewFile)
	app.DELETE("/api/home/delete-folder", homeHandler.DeleteFolder)
	app.DELETE("/api/home/delete-file", homeHandler.DeleteFile)
	app.GET("/*", wildCardHandler.FallbackToIndex)

	fmt.Print("🛎️  dev_nikki 🛎️" + "\n")
	fmt.Println("#############################################################")

	// 開発環境
	app.Logger.Fatal(app.Start(":8080"))

	// 本番環境
	// certPath := "/usr/src/certs/fullchain.pem"
	// keyPath := "/usr/src/certs/privkey.pem"
	//
	// if err := app.StartTLS(":443", certPath, keyPath); err != nil {
	// 	app.Logger.Fatal(err)
	// }
	//
	// s := app.TLSServer
	// fmt.Println("opened port: ", s.Addr)
}
