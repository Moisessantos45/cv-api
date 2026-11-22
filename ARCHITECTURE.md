# Guía de Estructura de Funcionalidad - cv-api

## 1. Arquitectura General

Cada funcionalidad vive en `internal/features/<nombre>/` y sigue el patrón **Clean Architecture**:

```
<nombre>/
├── domain.go      # Interfaces + Structs de dominio + Lógica de validación
├── repository.go # Implementación de acceso a datos (PostgreSQL)
├── usecase.go     # Lógica de negocio (inyección de dependencias)
└── handler.go    # Controladores HTTP (Gin)
```

---

## 2. Nomenclatura y Convenciones

### 2.1 Estructuras

| Tipo | Nomenclatura | Ejemplo |
|------|-------------|---------|
| Modelo de datos | `models.<Nombre>` | `models.Project` |
| Struct básica (lista pública) | `<Nombre>Basic` | `ProjectBasic` |
| Struct info mínima | `<Nombre>InfoBasic` | `ProjectInfoBasic` |
| Struct paginada | `<Nombre>sPaginated` | `PaginatedProjects` |
| Struct uso interno | `<Nombre>UseCase` | `ProjectUseCase` |
| Struct repository | `PostgresRepository` | `PostgresRepository` |
| Struct handler | `<Nombre>Handler` | `ProjectHandler` |

### 2.2 Interfaces

| Tipo | Nomenclatura | Ejemplo |
|------|-------------|---------|
| Repositorio | `<Nombre>Repository` | `ProjectRepository` |
| Servicio/UseCase | `<Nombre>Service` | `ProjectService` |

### 2.3 Métodos

| Visibilidad | Prefijo | Ejemplo |
|-------------|---------|---------|
| Público (sin auth) | `GetAllPublic`, `GetBySlugPublic` | `GetAllPublic`, `GetBySlugPublic` |
| Privado (requiere auth) | `GetAll`, `GetBySlug` | `GetAll`, `GetBySlug` |

---

## 3. domain.go

### 3.1 Estructuras obligatorias

```go
// Struct para listado básico (paginación pública)
type <Nombre>Basic struct {
    ID        uint64 `json:"id" gorm:"column:id"`
    Slug      string `json:"slug" gorm:"column:slug"`
    Title     string `json:"title" gorm:"column:title"`
    // ... campos públicos mínimos para listados
    CreatedAt string `json:"created_at" gorm:"column:created_at"`
}

// Struct para info mínima (actualización, comparaciones)
type <Nombre>InfoBasic struct {
    ID    uint64 `json:"id" gorm:"column:id"`
    Slug  string `json:"slug" gorm:"column:slug"`
    Title string `json:"title" gorm:"column:title"`
}
```

### 3.2 Interfaces

```go
type <Nombre>Repository interface {
    WithTransaction(fn func(repo *PostgresRepository) error) error

    // Lecturas
    GetAllData(ctx context.Context) ([]models.<Nombre>, error)
    GetAll(ctx context.Context, offset int, limit int) ([]<Nombre>Basic, int64, error)
    GetAllPublic(ctx context.Context, offset int, limit int) ([]models.<Nombre>, int64, error)
    GetAllRecents(ctx context.Context, max int) ([]models.<Nombre>, error)
    GetByID(ctx context.Context, id uint64) (*models.<Nombre>, error)
    GetBySlug(ctx context.Context, slug string) (*models.<Nombre>, error)
    GetBySlugPublic(ctx context.Context, slug string) (*models.<Nombre>, error)
    GetBasicInfoBySlug(ctx context.Context, slug string) (*<Nombre>InfoBasic, error)

    // Existencias
    ExistsById(ctx context.Context, id uint64) (bool, error)

    // Escritura
    Create(ctx context.Context, data *models.<Nombre>) error
    Update(ctx context.Context, id uint64, data map[string]any) error
    Delete(ctx context.Context, id uint64) error
}

type <Nombre>Service interface {
    // Lecturas
    GetAllData(ctx context.Context) ([]models.<Nombre>, error)
    GetAll(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[<Nombre>Basic], error)
    GetAllPublic(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[models.<Nombre>], error)
    GetAllRecents(ctx context.Context, max int) (*service.PaginatedResult[models.<Nombre>], error)
    GetBySlug(ctx context.Context, slug string) (*models.<Nombre>, error)
    GetBySlugPublic(ctx context.Context, slug string) (*models.<Nombre>, error)
    GetByID(ctx context.Context, id uint64) (*models.<Nombre>, error)

    // Escritura
    Create(ctx context.Context, data *models.<Nombre>) error
    CreateAll(ctx context.Context, data []models.<Nombre>) error
    Update(ctx context.Context, slug string, data *models.<Nombre>) error
    Delete(ctx context.Context, id uint64) error

    // Especiales
    UpdateState(ctx context.Context, id uint64, stateID uint64) error
}
```

