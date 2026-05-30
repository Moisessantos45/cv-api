# CV-API

API REST para mi portafolio profesional. Construida con Go, PostgreSQL, Redis y Gin Framework.

## Stack Tecnológico

| Categoría | Tecnología | Versión |
|-----------|------------|---------|
| **Lenguaje** | Go | 1.25.6 |
| **Framework Web** | Gin | v1.11.0 |
| **ORM** | GORM | v1.31.1 |
| **Base de Datos** | PostgreSQL | - |
| **Caché/Sesiones** | Redis | v9.17.3 |
| **Tokens** | PASETO | v2 |
| **Email** | Resend + Google SMTP | v3.6.0 |

### Librerías Adicionales

| Librería | Propósito |
|----------|-----------|
| `paseto` | Tokens PASETO (alternativa moderna a JWT) |
| `golang.org/x/crypto` | Hash bcrypt para contraseñas |
| `gotp` | Generador TOTP para 2FA |
| `gin-contrib/cors` | Configuración CORS |
| `gin-contrib/gzip` | Compresión Gzip |
| `gin-helmet` | Headers de seguridad HTTP |
| `golang.org/x/time` | Rate limiting |
| `resend-go` | Envío de emails transaccionales |
| `godotenv` | Carga de variables de entorno |

---

## Arquitectura

El proyecto sigue una **Arquitectura Limpia (Clean Architecture)** con organización **Feature-Based**, separando el código por dominio de negocio.

### Flujo de una Request

```
HTTP Request
    ↓
Middleware (CORS, Gzip, Helmet)
    ↓
Route Handler (routes/)
    ↓
Feature Handler (features/*/handler.go)
    ↓
Feature UseCase (features/*/usecase.go)
    ↓
Repository (features/*/repository.go)
    ↓
GORM → PostgreSQL / Redis
```

### Estructura de Carpetas

```
cv-api/
├── main.go                     # Punto de entrada
│
├── config/                     # Configuración global
│   ├── db/
│   │   ├── connection.go      # Conexión PostgreSQL (GORM)
│   │   ├── initialize.go      # Migraciones automáticas
│   │   └── seed.go            # Datos iniciales (estados de posts)
│   ├── redis.go               # Cliente Redis
│   └── email.go               # Configuración de email (Brevo/SMTP)
│
├── internal/
│   ├── models/                # Modelos de datos compartidos
│   │   ├── models.go          # Auth, Profile, Project, Video, Post, Experience
│   │   ├── paginate.go        # Estructura de paginación
│   │   └── serializer.go       # Tipos JSON personalizados
│   │
│   ├── routes/                # Definición de rutas e inyección de dependencias
│   │   ├── deps.go            # Inicialización de repos y services
│   │   ├── auth.go            # Rutas de autenticación
│   │   ├── profile.go         # Rutas de perfil
│   │   ├── project.go         # Rutas de proyectos
│   │   ├── video.go           # Rutas de videos
│   │   └── post.go            # Rutas de posts
│   │
│   ├── features/              # Módulos de negocio por dominio
│   │   ├── auth/             # Autenticación, tokens, 2FA
│   │   │   ├── domain.go     # Interfaces Repository/Service, factory NewUser
│   │   │   ├── handler.go    # Handlers HTTP
│   │   │   ├── repository.go # Implementación PostgreSQL
│   │   │   └── usecase.go    # Lógica de negocio
│   │   │
│   │   ├── profile/          # Perfil profesional
│   │   ├── post/             # Blog y artículos
│   │   ├── project/          # Portafolio de proyectos
│   │   └── stream/           # Videos y streaming
│   │
│   └── shared/               # Utilidades compartidas
│       ├── utils/            # Helpers (JWT, bcrypt, email, normalización)
│       ├── service/          # CacheService (caché genérico)
│       ├── middleware/       # AuthMiddleware, RateLimiterMiddleware
│       ├── templates/        # Plantillas HTML para emails
│       └── errorsx/          # Errores predefinidos
```

---

## Patrones de Diseño

### Repository Pattern
Cada feature define interfaces que son implementadas por repositorios PostgreSQL:

```go
type AuthRepository interface {
    Register(user *models.Auth) error
    GetByEmail(email string) (*models.Auth, error)
    Update(ctx context.Context, auth *models.Auth) error
}
```

