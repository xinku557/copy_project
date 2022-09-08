package controllers

import (
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sheinko.tk/copy_project/models"
	"sheinko.tk/copy_project/utils/auth"

    "github.com/joho/godotenv"
)

type Handler struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (handler *Handler) Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}else{
		auth.TokenSecrect = os.Getenv("SECRET_KEY")
	}
	handler.initializeDatabase()
	handler.initializeRoutes()
}

func (handler *Handler) initializeDatabase() {
	log.Println("Initializing the database...")
	dsn := "admin:Admin@2022@tcp(localhost:3306)/copyp?charset=utf8&parseTime=True&loc=Local"
	var err error
	handler.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v\n", err)
	} else {
		log.Println("Database connection is successful")
	}
	if err := handler.DB.AutoMigrate(&models.User{}, &models.Post{}); err != nil {
		log.Fatalf("Error auto migration: %v", err)
	}
}

func (handler *Handler) Run(addr string) {
	log.Printf(" Server is listening on port %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler.Router))
}
