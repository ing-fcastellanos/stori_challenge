package services

import (
	"fintech-backend/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Balance godoc
// @Summary Get the balance of a user within a date range
// @Description Get the balance of a user by user_id within a date range
// @Produce json
// @Param user_id path int true "User ID"
// @Param from query string false "From date" example("2024-01-01T00:00:00Z")
// @Param to query string false "To date" example("2024-07-01T00:00:00Z")
// @Success 200 {object} models.BalanceResult
// @Failure 400 {string} string "Error al consultar el balance"
// @Failure 404 {string} string "Usuario no encontrado"
// @Router /users/{user_id}/balance [get]
func BalanceHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from != "" && to != "" {
		GetUserBalanceInRange(db, w, r)
	} else {
		GetUserBalance(db, w, r)
	}
}

// Función para obtener el balance total de un usuario
func GetUserBalance(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
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

	if creditsResult.TotalCredits == 0 && debitsResult.TotalDebits == 0 {
		http.Error(w, "Usuario no encontrado o sin transacciones", http.StatusNotFound)
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
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Obtener las fechas desde los parámetros de la URL
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

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

	if creditsResult.TotalCredits == 0 && debitsResult.TotalDebits == 0 {
		http.Error(w, "Usuario no encontrado o sin transacciones", http.StatusNotFound)
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
