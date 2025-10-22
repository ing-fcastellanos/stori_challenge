package services_test

import (
	"fintech-backend/services"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestBalanceHandler_WithRange(t *testing.T) {
	// Crear una conexión simulada a la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// Datos de prueba
	userID := "1"
	from := "2024-01-01T00:00:00Z"
	to := "2024-12-31T23:59:59Z"

	// Crear una solicitud HTTP con parámetros `from` y `to`
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%s/balance?from=%s&to=%s", userID, from, to), nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	// Crear un recorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Simular que BalanceHandler llama a GetUserBalanceInRange
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_credits`).
		WithArgs(userID, from, to).
		WillReturnRows(sqlmock.NewRows([]string{"total_credits"}).AddRow(200.0))

	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_debits`).
		WithArgs(userID, from, to).
		WillReturnRows(sqlmock.NewRows([]string{"total_debits"}).AddRow(-50.0))

	// Crear un router y registrar el handler
	router := mux.NewRouter()
	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.BalanceHandler(gormDB, w, r)
	}).Methods("GET")

	// Llamar al handler
	router.ServeHTTP(rr, req)

	// Verificar las expectativas de la base de datos
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar el código de estado
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verificar la respuesta
	expected := `{
		"balance": 150.00,
		"total_debits": -50.00,
		"total_credits": 200.00
	}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestBalanceHandler_WithoutRange(t *testing.T) {
	// Crear una conexión simulada a la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// Datos de prueba
	userID := "1"

	// Crear una solicitud HTTP sin parámetros `from` y `to`
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%s/balance", userID), nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	// Crear un recorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Simular que BalanceHandler llama a GetUserBalance
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_credits`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"total_credits"}).AddRow(100.0))

	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_debits`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"total_debits"}).AddRow(-25.0))

	// Crear un router y registrar el handler
	router := mux.NewRouter()
	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.BalanceHandler(gormDB, w, r)
	}).Methods("GET")

	// Llamar al handler
	router.ServeHTTP(rr, req)

	// Verificar las expectativas de la base de datos
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar el código de estado
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verificar la respuesta
	expected := `{
		"balance": 75.00,
		"total_debits": -25.00,
		"total_credits": 100.00
	}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestBalanceHandler_DBError(t *testing.T) {
	// Crear una conexión simulada a la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// Datos de prueba
	userID := "1"
	from := "2024-01-01T00:00:00Z"
	to := "2024-12-31T23:59:59Z"

	// Simular un error de base de datos en la consulta de créditos
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_credits`).
		WithArgs(userID, from, to).
		WillReturnError(fmt.Errorf("error en la base de datos"))

	// Crear la solicitud HTTP con parámetros `from` y `to`
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%s/balance?from=%s&to=%s", userID, from, to), nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	// Crear un recorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Crear un router y registrar el handler
	router := mux.NewRouter()
	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.BalanceHandler(gormDB, w, r)
	}).Methods("GET")

	// Llamar al handler
	router.ServeHTTP(rr, req)

	// Verificar las expectativas de la base de datos
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar el código de estado
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verificar el mensaje de error
	expected := `Error al consultar los créditos: error en la base de datos`
	assert.Contains(t, rr.Body.String(), expected)
}

func TestGetUserBalance(t *testing.T) {
	// Crear una conexión simulada a la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// Datos de prueba
	userID := "1"
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_credits`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"total_credits"}).AddRow(100.0))

	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_debits`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"total_debits"}).AddRow(-50.0))

	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%s/balance", userID), nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	// Crear un recorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Crear un router y registrar el handler
	router := mux.NewRouter()
	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.GetUserBalance(gormDB, w, r)
	}).Methods("GET")

	// Llamar al handler
	router.ServeHTTP(rr, req)

	// Verificar las expectativas de la base de datos
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar el código de estado
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verificar la respuesta
	expected := `{
		"balance": 50.00,
		"total_debits": -50.00,
		"total_credits": 100.00
	}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestGetUserBalance_UserNotFound(t *testing.T) {
	// Crear una conexión simulada a la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// Datos de prueba
	userID := "1"
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_credits`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"total_credits"}).AddRow(0))

	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_debits`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"total_debits"}).AddRow(0))

	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%s/balance", userID), nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	// Crear un recorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Crear un router y registrar el handler
	router := mux.NewRouter()
	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.GetUserBalance(gormDB, w, r)
	}).Methods("GET")

	// Llamar al handler
	router.ServeHTTP(rr, req)

	// Verificar las expectativas de la base de datos
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar el código de estado
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Verificar el mensaje de respuesta
	expected := `Usuario no encontrado o sin transacciones`
	assert.Contains(t, rr.Body.String(), expected)
}

