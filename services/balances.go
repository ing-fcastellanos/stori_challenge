package services

import (
	"fintech-backend/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func BalanceHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	log.Print(fmt.Sprintf("handdler: %s, %s", from, to))

	if from != "" && to != "" {
		GetUserBalanceInRange(db, w, r)
	} else {
		GetUserBalance(db, w, r)
	}
}

// Función para obtener el balance total de un usuario
func GetUserBalance(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("Regular Balances")

	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Consultar los créditos (transacciones con monto positivo)
	var creditsResult models.BalanceResult
	err := db.Table("transactions").
		Select("SUM(amount) AS total_credits").
		Where("user_id = ?", userID).
		Where("amount > 0").
		Scan(&creditsResult).
		Error

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al consultar los créditos: %v", err), http.StatusInternalServerError)
		return
	}

	// Consultar los débitos (transacciones con monto negativo)
	var debitsResult models.BalanceResult
	err = db.Table("transactions").
		Select("SUM(amount) AS total_debits").
		Where("user_id = ?", userID).
		Where("amount < 0").
		Scan(&debitsResult).
		Error

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al consultar los débitos: %v", err), http.StatusInternalServerError)
		return
	}

	// Calcular el balance final
	balance := creditsResult.TotalCredits + debitsResult.TotalDebits

	// Enviar respuesta con el balance
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"balance": %.2f,
		"total_debits": %.2f,
		"total_credits": %.2f
	}`, balance, debitsResult.TotalDebits, creditsResult.TotalCredits)))
}

// Función para obtener el balance de un usuario en un rango de fechas
func GetUserBalanceInRange(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("Balance in range")
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Obtener las fechas desde los parámetros de la URL
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	log.Print(fmt.Sprintf("%s, %s", from, to))

	// Consultar los créditos (transacciones con monto positivo) en el rango de fechas
	var creditsResult models.BalanceResult
	err := db.Table("transactions").
		Select("SUM(amount) AS total_credits").
		Where("user_id = ?", userID).
		Where("amount > 0").
		Where("timestamp BETWEEN ? AND ?", from, to).
		Scan(&creditsResult).
		Error

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al consultar los créditos: %v", err), http.StatusInternalServerError)
		return
	}

	// Consultar los débitos (transacciones con monto negativo) en el rango de fechas
	var debitsResult models.BalanceResult
	err = db.Table("transactions").
		Select("SUM(amount) AS total_debits").
		Where("user_id = ?", userID).
		Where("amount < 0").
		Where("timestamp BETWEEN ? AND ?", from, to).
		Scan(&debitsResult).
		Error

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al consultar los débitos: %v", err), http.StatusInternalServerError)
		return
	}

	// Calcular el balance final
	balance := creditsResult.TotalCredits + debitsResult.TotalDebits

	// Enviar respuesta con el balance en el rango de fechas
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"balance": %.2f,
		"total_debits": %.2f,
		"total_credits": %.2f
	}`, balance, debitsResult.TotalDebits, creditsResult.TotalCredits)))
}
