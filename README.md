# 🚀 MELI Product API - Go Implementation

API REST de productos estilo MercadoLibre implementada en **Go** con arquitectura hexagonal y patrón BFF (Backend For Frontend).

[![Go](https://img.shields.io/badge/Go-1.21-00ADD8.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## ✨ Características

### Funcionales
- ✅ **Detalles completos de producto** con agregación de múltiples fuentes
- ✅ **Búsqueda de productos** con paginación
- ✅ **Productos relacionados** por categoría
- ✅ **Reviews y calificaciones** con estadísticas
- ✅ **Preguntas y respuestas** de usuarios
- ✅ **Cálculo de envío** con lógica de envío gratis

### Técnicas
- 🏗️ **Arquitectura Hexagonal** (Ports & Adapters)
- 🔄 **Patrón BFF** para agregación de datos
- ⚡ **Procesamiento asíncrono** con Goroutines
- 📦 **JSON como base de datos** (simulación de microservicios)
- 🐳 **Dockerizado** con multi-stage builds
- 📝 **Logging estructurado** con slog
- 🔒 **Non-root container** (security best practice)
- 🚀 **Alta performance** con Go concurrency

---

## 📦 Requisitos

- **Go 1.21+**
- **Docker** (opcional)
- **Make** (opcional, para comandos simplificados)

---

## 🚀 Quick Start

### Opción 1: Ejecutar localmente
```bash
# Clonar repositorio
git clone https://github.com/tu-usuario/meli-product-api-go.git
cd meli-product-api-go

# Instalar dependencias
go mod download

# Ejecutar aplicación
go run cmd/api/main.go
```

### Opción 2: Docker
```bash
# Construir imagen
docker build -t meli-product-api-go .

# Ejecutar contenedor
docker run -p 8080:8080 meli-product-api-go
```

### Opción 3: Docker Compose
```bash
# Iniciar servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down
```

### Opción 4: Makefile (Recomendado)
```bash
# Ver comandos disponibles
make help

# Ejecutar localmente
make run

# Ejecutar con hot reload
make dev

# Construir y ejecutar con Docker
make docker-build
make docker-run

# Docker Compose
make docker-compose-up
```

---

## 📋 Estructura del Proyecto
```
meli-product-api-go/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point
├── internal/
│   ├── domain/                     # Domain Layer
│   │   ├── model/                  # Entities
│   │   └── port/                   # Ports (Interfaces)
│   ├── application/                # Application Layer
│   │   └── service/                # Use Cases
│   └── infrastructure/             # Infrastructure Layer
│       ├── adapter/
│       │   ├── http/               # HTTP Handlers
│       │   └── repository/         # Data Access
│       ├── config/                 # Configuration
│       └── router/                 # HTTP Router
├── data/                           # JSON Database
│   ├── products.json
│   ├── sellers.json
│   ├── reviews.json
│   └── questions.json
├── Dockerfile                      # Production image
├── docker-compose.yml              # Docker Compose config
├── Makefile                        # Build automation
├── go.mod                          # Go modules
└── README.md                       # This file
```

---

## 🔌 API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### 1. Obtener Detalles de Producto
```bash
GET /products/{id}

# Ejemplo
curl http://localhost:8080/api/v1/products/MLA123456
```

**Respuesta 200 OK:**
```json
{
  "product": {
    "id": "MLA123456",
    "title": "iPhone 14 Pro Max 256GB",
    "price": 899999.00,
    ...
  },
  "seller": {...},
  "shipping": {...},
  "reviews": {...},
  "questions": [...],
  "related_products": [...]
}
```

### 2. Buscar Productos
```bash
GET /products/search?q={query}&limit={limit}&offset={offset}

# Ejemplo
curl "http://localhost:8080/api/v1/products/search?q=iphone&limit=5&offset=0"
```

**Respuesta 200 OK:**
```json
{
  "query": "iphone",
  "total_results": 2,
  "limit": 5,
  "offset": 0,
  "results": [...]
}
```

### 3. Health Check
```bash
GET /health

# Ejemplo
curl http://localhost:8080/health
```

---

## 🧪 Testing
```bash
# Ejecutar tests
make test

# Tests con coverage
make test-coverage

# Benchmarks
make benchmark

# Linter
make lint

# Todas las verificaciones
make check
```

---

## 🐳 Docker

### Características del Dockerfile

- ✅ **Multi-stage build** (~15MB final image)
- ✅ **Alpine Linux** base image
- ✅ **Non-root user** (seguridad)
- ✅ **Health check** integrado
- ✅ **Optimizado** con CGO_ENABLED=0

### Comandos Docker
```bash
# Build
make docker-build

# Run
make docker-run

# Logs
make docker-logs

# Stop
make docker-stop

# Clean
make docker-clean
```

---

## 📊 Performance

### Go vs Java Comparison

| Métrica | Go | Java (Spring Boot) |
|---------|----|--------------------|
| **Startup time** | ~100ms | ~5-10s |
| **Memory usage** | ~20MB | ~200-500MB |
| **Image size** | ~15MB | ~200MB |
| **Request latency** | ~10ms | ~50ms |
| **Concurrency** | Goroutines (millions) | Threads (thousands) |

### Benchmarks
```bash
# Load test (requires apache bench)
make load-test

# Manual benchmark
ab -n 10000 -c 100 http://localhost:8080/api/v1/products/MLA123456
```

---

## 🛠️ Desarrollo

### Hot Reload (Air)
```bash
# Instalar Air
go install github.com/cosmtrek/air@latest

# Ejecutar con hot reload
make dev
```

### Instalar herramientas de desarrollo
```bash
make install-tools
```

Instala:
- Air (hot reload)
- golangci-lint (linter)
- swag (Swagger generator)

---

## 🎯 Decisiones de Diseño

### ¿Por qué Go?

**Ventajas:**
- ✅ **Performance** - Compilado, sin VM
- ✅ **Concurrencia nativa** - Goroutines y channels
- ✅ **Deployment simple** - Single binary
- ✅ **Memory footprint** - 10-20x menor que Java
- ✅ **Startup rápido** - Ideal para containers
- ✅ **Usado por MELI** - Arquitectura real

### Arquitectura Hexagonal

- **Domain** - Reglas de negocio puras
- **Application** - Casos de uso
- **Infrastructure** - Adaptadores (HTTP, JSON, DB)

### Concurrencia
```go
// Llamadas paralelas con Goroutines
var wg sync.WaitGroup
wg.Add(4)

go func() { defer wg.Done(); fetchSeller() }()
go func() { defer wg.Done(); fetchReviews() }()
go func() { defer wg.Done(); fetchQuestions() }()
go func() { defer wg.Done(); fetchRelated() }()

wg.Wait()
```

---

## 🚀 Migración a Producción

### PostgreSQL
```go
// Cambiar de JSON a PostgreSQL
productRepo := postgres.NewProductRepository(db)
```

### Redis Cache
```go
// Agregar cache layer
cachedRepo := cache.NewCachedRepository(productRepo, redisClient)
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: meli-product-api
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: api
        image: meli-product-api-go:latest
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
```

---

## 📚 Recursos

- [Go Documentation](https://go.dev/doc/)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [BFF Pattern](https://samnewman.io/patterns/architectural/bff/)
- [Effective Go](https://go.dev/doc/effective_go)

---

## 👤 Autor

**Juan** - Senior Backend Developer
- 8+ años de experiencia
- Java/Spring Boot + Go
- Arquitectura de microservicios

---

## 📄 Licencia

Este proyecto es de uso educativo y evaluación técnica.

---

**¡Gracias por revisar este proyecto!** 🚀