func TestGetUserBalance_DBError(t *testing.T) {
	// Crear una conexión simulada a la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// Datos de prueba
	userID := "1"
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_credits`).
		WithArgs(userID).
		WillReturnError(fmt.Errorf("error en la base de datos"))

	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%s/balance", userID), nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	// Crear un recorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Crear un router y registrar el handler
	router := mux.NewRouter()
	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.GetUserBalance(gormDB, w, r)
	}).Methods("GET")

	// Llamar al handler
	router.ServeHTTP(rr, req)

	// Verificar las expectativas de la base de datos
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar el código de estado
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verificar el mensaje de error
	expected := `Error al consultar los créditos: error en la base de datos`
	assert.Contains(t, rr.Body.String(), expected)
}

func TestGetUserBalanceInRange(t *testing.T) {
	// Crear una conexión simulada a la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// Datos de prueba
	userID := "1"
	from := "2024-01-01T00:00:00Z"
	to := "2024-12-31T23:59:59Z"

	// Simular la consulta de créditos (transacciones con monto positivo) dentro del rango de fechas
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_credits`).
		WithArgs(userID, from, to).
		WillReturnRows(sqlmock.NewRows([]string{"total_credits"}).AddRow(200.0))

	// Simular la consulta de débitos (transacciones con monto negativo) dentro del rango de fechas
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_debits`).
		WithArgs(userID, from, to).
		WillReturnRows(sqlmock.NewRows([]string{"total_debits"}).AddRow(-50.0))

	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%s/balance?from=%s&to=%s", userID, from, to), nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	// Crear un recorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Crear un router y registrar el handler
	router := mux.NewRouter()
	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.GetUserBalanceInRange(gormDB, w, r)
	}).Methods("GET")

	// Llamar al handler
	router.ServeHTTP(rr, req)

	// Verificar las expectativas de la base de datos
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar el código de estado
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verificar la respuesta
	expected := `{
		"balance": 150.00,
		"total_debits": -50.00,
		"total_credits": 200.00
	}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestGetUserBalanceInRange_UserNotFound(t *testing.T) {
	// Crear una conexión simulada a la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// Datos de prueba
	userID := "1"
	from := "2024-01-01T00:00:00Z"
	to := "2024-12-31T23:59:59Z"

	// Simular la consulta de créditos (transacciones con monto positivo) dentro del rango de fechas
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_credits`).
		WithArgs(userID, from, to).
		WillReturnRows(sqlmock.NewRows([]string{"total_credits"}).AddRow(0))

	// Simular la consulta de débitos (transacciones con monto negativo) dentro del rango de fechas
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_debits`).
		WithArgs(userID, from, to).
		WillReturnRows(sqlmock.NewRows([]string{"total_debits"}).AddRow(0))

	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%s/balance?from=%s&to=%s", userID, from, to), nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	// Crear un recorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Crear un router y registrar el handler
	router := mux.NewRouter()
	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.GetUserBalanceInRange(gormDB, w, r)
	}).Methods("GET")

	// Llamar al handler
	router.ServeHTTP(rr, req)

	// Verificar las expectativas de la base de datos
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar el código de estado
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Verificar el mensaje de respuesta
	expected := `Usuario no encontrado o sin transacciones`
	assert.Contains(t, rr.Body.String(), expected)
}

func TestGetUserBalanceInRange_DBError(t *testing.T) {
	// Crear una conexión simulada a la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// Datos de prueba
	userID := "1"
	from := "2024-01-01T00:00:00Z"
	to := "2024-12-31T23:59:59Z"

	// Simular un error en la consulta de créditos (transacciones con monto positivo) dentro del rango de fechas
	mock.ExpectQuery(`SELECT SUM\(amount\) AS total_credits`).
		WithArgs(userID, from, to).
		WillReturnError(fmt.Errorf("error en la base de datos"))

	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%s/balance?from=%s&to=%s", userID, from, to), nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	// Crear un recorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Crear un router y registrar el handler
	router := mux.NewRouter()
	router.HandleFunc("/users/{user_id}/balance", func(w http.ResponseWriter, r *http.Request) {
		services.GetUserBalanceInRange(gormDB, w, r)
	}).Methods("GET")

	// Llamar al handler
	router.ServeHTTP(rr, req)

	// Verificar las expectativas de la base de datos
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar el código de estado
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verificar el mensaje de error
	expected := `Error al consultar los créditos: error en la base de datos`
	assert.Contains(t, rr.Body.String(), expected)
}
