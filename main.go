package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/config"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/middleware"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           			Eventhub Backend
// @version         			1.0
// @description     			API docs for eventhub
// @host      					localhost:9000
// @BasePath  					/
// @securityDefinitions.apikey	BearerToken
// @in							header
// @name						Authorization
// @description					Bearer Token used as identity for accessing backend
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err.Error())
		// return
	}

	pdb := config.NewPsqlDb(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME"))
	pool, err := pdb.Connect()
	if err != nil {
		log.Println("Cannot connect to database \n", err.Error())
		return
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		log.Println("Database is not ready \n", err.Error())
	}

	rdbConf := config.NewRdb(os.Getenv("REDIS_USER"), os.Getenv("REDIS_PASS"), os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	rdb := rdbConf.Connect()
	defer rdb.Close()

	r := gin.Default()
	r.Use(middleware.Cors)
	router.InitMainRouter(r, pool, rdb)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
