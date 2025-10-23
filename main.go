package main

import (
	"fintech-backend/services"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "fintech-backend/docs"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	httpSwagger "github.com/swaggo/http-swagger"
)

var db *gorm.DB

func init() {
	//TODO: mover esta función hacia su propio fichero de clase DB
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// Obtener las variables de entorno
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Conectarse a la base de datos PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	// Ejecutar las migraciones con Goose
	err = goose.Up(db.DB(), "./migrations")
	if err != nil {
		log.Fatal("Error al ejecutar migraciones:", err)
	}
}

// @title FinTech Backend API
// @version 1.0
// @description This is a simple API for handling financial transactions
// @termsOfService https://example.com/terms/
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	// Inicializar el router
	// TODO: Cargar las variables de DB desde una interface e iniciarlizar con el proyecto
	router := mux.NewRouter()

	router.HandleFunc("/migrate", func(w http.ResponseWriter, r *http.Request) {
		services.MigrateTransactions(db, w, r)
	}).Methods("POST")

	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.BalanceHandler(db, w, r)
	}).Methods("GET")

	// Documentación Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Arrancar el servidor
	log.Fatal(http.ListenAndServe(":8080", router))
}
