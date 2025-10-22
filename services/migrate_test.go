package services_test

import (
	"bytes"
	"fintech-backend/services"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestMigrateTransactions_Success(t *testing.T) {
	// Crear mock de base de datos
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

	// Contenido del CSV válido
	csvContent := `id,user_id,amount,datetime
1,101,150.50,2024-01-15T10:30:00Z
2,102,-25.75,2024-01-16T14:45:00Z
3,101,300.00,2024-01-17T09:00:00Z`

	// Crear multipart form con el archivo CSV
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "transactions.csv")
	if err != nil {
		t.Fatalf("Error al crear form file: %v", err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	// Expectativas de la base de datos (3 inserciones)
	// GORM usa Query con RETURNING para obtener el ID generado
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "transactions"`).
		WithArgs(1, 101, 150.50, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "transactions"`).
		WithArgs(2, 102, -25.75, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "transactions"`).
		WithArgs(3, 101, 300.00, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
	mock.ExpectCommit()

	// Crear request HTTP
	req, err := http.NewRequest("POST", "/transactions/migrate", body)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Crear recorder
	rr := httptest.NewRecorder()

	// Crear router
	router := mux.NewRouter()
	router.HandleFunc("/transactions/migrate", func(w http.ResponseWriter, r *http.Request) {
		services.MigrateTransactions(gormDB, w, r)
	}).Methods("POST")

	// Ejecutar handler
	router.ServeHTTP(rr, req)

	// Verificar expectativas
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar código de estado
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verificar mensaje de éxito
	assert.Contains(t, rr.Body.String(), "Migración completada con éxito")
}

func TestMigrateTransactions_NoFileProvided(t *testing.T) {
	// Crear mock de base de datos
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// Crear request sin archivo
	req, err := http.NewRequest("POST", "/transactions/migrate", nil)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}

	// Crear recorder
	rr := httptest.NewRecorder()

	// Crear router
	router := mux.NewRouter()
	router.HandleFunc("/transactions/migrate", func(w http.ResponseWriter, r *http.Request) {
		services.MigrateTransactions(gormDB, w, r)
	}).Methods("POST")

	// Ejecutar handler
	router.ServeHTTP(rr, req)

	// Verificar código de estado
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Verificar mensaje de error
	assert.Contains(t, rr.Body.String(), "Error al recibir el archivo")
}

func TestMigrateTransactions_EmptyCSV(t *testing.T) {
	// Crear mock de base de datos
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// CSV vacío
	csvContent := ``

	// Crear multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "transactions.csv")
	if err != nil {
		t.Fatalf("Error al crear form file: %v", err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	// Crear request HTTP
	req, err := http.NewRequest("POST", "/transactions/migrate", body)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Crear recorder
	rr := httptest.NewRecorder()

	// Crear router
	router := mux.NewRouter()
	router.HandleFunc("/transactions/migrate", func(w http.ResponseWriter, r *http.Request) {
		services.MigrateTransactions(gormDB, w, r)
	}).Methods("POST")

	// Ejecutar handler
	router.ServeHTTP(rr, req)

	// Verificar código de estado
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verificar mensaje de error
	assert.Contains(t, rr.Body.String(), "Error al leer el archivo CSV")
}

func TestMigrateTransactions_InvalidDateFormat(t *testing.T) {
	// Crear mock de base de datos
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// CSV con fecha inválida
	csvContent := `id,user_id,amount,datetime
1,101,150.50,2024/01/15 10:30:00`

	// Crear multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "transactions.csv")
	if err != nil {
		t.Fatalf("Error al crear form file: %v", err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	// Crear request HTTP
	req, err := http.NewRequest("POST", "/transactions/migrate", body)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Crear recorder
	rr := httptest.NewRecorder()

	// Crear router
	router := mux.NewRouter()
	router.HandleFunc("/transactions/migrate", func(w http.ResponseWriter, r *http.Request) {
		services.MigrateTransactions(gormDB, w, r)
	}).Methods("POST")

	// Ejecutar handler
	router.ServeHTTP(rr, req)

	// Verificar código de estado
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Verificar mensaje de error
	assert.Contains(t, rr.Body.String(), "Error de formato en la fecha")
}

func TestMigrateTransactions_DatabaseError(t *testing.T) {
	// Crear mock de base de datos
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

	// CSV válido
	csvContent := `id,user_id,amount,datetime
1,101,150.50,2024-01-15T10:30:00Z`

	// Crear multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "transactions.csv")
	if err != nil {
		t.Fatalf("Error al crear form file: %v", err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	// Simular error en la base de datos
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "transactions"`).
		WithArgs(1, 101, 150.50, sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("error de conexión a la base de datos"))
	mock.ExpectRollback()

	// Crear request HTTP
	req, err := http.NewRequest("POST", "/transactions/migrate", body)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Crear recorder
	rr := httptest.NewRecorder()

	// Crear router
	router := mux.NewRouter()
	router.HandleFunc("/transactions/migrate", func(w http.ResponseWriter, r *http.Request) {
		services.MigrateTransactions(gormDB, w, r)
	}).Methods("POST")

	// Ejecutar handler
	router.ServeHTTP(rr, req)

	// Verificar expectativas
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Las expectativas de la base de datos no fueron cumplidas: %v", err)
	}

	// Verificar código de estado
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verificar mensaje de error
	assert.Contains(t, rr.Body.String(), "Error al guardar la transacción")
}

