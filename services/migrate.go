package services

import (
	"encoding/csv"
	"fintech-backend/models"
	"fintech-backend/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DataMigrate godoc
// @Summary Upload transactions CSV and migrate to DB
// @Description Upload a CSV file with transaction data and migrate it to the database
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV File"
// @Success 200 {string} string "Migración completada con éxito"
// @Failure 400 {string} string "Error al recibir el archivo"
// @Failure 500 {string} string "Error al guardar la transacción"
// @Router /migrate [post]
func MigrateTransactions(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Leer el archivo CSV desde el request
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error al recibir el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parsear el archivo CSV
	reader := csv.NewReader(file)
	// Leer la primera línea (encabezado) para asegurarnos de que las columnas sean las correctas
	headers, err := reader.Read()
	if err != nil {
		http.Error(w, "Error al leer el archivo CSV", http.StatusInternalServerError)
		return
	}

	// Validar que el CSV tenga las columnas necesarias
	expectedHeaders := []string{"id", "user_id", "amount", "datetime"}
	if !isValidCSVHeader(headers, expectedHeaders) {
		http.Error(w, "El archivo CSV tiene un formato incorrecto", http.StatusBadRequest)
		return
	}

	// Procesar cada línea del CSV
	var transactions []models.Transaction
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Convertir los datos a la estructura correspondiente
		timestamp, err := time.Parse(time.RFC3339, record[3])
		if err != nil {
			http.Error(w, "Error de formato en la fecha", http.StatusBadRequest)
			return
		}

		// Crear la transacción y agregarla a la lista
		transaction := models.Transaction{
			LegacyId:  utils.ParseInt(record[0]),
			UserID:    utils.ParseInt(record[1]),
			Amount:    utils.ParseFloat(record[2]),
			Timestamp: timestamp,
		}

		transactions = append(transactions, transaction)
	}

	// Insertar las transacciones en la base de datos
	for _, txn := range transactions {
		if err := db.Create(&txn).Error; err != nil {
			log.Println("Error al guardar la transacción:", err)
			http.Error(w, fmt.Sprintf("Error al guardar la transacción: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Enviar respuesta de éxito
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Migración completada con éxito"))
}

// Verificar que las columnas del archivo CSV sean las esperadas
func isValidCSVHeader(headers, expectedHeaders []string) bool {
	if len(headers) != len(expectedHeaders) {
		return false
	}
	for i := range headers {
		if headers[i] != expectedHeaders[i] {
			return false
		}
	}
	return true
}
