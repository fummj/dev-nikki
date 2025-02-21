package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/api/login"
	"dev_nikki/internal/api/signup"
	"dev_nikki/internal/handlers/index"
)

func main() {
	app := echo.New()
	app.Static("/static", "./static/dist")

	wildCardHandler := index.WildCardHandler{}

	app.POST("/api/login", login.ReturnFormData)
	app.POST("/api/signup", signup.SendUserData)
	app.GET("/*", wildCardHandler.FallbackToIndex)

	fmt.Print("🛎️  dev_nikki 🛎️" + "\n")
	fmt.Println("#############################################################")

	// jwt
	// claim := authn.NewClaim(12, "masa", "email@email.com")
	// k := authn.NewJWTKeysKeeper()
	// t := authn.CreatePreSignedToken(claim)
	// fmt.Println(t)
	// signedJWT, _ := authn.CreateJWT(t, k)
	// fmt.Println("$$$$$$$$$$$$$$$")
	// fmt.Println(authn.ParseJWT(signedJWT, k.Publ))

	// create user
	// m := map[string]string{
	// 	"username": "hogehoge",
	// 	"email":    "hogehoge@hoge.com",
	// 	"password": "Ufis2f9j",
	// 	"salt":     models.GenerateSalt(),
	// }
	// fmt.Println(models.PasswordHashing(m["password"], m["salt"]))
	// models.CreateUser(DBC.DB, m)

	app.Logger.Fatal(app.Start(":8080"))
}
