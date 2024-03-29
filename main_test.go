package main

import (
	"net/http"
	"testing"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"go-boilerplate/db"
	"go-boilerplate/forms"
	"go-boilerplate/models"
	"go-boilerplate/migrations"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"log"
	"github.com/brianvoe/gofakeit/v6"
	"os"
)

func TestMain(t *testing.T) {

	err := godotenv.Load(".env.test", ".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := gin.New()

	//Custom form validator
	binding.Validator = new(forms.DefaultValidator)

	//app.Use(CORSMiddleware())
	//app.Use(RequestIDMiddleware())

	db.ConnectDb()

	migrations.Migrate()

	gin.SetMode(gin.TestMode)

	//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	//Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()

	if os.Getenv("REDIS_HOST") != "" {
		db.InitRedis(1)
	}

	// Add /example route via handler function to the gin instance
	handler := RouteHandler(app)

	// Create httpexpect instance
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	// Assert response
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("message", "pong")

	e.POST("/v1/user/register").
		WithForm(gin.H{
			"email": gofakeit.Email(),
			"password": "demo1admin",
			"name": "tutan khamun",
		}).
		Expect().
		Status(http.StatusOK)

	e.POST("/v1/user/login").
		WithForm(gin.H{"email": "6demo@localhost.com", "password": "1222"}).
		Expect().
		Status(http.StatusNotAcceptable)//StatusUnauthorized

	type LoginResponse struct {
		Message   string `json:"message"`
		Token  models.Token    `json:"token"`
		User models.User `json:"user"`
	}

	var loginResponse LoginResponse

	e.POST("/v1/user/login").
		WithForm(gin.H{"email": "demo@localhost.com", "password": "demo1admin"}).
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("message", "Successfully logged in").
		Decode(&loginResponse)

	/*e.GET("/v1/user/logout").
		Expect().
		Status(http.StatusOK)*/

	e.POST("/v1/token/refresh").
		WithForm(gin.H{"access_token": "expired_token"}).
		Expect().
		Status(http.StatusNotAcceptable)
}
