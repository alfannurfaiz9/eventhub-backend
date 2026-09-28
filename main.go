package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/config"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err.Error())
		return
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

	r := gin.Default()

	router.InitMainRouter(r, pool)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
