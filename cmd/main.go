package main

import (
	"log"
	"os"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	// create jwt handler
	privKey, err := os.ReadFile(os.Getenv("RSA_PRIVATE_KEY_PATH"))
	if err != nil {
		log.Fatalln(err)
	}
	pubKey, err := os.ReadFile(os.Getenv("RSA_PUBLIC_KEY_PATH"))
	if err != nil {
		log.Fatalln(err)
	}
	jwtHandler := handler.NewJWT(pubKey, privKey)

	opts := handler.NewServerOptions{
		Repository: repo,
		JWT:        jwtHandler,
	}
	return handler.NewServer(opts)
}
