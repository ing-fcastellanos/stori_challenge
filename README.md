# FinTech Backend API

## Descripción

Esta es una API backend para la gestión de transacciones financieras en una aplicación de tecnología financiera (FinTech). Permite cargar archivos CSV con registros de transacciones y obtener información sobre el balance de los usuarios.

## Requisitos

- Docker
- Docker Compose
- Go (si deseas compilar localmente)

## Instalación

### 1. Clonar el repositorio

Primero, clona el repositorio de la API:

```bash
git clone https://github.com/ing-fcastellanos/stori_challenge.git
cd stori_challenge
```

### 2. Configuración del archivo `.env`

Crea un archivo `.env` en la raíz del proyecto con las siguientes variables de entorno para configurar la base de datos:

```dotenv
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=fintech
DB_SSLMODE=disable
```

### 3. Instalar dependencias

Si estás ejecutando la aplicación sin Docker, necesitarás instalar las dependencias de Go:

```bash
go mod tidy
```

### 4. Iniciar el proyecto con Docker

Si estás utilizando Docker, simplemente ejecuta el siguiente comando para construir e iniciar los contenedores:

```bash
docker-compose up --build
```

Esto construirá la imagen de Docker para la API y levantará el contenedor de PostgreSQL.

### 5. Iniciar el servidor localmente (si no usas Docker)

Si prefieres no usar Docker, puedes iniciar el servidor localmente ejecutando:

```bash
go run main.go
```

### 6. Acceder a la API

La API estará disponible en `http://localhost:8080`.

---

## Endpoints

### **POST /migrate**

Este endpoint recibe un archivo CSV con las transacciones y las guarda en la base de datos.

#### Ejemplo de solicitud:

```bash
curl -X POST -F "file=@/path/to/transactions.csv" http://localhost:8080/migrate
```

#### Respuesta exitosa:

```json
{
  "message": "Migración completada con éxito"
}
```

#### Errores comunes:

- **400 Bad Request**: Si el archivo no es válido o tiene un formato incorrecto.

### **GET /users/{user_id}/balance**

Este endpoint devuelve el balance de un usuario. Calcula el total de los créditos (transacciones con monto positivo) y los débitos (transacciones con monto negativo).

#### Ejemplo de solicitud:

```bash
curl http://localhost:8080/users/1/balance
```

#### Respuesta exitosa:

```json
{
  "balance": 25.21,
  "total_debits": 10,
  "total_credits": 15
}
```

#### Errores comunes:

- **404 Not Found**: Si no se encuentra el usuario o no tiene transacciones.

### **GET /users/{user_id}/balance?from=YYYY-MM-DDThh:mm:ssZ&to=YYYY-MM-DDThh:mm:ssZ**

Este endpoint devuelve el balance de un usuario en un rango de fechas específico.

#### Ejemplo de solicitud:

```bash
curl "http://localhost:8080/users/1/balance?from=2024-01-01T00:00:00Z&to=2024-07-01T00:00:00Z"
```

#### Respuesta exitosa:

```json
{
  "balance": 15.0,
  "total_debits": 5,
  "total_credits": 20
}
```

#### Errores comunes:

- **400 Bad Request**: Si el formato de la fecha no es válido.
- **404 Not Found**: Si no se encuentra el usuario o no tiene transacciones dentro del rango de fechas.

---

## Swagger UI

La documentación interactiva de la API está disponible a través de Swagger UI. Para acceder a la interfaz de Swagger, navega a:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Troubleshooting

### 1. **El archivo CSV no es procesado correctamente**

Si recibes un error 400 al intentar procesar el archivo CSV, asegúrate de que el archivo esté en el formato correcto y contenga las siguientes columnas:

- `id`
- `user_id`
- `amount`
- `datetime` (en formato `YYYY-MM-DDThh:mm:ssZ`)

### 2. **No puedo acceder a Swagger UI**

Si no puedes acceder a la interfaz de Swagger UI, verifica lo siguiente:

- Asegúrate de que el servidor esté corriendo en `http://localhost:8080`.
- Verifica que el contenedor de Docker esté corriendo correctamente.

Si el error persiste, revisa los logs del contenedor con:

```bash
docker-compose logs backend
```

### 3. **Errores de conexión a la base de datos**

Si no puedes conectar a la base de datos, asegúrate de que las credenciales en el archivo `.env` sean correctas. Además, asegúrate de que el contenedor de PostgreSQL esté levantado correctamente:

```bash
docker-compose ps
```

### 4. **Problemas con la migración**

Si las migraciones no se están aplicando correctamente, verifica los logs para identificar cualquier error en el proceso de migración.

Si las migraciones fallan, puedes intentar eliminar la base de datos y crearla de nuevo:

```bash
docker-compose down
docker-compose up --build
```
