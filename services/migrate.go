package services

import (
	"encoding/csv"
	"fintech-backend/models"
	"fintech-backend/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func MigrateTransactions(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Leer el archivo CSV desde el request
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error al recibir el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parsear el archivo CSV y omite la primera linea (encabezados)
	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		http.Error(w, "Error al leer el archivo CSV", http.StatusInternalServerError)
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
			http.Error(w, fmt.Sprintf("Error al guardar la transacción: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Enviar respuesta de éxito
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Migración completada con éxito"))
}