### Service/UseCase Pattern
La lógica de negocio se encapsula en UseCases, separados de los Handlers HTTP.

### Factory Pattern
Cada dominio tiene funciones factory (`NewUser`, `NewProfile`, `NewPost`, etc.) para crear entidades con validación.

### Cache-Aside Pattern
El `CacheService` implementa el patrón cache-aside con Redis:
1. Buscar en caché primero
2. Si miss, consultar DB
3. Guardar en caché con TTL
4. Retornar resultado

---

## Modelos de Datos

| Entidad | Tabla | Descripción |
|---------|-------|-------------|
| **Auth** | `auths` | Usuarios (email, password_hash, 2FA, email_confirmed) |
| **Profile** | `profiles` | Perfil profesional (nombre, CV, links sociales, descripción) |
| **Project** | `projects` | Proyectos (slug, título, tecnologías, banner, imágenes) |
| **Video** | `videos` | Videos (título, URL, descripción, estado) |
| **Post** | `posts` | Artículos (slug, título, contenido, banner, tags, categoría) |
| **Experience** | `experiences` | Experiencia laboral |
| **StatePost** | `state_posts` | Estados (draft, published, archived, private, unlisted) |

---

## Seguridad

### Tokens PASETO
Se usa **PASETO** (Protocol for Attributable SEAs with TOKen) en lugar de JWT:
- Más seguro: no vulnerable a ataques de algoritmo
- Access token: 15 minutos de TTL
- Refresh token: 8 horas de TTL

### Autenticación de Dos Factores (2FA)
- Implementación con **TOTP** (Time-based One-Time Password)
- Secret generado por usuario
- Códigos temporales de 6 dígitos

### Medidas de Seguridad
- **bcrypt**: Hash de contraseñas con salt automático
- **Helmet**: Headers de seguridad HTTP
- **CORS**: Whitelist de orígenes permitidos
- **Gzip**: Compresión de respuestas
- **Rate Limiting**: Limitación de requests por IP (configurable)

---

## Endpoints Principales

### Rutas Públicas

| Método | Path | Descripción |
|--------|------|-------------|
| GET | `/` | Mensaje de bienvenida |
| GET | `/health` | Health check |
| GET | `/api/v1/project/public` | Listar proyectos públicos |
| GET | `/api/v1/project/:slug` | Detalle proyecto |
| GET | `/api/v1/project/recent` | Proyectos recientes |
| GET | `/api/v1/post/public` | Listar posts públicos |
| GET | `/api/v1/post/:slug` | Detalle post |
| GET | `/api/v1/post/recent` | Posts recientes |
| GET | `/api/v1/stream` | Listar videos (paginado) |
| GET | `/api/v1/stream/:id` | Detalle video |

### Autenticación (`/api/v1/auth`)

| Método | Path | Descripción |
|--------|------|-------------|
| POST | `/login` | Login con email/password |
| POST | `/register` | Registro de usuario |
| POST | `/logout` | Cerrar sesión |
| POST | `/refresh-token` | Refrescar access token |
| POST | `/forgot-password` | Solicitar reset de contraseña |
| POST | `/forward-email-verification` | Reenviar email de confirmación |
| PATCH | `/reset-password` | Resetear contraseña (con token) |
| PATCH | `/change-password` | Cambiar contraseña (autenticado) |
| GET | `/session` | Obtener sesión actual |
| GET | `/confirm-account` | Confirmar cuenta por email |
| POST | `/verify-email` | Verificar email |
| GET | `/generate-two-factor` | Generar secret 2FA |

### Rutas Protegidas (requieren JWT)

| Método | Path | Descripción |
|--------|------|-------------|
| POST | `/api/v1/profile` | Crear perfil |
| GET | `/api/v1/profile` | Obtener mi perfil |
| PUT | `/api/v1/profile` | Actualizar perfil |
| POST | `/api/v1/post` | Crear post |
| GET | `/api/v1/post/user` | Listar mis posts |
| PUT | `/api/v1/post/:slug` | Actualizar post |
| POST | `/api/v1/project` | Crear proyecto |
| PUT | `/api/v1/project/:id` | Actualizar proyecto |
| PATCH | `/api/v1/project/:id/state` | Cambiar estado |
| POST | `/api/v1/stream` | Crear video |
| PUT | `/api/v1/stream` | Actualizar video |

