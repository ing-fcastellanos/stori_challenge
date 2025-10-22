package main

import (
	"fintech-backend/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
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

func main() {
	// Inicializar el router
	router := mux.NewRouter()

	// Definir los endpoints
	// TODO: Cargar las variables de DB desde una interface e iniciarlizar con el proyecto
	router.HandleFunc("/migrate", func(w http.ResponseWriter, r *http.Request) {
		services.MigrateTransactions(db, w, r)
	}).Methods("POST")

	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.BalanceHandler(db, w, r)
	}).Methods("GET")

	// Arrancar el servidor
	log.Fatal(http.ListenAndServe(":8080", router))
}