### 3.3 Implementación del UseCase

```go
type <Nombre>UseCase struct {
    repo  <Nombre>Repository
    cache *service.CacheService
}

func New<Nombre>Service(repo <Nombre>Repository, cache *service.CacheService) <Nombre>Service {
    return &<Nombre>UseCase{
        repo:  repo,
        cache: cache,
    }
}
```

### 3.4 Funciones de dominio obligatorias

```go
// Constructor con validación para creación
func New<Nombre>(data *models.<Nombre>) (*models.<Nombre>, error) {
    if strings.TrimSpace(data.Title) == "" {
        return nil, errors.New("el título es obligatorio")
    }
    // ... más validaciones

    data.Slug = generateSlug(data.Title)
    return &models.<Nombre>{...}, nil
}

// Builder para actualización (solo campos presentes)
func Build<Nombre>UpdateData(data *models.<Nombre>, changeTitle bool) map[string]any {
    updateData := make(map[string]any)

    if strings.TrimSpace(data.Title) != "" && !changeTitle {
        updateData["title"] = utils.NormalizeText(data.Title)
        updateData["slug"] = generateSlug(data.Title)
    }

    if strings.TrimSpace(data.Description) != "" {
        updateData["description"] = utils.NormalizeText(data.Description)
    }

    // ... otros campos opcionales
    return updateData
}

// Helper para paginación
func NewPaginationResult[T any](data []T, total int64, totalPages int, page int, pageSize int) *service.PaginatedResult[T] {
    return &service.PaginatedResult[T]{
        Data: data,
        Paginate: models.Pagination{
            Total:      total,
            TotalPages: totalPages,
            Page:       page,
            PageSize:   pageSize,
        },
    }
}

// Generador de slug
func generateSlug(title string) string {
    slugNormalized := utils.NormalizeText(title)
    slugBase := strings.ToLower(strings.TrimSpace(slugNormalized))
    slugBase = strings.ReplaceAll(slugBase, " ", "-")

    suffix := fmt.Sprintf("-%d", time.Now().UnixNano())
    // truncamiento a máximo 55 caracteres
    return slugBase + suffix
}
```

---

## 4. repository.go

### 4.1 Estructura y constructor

```go
type PostgresRepository struct {
    db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) <Nombre>Repository {
    return &PostgresRepository{db: db}
}
```

### 4.2 Métodos de lectura

```go
// Listado completo (sin paginación, para exportación)
func (r *PostgresRepository) GetAllData(ctx context.Context) ([]models.<Nombre>, error) {
    var items []models.<Nombre>
    err := r.db.WithContext(ctx).Find(&items).Error
    return items, err
}

// Listado paginado básico
func (r *PostgresRepository) GetAll(ctx context.Context, offset int, limit int) ([]<Nombre>Basic, int64, error) {
    var items []<Nombre>Basic
    var total int64

    db := r.db.WithContext(ctx).
        Select("id", "slug", "title", "created_at").
        Order("created_at DESC")

    err := db.Model(&models.<Nombre>{}).Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    if total == 0 {
        return []<Nombre>Basic{}, 0, nil
    }

    err = db.Offset(offset).Limit(limit).Find(&items).Error
    return items, total, err
}

// Listado público (puede incluir más campos que Basic)
func (r *PostgresRepository) GetAllPublic(ctx context.Context, offset int, limit int) ([]models.<Nombre>, int64, error) {
    var items []models.<Nombre>
    var total int64

    db := r.db.WithContext(ctx).
        Select("id", "slug", "title", "description", "banner", "created_at").
        Order("created_at DESC")

    err := db.Model(&models.<Nombre>{}).Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    err = db.Offset(offset).Limit(limit).Find(&items).Error
    return items, total, err
}

// Por ID
func (r *PostgresRepository) GetByID(ctx context.Context, id uint64) (*models.<Nombre>, error) {
    var item models.<Nombre>
    if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("<nombre> no encontrado para el id: %d", id)
        }
        return nil, err
    }
    return &item, nil
}

// Por slug
func (r *PostgresRepository) GetBySlug(ctx context.Context, slug string) (*models.<Nombre>, error) {
    var item models.<Nombre>
    if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&item).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("<nombre> no encontrado para el slug: %s", slug)
        }
        return nil, err
    }
    return &item, nil
}

// Por slug público (con filtro de estado)
func (r *PostgresRepository) GetBySlugPublic(ctx context.Context, slug string) (*models.<Nombre>, error) {
    var item models.<Nombre>
    err := r.db.WithContext(ctx).
        Select("id", "slug", "title", "description", "...").
        Where("slug = ? AND state_id = ?", slug, 2).
        First(&item).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}
```