func TestMigrateTransactions_OnlyHeaders(t *testing.T) {
	// Crear mock de base de datos
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer db.Close()

	// Crear una instancia de GORM usando el mock de SQL
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Error al abrir la conexión con GORM: %v", err)
	}

	// CSV solo con encabezados
	csvContent := `id,user_id,amount,datetime`

	// Crear multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "transactions.csv")
	if err != nil {
		t.Fatalf("Error al crear form file: %v", err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	// Crear request HTTP
	req, err := http.NewRequest("POST", "/transactions/migrate", body)
	if err != nil {
		t.Fatalf("Error al crear solicitud: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Crear recorder
	rr := httptest.NewRecorder()

	// Crear router
	router := mux.NewRouter()
	router.HandleFunc("/transactions/migrate", func(w http.ResponseWriter, r *http.Request) {
		services.MigrateTransactions(gormDB, w, r)
	}).Methods("POST")

	// Ejecutar handler
	router.ServeHTTP(rr, req)

	// Verificar código de estado
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verificar mensaje de éxito
	assert.Contains(t, rr.Body.String(), "Migración completada con éxito")
}

func TestIsValidCSVHeader_ValidHeaders(t *testing.T) {
	// Encabezados válidos que coinciden exactamente
	headers := []string{"id", "user_id", "amount", "datetime"}
	expectedHeaders := []string{"id", "user_id", "amount", "datetime"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.True(t, result, "Los encabezados válidos deberían retornar true")
}

func TestIsValidCSVHeader_InvalidHeaders_DifferentNames(t *testing.T) {
	// Encabezados con nombres diferentes
	headers := []string{"id", "usuario", "monto", "fecha"}
	expectedHeaders := []string{"id", "user_id", "amount", "datetime"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.False(t, result, "Los encabezados con nombres diferentes deberían retornar false")
}

func TestIsValidCSVHeader_InvalidHeaders_DifferentOrder(t *testing.T) {
	// Encabezados en diferente orden
	headers := []string{"user_id", "id", "amount", "datetime"}
	expectedHeaders := []string{"id", "user_id", "amount", "datetime"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.False(t, result, "Los encabezados en diferente orden deberían retornar false")
}

func TestIsValidCSVHeader_InvalidHeaders_FewerColumns(t *testing.T) {
	// Menos columnas de las esperadas
	headers := []string{"id", "user_id", "amount"}
	expectedHeaders := []string{"id", "user_id", "amount", "datetime"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.False(t, result, "Los encabezados con menos columnas deberían retornar false")
}

func TestIsValidCSVHeader_InvalidHeaders_MoreColumns(t *testing.T) {
	// Más columnas de las esperadas
	headers := []string{"id", "user_id", "amount", "datetime", "extra_column"}
	expectedHeaders := []string{"id", "user_id", "amount", "datetime"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.False(t, result, "Los encabezados con más columnas deberían retornar false")
}

func TestIsValidCSVHeader_EmptyHeaders(t *testing.T) {
	// Ambos arreglos vacíos
	headers := []string{}
	expectedHeaders := []string{}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.True(t, result, "Dos arreglos vacíos deberían retornar true")
}

func TestIsValidCSVHeader_EmptyHeadersWithExpected(t *testing.T) {
	// Headers vacío pero se esperan columnas
	headers := []string{}
	expectedHeaders := []string{"id", "user_id", "amount", "datetime"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.False(t, result, "Headers vacío con expected no vacío debería retornar false")
}

func TestIsValidCSVHeader_CaseSensitive(t *testing.T) {
	// Verificar que es case sensitive
	headers := []string{"ID", "USER_ID", "AMOUNT", "DATETIME"}
	expectedHeaders := []string{"id", "user_id", "amount", "datetime"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.False(t, result, "Los encabezados con diferente capitalización deberían retornar false")
}

func TestIsValidCSVHeader_WithWhitespace(t *testing.T) {
	// Encabezados con espacios en blanco
	headers := []string{"id ", " user_id", "amount ", " datetime"}
	expectedHeaders := []string{"id", "user_id", "amount", "datetime"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.False(t, result, "Los encabezados con espacios en blanco deberían retornar false")
}

func TestIsValidCSVHeader_PartialMatch(t *testing.T) {
	// Solo las primeras columnas coinciden
	headers := []string{"id", "user_id", "monto", "fecha"}
	expectedHeaders := []string{"id", "user_id", "amount", "datetime"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.False(t, result, "Una coincidencia parcial debería retornar false")
}

func TestIsValidCSVHeader_SingleColumn_Valid(t *testing.T) {
	// Una sola columna que coincide
	headers := []string{"id"}
	expectedHeaders := []string{"id"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.True(t, result, "Una sola columna válida debería retornar true")
}

func TestIsValidCSVHeader_SingleColumn_Invalid(t *testing.T) {
	// Una sola columna que no coincide
	headers := []string{"identifier"}
	expectedHeaders := []string{"id"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.False(t, result, "Una sola columna inválida debería retornar false")
}

func TestIsValidCSVHeader_SpecialCharacters(t *testing.T) {
	// Encabezados con caracteres especiales
	headers := []string{"id", "user-id", "amount$", "datetime!"}
	expectedHeaders := []string{"id", "user-id", "amount$", "datetime!"}

	result := services.IsValidCSVHeader(headers, expectedHeaders)

	assert.True(t, result, "Los encabezados con caracteres especiales que coinciden deberían retornar true")
}
