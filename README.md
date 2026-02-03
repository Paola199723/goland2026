# Backend Goland 2026

## ejecucion rapida de backend y frontend 
Puedes ejecutar todos los procesos en cponsola ejecutando lo siguinte comando 
```go
bash START.sh
```

Backend API desarrollado en Go con Clean Architecture.


## docuemntación de ejecición del proyecto e instalacion de dependencias en el backend
verifica si tienes instalado crocroachdb con el siguiente comando 
```go
cockroach version
```
en caso de no tenerlo instalado puedes ejecura el siguiente comando para mac 
```go
brew install cockroach 
```
y para windows Descarga CockroachDB
Ve a la página oficial de CockroachDB,
Descarga el ZIP para Windows (64-bit)
Dentro verás cockroach.exe.
Abre Configuración del sistema
Variables de entorno
En Path → Editar → Nuevo
En PowerShell o CMD verifica que este instalado con 
```go
cockroach version
```

1. ejecutar desde la consola  localmente crockroachdb con el comando en mac o windows y lo dejas corriendo en la consola sin cerrarla
```go
cockroach start-single-node --insecure
```

2. ahora abre otra consola y ejecuta el comando 
``` go 
go run cmd/api/main.go            
```

# Frontend Goland 2026

Frontend moderno desarrollado con Vue 3, TypeScript, Pinia y Tailwind CSS.

## Características

- ✅ Autenticación de usuarios
- ✅ Visualización de challenges en tabla
- ✅ Paginación
- ✅ Integración con API backend
- ✅ Diseño responsivo con Tailwind CSS
- ✅ Gestión de estado con Pinia
- ✅ TypeScript para type safety

## Requisitos

- Node.js 16+
- npm o yarn

## Instalación

```bash
cd frontEndGoland2026
npm install
```

## Desarrollo

```bash
npm run dev
```

La aplicación estará disponible en `http://localhost:5173`

## Build para producción

```bash
npm run build
```

## Estructura del Proyecto frontEnd

```
src/
├── pages/              # Páginas principales
│   ├── LoginPage.vue   # Página de login
│   └── DashboardPage.vue # Panel principal
├── components/         # Componentes reutilizables
├── services/          # Servicios API
├── stores/            # Store de Pinia
├── App.vue            # Componente raíz
├── main.ts            # Punto de entrada
└── style.css          # Estilos globales
```

## Variables de Entorno

Configurar en un archivo `.env` si es necesario:

```
VITE_API_BASE_URL=http://localhost:8080
```



## 🏗️ Arquitectura del backend

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

Flujo de Consulta de lista en backend y frontend
![Flujo del backend](flujo_lista.png)
Flujo de creación de recomendación en backend y frontend
![Flujo del backend](recomendacion.png)