### 4.3 Métodos de escritura

```go
func (r *PostgresRepository) Create(ctx context.Context, data *models.<Nombre>) error {
    return r.db.WithContext(ctx).Create(data).Error
}

func (r *PostgresRepository) Update(ctx context.Context, id uint64, data map[string]any) error {
    result := r.db.WithContext(ctx).Model(&models.<Nombre>{}).Where("id = ?", id).Updates(data)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("<nombre> no encontrado para el id: %d", id)
    }
    return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint64) error {
    result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.<Nombre>{})
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("<nombre> no encontrado para el id: %d", id)
    }
    return nil
}

func (r *PostgresRepository) ExistsById(ctx context.Context, id uint64) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&models.<Nombre>{}).Where("id = ?", id).Count(&count).Error
    return count > 0, err
}
```

### 4.4 Transacciones

```go
func (r *PostgresRepository) WithTransaction(fn func(repo *PostgresRepository) error) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        txRepo := &PostgresRepository{db: tx}
        return fn(txRepo)
    })
}
```

---

## 5. usecase.go

### 5.1 Estructura con dependencias

```go
const <nombre>CacheKey = "<nombre>s:all"

type <Nombre>UseCase struct {
    repo  <Nombre>Repository
    cache *service.CacheService
}
```

### 5.2 Métodos de lectura con paginación

```go
// Listado con paginación y caché
func (s *<Nombre>UseCase) GetAll(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[<Nombre>Basic], error) {
    return service.Paginate[<Nombre>Basic](
        ctx,
        s.cache,
        <nombre>CacheKey,
        page,
        pageSize,
        s.repo.GetAll,
        NewPaginationResult[<Nombre>Basic],
    )
}

// Listado público con paginación y caché
func (s *<Nombre>UseCase) GetAllPublic(ctx context.Context, page int, pageSize int) (*service.PaginatedResult[models.<Nombre>], error) {
    return service.Paginate[models.<Nombre>](
        ctx,
        s.cache,
        <nombre>CacheKey+":public",
        page,
        pageSize,
        s.repo.GetAllPublic,
        NewPaginationResult[models.<Nombre>],
    )
}

// Recientes con límite
func (s *<Nombre>UseCase) GetAllRecents(ctx context.Context, max int) (*service.PaginatedResult[models.<Nombre>], error) {
    return service.GetRecents[models.<Nombre>](
        ctx,
        s.cache,
        <nombre>CacheKey,
        max,
        s.repo.GetAllRecents,
        NewPaginationResult[models.<Nombre>],
    )
}
```

### 5.3 Métodos de lectura por identificador

```go
func (s *<Nombre>UseCase) GetByID(ctx context.Context, id uint64) (*models.<Nombre>, error) {
    if id == 0 {
        return nil, fmt.Errorf("ID inválido")
    }
    return s.repo.GetByID(ctx, id)
}

func (s *<Nombre>UseCase) GetBySlug(ctx context.Context, slug string) (*models.<Nombre>, error) {
    if slug == "" {
        return nil, fmt.Errorf("el slug es requerido")
    }
    return s.repo.GetBySlug(ctx, slug)
}

func (s *<Nombre>UseCase) GetBySlugPublic(ctx context.Context, slug string) (*models.<Nombre>, error) {
    if slug == "" {
        return nil, fmt.Errorf("el slug es requerido")
    }
    return s.repo.GetBySlugPublic(ctx, slug)
}
```

### 5.4 Métodos de escritura

