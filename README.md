# Technical Challenge - Talsory Recruiting

Solución al desafío técnico de dos APIs RESTful comunicadas por HTTP, con factorización QR en Go y cálculo de estadísticas en Node.js. Incluye Docker, JWT, pruebas unitarias y un frontend desarrollado con Angular 22 y Tailwind CSS.

## 📋 Tabla de Contenidos

1. [Descripción del Proyecto](#descripción-del-proyecto)
2. [Arquitectura](#arquitectura)
3. [Tecnologías Utilizadas](#tecnologías-utilizadas)
4. [Justificación de Decisiones Técnicas](#justificación-de-decisiones-técnicas)
5. [Estructura de Carpetas](#estructura-de-carpetas)
6. [Requisitos de Instalación](#requisitos-de-instalación)
7. [Ejecución Local (sin Docker)](#ejecución-local-sin-docker)
8. [Ejecución con Docker](#ejecución-con-docker)
9. [Endpoints y Ejemplos](#endpoints-y-ejemplos)
10. [Pruebas](#pruebas)
11. [Despliegue en la Nube](#despliegue-en-la-nube)
12. [Limitaciones Conocidas](#limitaciones-conocidas)

---

## Descripción del Proyecto

El sistema consta de dos microservicios backend y un frontend:

- **API Go (Fiber):** Recibe una matriz rectangular, calcula su factorización QR (Q y R) mediante el algoritmo de Householder y envía las matrices resultantes a la API de Node.js.
- **API Node.js (Express):** Recibe las matrices Q y R, calcula estadísticas (máximo, mínimo, promedio, suma total y si es diagonal) y devuelve la respuesta a la API de Go, que consolida todo y responde al cliente final.
- **Frontend (Angular 22 + Tailwind CSS):** Proporciona la interfaz gráfica para interactuar con el sistema, consumir los servicios backend y visualizar los resultados de la factorización QR y sus estadísticas.

Las APIs están protegidas por **JWT**, excepto el endpoint de login y los health checks.

## Arquitectura

```text
                       [ Usuario ]
                            |
                            v
                  [ Frontend Angular ]
                       Cloudflare
                            |
                            | HTTP / HTTPS + JWT
                            v
                     [ API Go :8080 ]
                            |
                            | 1. Valida matriz
                            | 2. Calcula factorización QR
                            |    (Householder)
                            | 3. Envía Q y R con JWT
                            v
                    [ API Node :3000 ]
                            |
                            | 4. Calcula estadísticas
                            |    (máximo, mínimo, promedio,
                            |     suma e indicador diagonal)
                            v
                        [ API Go ]
                            |
                            | 5. Consolida resultados
                            v
                  [ Frontend Angular ]
                            |
                            v
                       [ Usuario ]
```

## Tecnologías Utilizadas

| Componente | Tecnología |
|---|---|
| API 1 | Go 1.22 + Fiber v2 |
| API 2 | Node.js 20 + Express v4 |
| Autenticación | JWT (HS256) |
| Contenedores | Docker + Docker Compose |
| Pruebas (Go) | testing nativo |
| Pruebas (Node) | Jest + Supertest |
| Frontend | Angular 22 + Tailwind CSS |

## Justificación de Decisiones Técnicas

### 1. Factorización QR vs. Rotación (Inconsistencia del enunciado)

El documento original menciona "rotación" en la sección de arquitectura, pero en la sección de funcionalidad requerida exige explícitamente la factorización QR.

Se optó por implementar QR porque:

- Es el requerimiento matemático explícito y más relevante para la evaluación.
- La rotación es una operación trivial en comparación con la factorización QR.
- QR demuestra dominio de álgebra lineal y estabilidad numérica.

### 2. Algoritmo de Householder

Se implementó Householder en lugar de Gram-Schmidt porque:

- Es numéricamente más estable y evita la pérdida de ortogonalidad.
- Maneja mejor matrices mal condicionadas.
- Es el estándar en librerías como LAPACK.

### 3. QR Reducida (Thin QR)

Para una matriz A de dimensiones m × n, donde m ≥ n, se devuelve:

- Q de dimensiones m × n.
- R de dimensiones n × n.

Esta variante es más eficiente y práctica para el cálculo de estadísticas.

### 4. JWT

Se implementó autenticación JWT para proteger los endpoints.

El login se realiza en Go, que emite el token. Node.js valida el mismo token utilizando la misma clave secreta (`JWT_SECRET`).

### 5. Comunicación HTTP con Timeouts

El cliente HTTP de Go tiene un timeout de 10 segundos para evitar bloqueos si Node.js no responde.

Además, se usa `context.Context` para permitir la cancelación desde el cliente original.

### 6. Manejo de Matrices Rectangulares y Rango Deficiente

Se requiere m ≥ n (filas ≥ columnas) para QR estándar.

Householder maneja matrices de rango deficiente (R tendrá ceros o valores cercanos a cero en la diagonal), pero no se considera un error; se documenta como propiedad matemática.

## Estructura de Carpetas

```text
coding-challenge/
├── api-go/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handler/       # Handlers HTTP (QR, Auth, Middleware JWT)
│   │   ├── service/       # Lógica de negocio (QRService, AuthService)
│   │   ├── math/          # Algoritmo QR (Householder)
│   │   ├── client/        # Cliente HTTP hacia Node.js
│   │   └── model/         # Tipos y validaciones
│   ├── tests/             # Pruebas unitarias (QR)
│   ├── Dockerfile
│   └── go.mod
├── api-node/
│   ├── src/
│   │   ├── routes/        # Definición de endpoints
│   │   ├── controllers/   # Controladores
│   │   ├── services/      # Lógica de estadísticas
│   │   ├── middlewares/   # Middleware JWT
│   │   └── app.js         # Configuración de Express
│   ├── tests/             # Pruebas unitarias (Jest)
│   ├── Dockerfile
│   └── package.json
├── frontend-angular/
│   ├── src/
│   └── package.json
├── docker-compose.yml
├── .env.example
├── .env                   # No subir al repositorio
└── README.md
```

## Requisitos de Instalación

Para ejecutar el proyecto se requiere:

- Go 1.22 o superior.
- Node.js 20 o superior.
- Angular CLI compatible con Angular 22.
- Docker y Docker Compose (recomendado).
- Postman o curl para realizar las pruebas.

## Ejecución Local (sin Docker)

### 1. API Node.js

```bash
cd api-node

npm install

cp .env.example .env

npm start
```

La API Node.js estará disponible en:

`http://localhost:3000`

### 2. API Go

En una nueva terminal:

```bash
cd api-go

go mod tidy

go run cmd/main.go
```

La API Go estará disponible en:

`http://localhost:8080`

### 3. Frontend Angular

En una nueva terminal:

```bash
cd frontend-angular

npm install

npm start
```

Si el proyecto no tiene configurado el script `start`, se puede ejecutar mediante Angular CLI:

```bash
ng serve
```

El frontend estará disponible en la dirección local indicada por Angular CLI, normalmente:

`http://localhost:4200`

**Nota:** Asegurarse de que las variables de entorno (`JWT_SECRET`, `NODE_API_URL`, etc.) estén configuradas correctamente y que el frontend apunte a la URL correspondiente de la API Go.

## Ejecución con Docker

Copiar `.env.example` a `.env` y ajustar los valores según sea necesario.

Ejecutar:

```bash
docker-compose up --build
```

Esto levantará ambos servicios backend:

| Servicio | URL |
|---|---|
| API Go | http://localhost:8090 (mapeado al puerto interno 8080) |
| API Node.js | http://localhost:3000 |

Para ejecutar los servicios en segundo plano:

```bash
docker-compose up --build -d
```

Para detener los servicios:

```bash
docker-compose down
```

**Nota:** El frontend se despliega de forma independiente mediante Cloudflare Pages.

## Endpoints y Ejemplos

### 1. Login (API Go)

**Endpoint:**

`POST http://localhost:8090/api/v1/auth/login`

**Request:**

```json
{
  "username": "admin",
  "password": "password123"
}
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

El token obtenido debe utilizarse para acceder a los endpoints protegidos.

### 2. Factorización QR + Estadísticas (API Go, protegido)

**Endpoint:**

`POST http://localhost:8090/api/v1/qr`

**Headers:**

```text
Authorization: Bearer <token>
Content-Type: application/json
```

**Request:**

```json
{
  "matrix": [
    [1, 2],
    [3, 4],
    [5, 6]
  ]
}
```

**Response (200 OK):**

```json
{
  "q": [
    [-0.169..., 0.897...],
    [-0.507..., 0.276...],
    [-0.845..., -0.345...]
  ],
  "r": [
    [-5.916..., -7.437...],
    [0, 0.828...]
  ],
  "statistics": {
    "q": {
      "max": 0.897,
      "min": -0.845,
      "average": -0.115,
      "sum": -0.693,
      "isDiagonal": false
    },
    "r": {
      "max": 0.828,
      "min": -7.437,
      "average": -3.131,
      "sum": -12.525,
      "isDiagonal": false
    },
    "global": {
      "max": 0.897,
      "min": -7.437,
      "average": -1.321,
      "sum": -13.218,
      "isDiagonal": false
    }
  }
}
```

**Nota:** Los valores numéricos de la respuesta son ilustrativos. La respuesta real dependerá de la implementación y del algoritmo utilizado.

### 3. Health Checks (públicos)

Permiten verificar que los servicios están funcionando correctamente.

**API Go:**

`GET http://localhost:8090/health`

Respuesta:

```json
{
  "status": "ok",
  "service": "api-go"
}
```

**API Node.js:**

`GET http://localhost:3000/health`

Respuesta:

```json
{
  "status": "ok",
  "service": "api-node"
}
```

> **Nota:** Para probar los endpoints en producción, reemplaza `http://localhost:8090` por `https://talsory-api-go.onrender.com`.

## Pruebas

### Go

Ejecutar las pruebas unitarias:

```bash
cd api-go

go test ./tests/ -v
```

Pruebas incluidas:

- Reconstrucción de A (`Q * R ≈ A`).
- Ortogonalidad de Q (`Q^T * Q ≈ I`).
- Triangularidad de R.
- Dimensiones correctas.
- Casos borde (matriz vacía, columna de ceros).

### Node.js

Ejecutar las pruebas unitarias:

```bash
cd api-node

npm test
```

Pruebas incluidas:

- `flattenMatrix`.
- `calculateBasicStats`.
- `isDiagonal`.
- `calculateStats`.
- Casos típicos y borde (números negativos, tolerancia numérica).

## Despliegue en la Nube

La aplicación está completamente desplegada y funcional en servicios cloud gratuitos:

### 🔗 Enlaces en Vivo

- **Frontend (Angular + Tailwind):** [https://talsory-frontend.pages.dev](https://talsory-frontend.pages.dev)
- **API Go (Render):** [https://talsory-api-go.onrender.com/health](https://talsory-api-go.onrender.com/health)
- **API Node.js (Render):** [https://talsory-api-node.onrender.com/health](https://talsory-api-node.onrender.com/health)

### Infraestructura

- **Backend:** Desplegado en **Render** utilizando los Dockerfiles de cada servicio. Las variables de entorno (`JWT_SECRET`, `NODE_API_URL`, etc.) están configuradas directamente en el panel de Render.
- **Frontend:** Desplegado en **Cloudflare Pages**, configurado para compilar Angular 22 automáticamente desde el repositorio de GitHub.

## Limitaciones Conocidas

- **JWT simple:** No se implementaron refresh tokens ni revocación. Es suficiente para el alcance del reto.
- **Sin base de datos:** Las credenciales están configuradas mediante variables de entorno. En producción se utilizaría una base de datos.
- **Cold Starts:** Al usar el plan gratuito de Render, los servicios backend se "duermen" tras un periodo de inactividad. La primera petición puede tardar hasta 50 segundos en responder mientras los contenedores se reinician.
- **Matrices grandes:** El algoritmo Householder tiene una complejidad de O(m × n²). Para matrices muy grandes se recomendaría utilizar una librería optimizada como LAPACK.

## Autor

Hernán Ríos - Technical Challenge para Talsory Recruiting