---

## Configuración

### Variables de Entorno

```env
# Environment
GO_ENV=dev

# URLs
HOST_URL_DEV=http://localhost:5173
HOST_URL_PROD=http://localhost:3000
HOST_API_DEV=http://localhost:4100

# Database PostgreSQL
HOST=localhost
PORT=5432
DBUSER=user
PASSWORD=tu_password
DBNAME=name

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Email (SMTP + Resend)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=tu_email@gmail.com
API_KEY_SMTP=tu_api_key

# JWT/Auth
SECRET_KEY_JWT=tu_secret_key
```

### Servidor

```go
Addr: ":4100"
```

---

## Flujo de Autenticación

```
1. POST /api/v1/auth/login
       ↓
   Validar email/password contra tabla auths
       ↓
   ¿2FA habilitado?
   ├── Sí → Guardar pre-auth en Redis → Retornar ErrPending2FA
   └── No → Continuar
       ↓
   Generar access token (15 min) + refresh token (8 horas)
       ↓
   Crear sesión en Redis: session:{sessionID}
       ↓
   Guardar refresh token en cookie httponly
       ↓
   Retornar tokens al cliente
```

---

## Caché con Redis

### Recursos Cacheados
- Perfiles: `profiles:all:{page}:{pageSize}`
- Posts: `posts:all:public:{page}:{pageSize}`
- Proyectos: `projects:all:public:{page}:{pageSize}`
- Videos: `videos:all:{page}:{pageSize}`

### Claves de Sesión
- `session:{sessionID}` - Datos de sesión (TTL: 8 horas)
- `preauth:{sessionID}` - Pre-auth para flujo 2FA (TTL: 5 minutos)
- `confirm:{userID}` - Token de confirmación email (TTL: 15 minutos)

---

## Migraciones y Seeds

Las migraciones se ejecutan automáticamente en modo desarrollo (`GO_ENV=dev`) mediante `DB.AutoMigrate()`.

En producción (`GO_ENV=production`) las migraciones se saltan.

Seed de datos iniciales: estados de posts (draft, published, archived, private, unlisted).

---

## Diagramas

### Arquitectura General

```
                    ┌──────────────────────────────────────────────────────┐
                    │                      Gin Engine                        │
                    │  ┌─────────┐ ┌─────────┐ ┌───────┐ ┌──────────────┐ │
                    │  │  CORS   │ │  Gzip   │ │Helmet │ │ Rate Limiter │ │
                    │  └────┬────┘ └────┬────┘ └───┬───┘ └──────┬───────┘ │
                    └───────┼──────────┼──────────┼────────────┼─────────┘
                            │          │          │            │
                    ┌───────▼──────────▼──────────▼────────────▼──────────┐
                    │                      Routes                           │
                    │  /api/v1/auth/*  │  /api/v1/profile/*               │
                    │  /api/v1/project/* │ /api/v1/post/*                 │
                    │  /api/v1/stream/*                                       │
                    └──────────────────────┬───────────────────────────────┘
                                           │
                    ┌──────────────────────▼───────────────────────────────┐
                    │               Feature Handlers                        │
                    │  auth.Handler │ profile.Handler │ post.Handler         │
                    │  project.Handler │ stream.Handler                    │
                    └──────────────────────┬───────────────────────────────┘
                                           │
                    ┌──────────────────────▼───────────────────────────────┐
                    │                  Use Cases                            │
                    │  auth.UseCase  │ profile.UseCase │ post.UseCase      │
                    │  project.UseCase │ stream.UseCase                    │
                    └──────────────────────┬───────────────────────────────┘
                                           │
                    ┌──────────────────────▼───────────────────────────────┐
                    │                 Repositories                           │
                    │  auth.PostgresRepository │ profile.PostgresRepo      │
                    │  post.PostgresRepository │ project.PostgresRepo      │
                    │  stream.PostgresRepository                            │
                    └──────────────────────┬───────────────────────────────┘
                                           │
                    ┌──────────────────────▼───────────────────────────────┐
                    │    Cache Service (Redis) │ PostgreSQL (GORM)          │
                    └─────────────────────────────────────────────────────┘
```