```go
// Crear uno
func (s *<Nombre>UseCase) Create(ctx context.Context, data *models.<Nombre>) error {
    if data == nil {
        return fmt.Errorf("los datos son requeridos")
    }

    validatedData, err := New<Nombre>(data)
    if err != nil {
        return err
    }

    err = s.repo.Create(ctx, validatedData)
    if err != nil {
        return fmt.Errorf("error al crear el <nombre>: %w", err)
    }

    *data = *validatedData
    return nil
}

// Crear múltiples (batch)
func (s *<Nombre>UseCase) CreateAll(ctx context.Context, data []models.<Nombre>) error {
    return s.repo.WithTransaction(func(repo *PostgresRepository) error {
        for i := range data {
            if err := s.Create(ctx, &data[i]); err != nil {
                return err
            }
        }
        return nil
    })
}

// Actualizar
func (s *<Nombre>UseCase) Update(ctx context.Context, slug string, data *models.<Nombre>) error {
    if slug == "" {
        return fmt.Errorf("el slug es requerido")
    }
    if data == nil {
        return fmt.Errorf("los datos son requeridos")
    }

    existing, err := s.repo.GetBasicInfoBySlug(ctx, slug)
    if err != nil {
        return fmt.Errorf("error al buscar el <nombre>: %v", err)
    }

    isEqualTitle := compareStrings(existing.Title, data.Title)
    updateData := Build<Nombre>UpdateData(data, isEqualTitle)

    err = s.repo.Update(ctx, existing.ID, updateData)
    if err != nil {
        return fmt.Errorf("error al actualizar el <nombre>: %w", err)
    }

    return nil
}
```

### 5.5 Métodos especiales

```go
func (s *<Nombre>UseCase) UpdateState(ctx context.Context, id uint64, stateID uint64) error {
    if id == 0 {
        return fmt.Errorf("el ID es requerido")
    }

    exists, err := s.repo.ExistsById(ctx, id)
    if err != nil {
        return fmt.Errorf("error al verificar existencia: %v", err)
    }
    if !exists {
        return fmt.Errorf("el <nombre> no existe")
    }

    err = s.repo.UpdateState(ctx, id, stateID)
    if err != nil {
        return fmt.Errorf("error al actualizar el estado: %w", err)
    }
    return nil
}
```

---

## 6. handler.go

### 6.1 Estructura y constructor

```go
type <Nombre>Handler struct {
    service <Nombre>Service
}

func New<Nombre>Handlers(service <Nombre>Service) *<Nombre>Handler {
    return &<Nombre>Handler{service: service}
}
```

### 6.2 Handlers de lectura

```go
// Listado con paginación
func (h *<Nombre>Handler) GetAll(c *gin.Context) {
    page, pageSize, _, err := utils.ValidateQueryPagination(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    result, err := h.service.GetAll(c.Request.Context(), page, pageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": result, "message": "Data Success"})
}

// Listado público
func (h *<Nombre>Handler) GetAllPublic(c *gin.Context) {
    page, pageSize, _, err := utils.ValidateQueryPagination(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    result, err := h.service.GetAllPublic(c.Request.Context(), page, pageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": result, "message": "Data Success"})
}

// Recientes
func (h *<Nombre>Handler) GetAllRecents(c *gin.Context) {
    maxSize, err := utils.ValidateParamsQuery[int](c, "max")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    result, err := h.service.GetAllRecents(c.Request.Context(), maxSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": result, "message": "Data Success"})
}

// Por ID
func (h *<Nombre>Handler) GetByID(c *gin.Context) {
    id, err := utils.ValidateParamsId(c, "")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    result, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": result})
}

// Por slug
func (h *<Nombre>Handler) GetBySlug(c *gin.Context) {
    slug := c.Param("slug")
    if slug == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Slug is required"})
        return
    }

    result, err := h.service.GetBySlug(c.Request.Context(), slug)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": result, "message": "Data Success"})
}

// Por slug público
func (h *<Nombre>Handler) GetBySlugPublic(c *gin.Context) {
    slug := c.Param("slug")
    if slug == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Slug is required"})
        return
    }

    result, err := h.service.GetBySlugPublic(c.Request.Context(), slug)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": result, "message": "Data Success"})
}
```

### 6.3 Handlers de escritura

```go
// Crear
func (h *<Nombre>Handler) Create(c *gin.Context) {
    var model models.<Nombre>
    if err := c.ShouldBindJSON(&model); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Datos inválidos", "error": err.Error()})
        return
    }

    err := h.service.Create(c.Request.Context(), &model)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"data": model, "message": "Creado correctamente"})
}

// Actualizar
func (h *<Nombre>Handler) Update(c *gin.Context) {
    slug := c.Param("slug")
    if slug == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Slug is required"})
        return
    }

    var model models.<Nombre>
    if err := c.ShouldBindJSON(&model); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Datos inválidos", "error": err.Error()})
        return
    }

    err := h.service.Update(c.Request.Context(), slug, &model)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Actualizado correctamente"})
}

// Eliminar
func (h *<Nombre>Handler) Delete(c *gin.Context) {
    id, err := utils.ValidateParamsId(c, "")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    err = h.service.Delete(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Eliminado correctamente"})
}
```

