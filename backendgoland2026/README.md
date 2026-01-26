# Backend Goland 2026

Backend API desarrollado en Go con Clean Architecture.

## 🏗️ Arquitectura

### Capas

1. **Domain**: Entidades y interfaces de repositorios
2. **Application**: Casos de uso (business logic)
3. **Infrastructure**: Implementación de persistencia y servicios externos
4. **Interfaces**: HTTP handlers y DTOs

## 📋 Endpoints

### POST /api/login
Autentica un usuario y obtiene los desafíos del día.

**Request:**
```json
{
  "email": "usuario@example.com",
  "password": "contraseña"
}
```

**Response:**
```json
{
  "email": "usuario@example.com",
  "challenges": [
    {
      "ticker": "GSHD",
      "target_from": "$110.00",
      "target_to": "$79.00",
      "company": "Goosehead Insurance",
      "action": "target lowered by",
      "brokerage": "",
      "rating_from": "Market Perform",
      "rating_to": "Market Perform",
      "time": "2025-11-26T00:30:05.005554348Z"
    }
  ],
  "total_pages": 5
}
```

### GET /api/challenges
Obtiene desafíos con paginación.

**Parameters:**
- `page` (int): Número de página (default: 1)
- `token` (string): Token de autenticación

**Response:**
```json
{
  "challenges": [...],
  "next_page": "PDSB",
  "total_pages": 5
}
```

## 🗄️ Base de Datos

Se conecta a CockroachDB usando PostgreSQL driver.

**Conexión:**
```
postgresql://user:pass@localhost:26257/mydb?sslmode=disable
```

Las tablas se crean automáticamente en la primera ejecución.

## 🚀 Ejecución

```bash
cd backendgoland2026
go run cmd/api/main.go
```

## 📦 Dependencias

- `github.com/gin-gonic/gin` - Framework HTTP
- `github.com/gin-contrib/cors` - CORS middleware
- `github.com/lib/pq` - PostgreSQL driver

## 🔄 Flujo de Autenticación

1. Usuario envía credenciales
2. Backend consulta BD por datos del día
3. Si no hay datos, consulta API de KarenAI
4. Guarda datos y cursor en BD
5. Retorna datos al cliente

## 📝 DTOs

### LoginRequest
```go
type LoginRequest struct {
  Email    string `json:"email" binding:"required,email"`
  Password string `json:"password" binding:"required"`
}
```

### LoginResponse
```go
type LoginResponse struct {
  Email      string
  Challenges []ChallengeDTO
  NextPage   string
  TotalPages int
}
```

### ChallengeDTO
```go
type ChallengeDTO struct {
  Ticker     string
  TargetFrom string
  TargetTo   string
  Company    string
  Action     string
  Brokerage  string
  RatingFrom string
  RatingTo   string
  Time       time.Time
}
```