---

## 7. Cuándo usar structs Basic vs models completas

| Contexto | Struct a usar |
|----------|---------------|
| Listado paginado (público) | `<Nombre>Basic` |
| Listado paginado (admin) | `<Nombre>Basic` o `models.<Nombre>` según campos necesarios |
| Detalle por ID/Slug | `models.<Nombre>` |
| Crear/Actualizar | `models.<Nombre>` (validada por `New<Nombre>`) |
| Transmisión de datos (API) | Según necesidad - no siempre todo el modelo |

---

## 8. Respuestas HTTP según tipo de operación

| Operación | Código | Estructura |
|-----------|--------|------------|
| Listar/Ver | `200 OK` | `{"data": ..., "message": "Data Success"}` |
| Crear | `201 Created` | `{"data": ..., "message": "Creado correctamente"}` |
| Actualizar | `200 OK` | `{"message": "Actualizado correctamente"}` |
| Eliminar | `200 OK` | `{"message": "Eliminado correctamente"}` |
| Error de validación | `400 Bad Request` | `{"message": "...", "error": "..."}` |
| No encontrado | `404 Not Found` | Manejado en el handler (error del repo) |
| Error interno | `500 Internal Server Error` | `{"message": "..."}` |

---

## 9. Validaciones comunes

### 9.1 ID inválido
```go
if id == 0 {
    return nil, fmt.Errorf("ID inválido")
}
```

### 9.2 Slug vacío
```go
if slug == "" {
    return nil, fmt.Errorf("el slug es requerido")
}
```

### 9.3 Datos nulos
```go
if data == nil {
    return fmt.Errorf("los datos son requeridos")
}
```

### 9.4 Validaciones de dominio (en `New<Nombre>`)
```go
if strings.TrimSpace(data.Title) == "" {
    return nil, errors.New("el título es obligatorio")
}

if len(data.Items) == 0 {
    return nil, errors.New("los items son obligatorios")
}
```

---

## 10. Inyección de dependencias

### 10.1 En repository.go
```go
func NewPostgresRepository(db *gorm.DB) <Nombre>Repository {
    return &PostgresRepository{db: db}
}
```

### 10.2 En usecase.go
```go
func New<Nombre>Service(repo <Nombre>Repository, cache *service.CacheService) <Nombre>Service {
    return &<Nombre>UseCase{repo: repo, cache: cache}
}
```

### 10.3 En handler.go
```go
func New<Nombre>Handlers(service <Nombre>Service) *<Nombre>Handler {
    return &<Nombre>Handler{service: service}
}
```

---

## 11. Métodos públicos vs privados

Se distingue por el nombre del método:

- **Público** (sin auth): `GetAllPublic`, `GetBySlugPublic`, `GetAllRecents`
  - Puede ser accedido sin autenticación
  - Filtra solo datos activos/publicados
  - En repository suelen filtrar por `state_id = 2`

- **Privado** (requiere auth): `GetAll`, `GetBySlug`
  - Puede requerir autenticación
  - Incluye datos en cualquier estado
  - Más campos en el SELECT

---

## 12. Métodos que todo módulo debe tener

### Repository
- `GetAllData` - Exportación completa
- `GetAll` - Listado paginado (básico)
- `GetAllPublic` - Listado público
- `GetAllRecents` - Listado recientes
- `GetByID` - Por ID
- `GetBySlug` - Por slug
- `GetBySlugPublic` - Por slug público
- `GetBasicInfoBySlug` - Info mínima para update
- `ExistsById` - Verificar existencia
- `Create` - Crear
- `Update` - Actualizar
- `UpdateState` - Cambiar estado
- `WithTransaction` - Transacción

### UseCase
- `GetAll` - Con paginación y caché
- `GetAllPublic` - Público con paginación
- `GetAllRecents` - Recientes
- `GetByID` - Por ID
- `GetBySlug` - Por slug
- `GetBySlugPublic` - Por slug público
- `Create` - Crear
- `CreateAll` - Crear batch
- `Update` - Actualizar
- `UpdateState` - Cambiar estado

### Handler
- `GetAll` - Listar
- `GetAllPublic` - Listar público
- `GetAllRecents` - Recientes
- `GetByID` - Por ID
- `GetBySlug` - Por slug
- `GetBySlugPublic` - Por slug público
- `Create` - Crear
- `CreateAll` - Crear desde archivo
- `Update` - Actualizar
- `UpdateState` - Cambiar estado