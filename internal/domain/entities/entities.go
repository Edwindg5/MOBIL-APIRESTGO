package entities

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Usuario representa un usuario del sistema
type Usuario struct {
	ID             int       `db:"id"`
	Email          string    `db:"email"`
	Password       string    `db:"password"`
	NombreCompleto string    `db:"nombre_completo"`
	Telefono       *string   `db:"telefono"`
	Rol            string    `db:"rol"`
	Estado         string    `db:"estado"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// RegisterRequest es la solicitud de registro de nuevo usuario
type RegisterRequest struct {
	Nombre   string `json:"nombre" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Telefono string `json:"telefono" validate:"required"`
}

// RegisterResponse es la respuesta tras registrar un usuario
type RegisterResponse struct {
	IDUsuario     int       `json:"id_usuario"`
	Nombre        string    `json:"nombre"`
	Email         string    `json:"email"`
	Rol           string    `json:"rol"`
	FechaRegistro time.Time `json:"fecha_registro"`
}

// PerfilResponse es la respuesta con el perfil del usuario autenticado
type PerfilResponse struct {
	IDUsuario     int       `json:"id_usuario"`
	Nombre        string    `json:"nombre"`
	Email         string    `json:"email"`
	Rol           string    `json:"rol"`
	Telefono      *string   `json:"telefono"`
	Estado        string    `json:"estado"`
	FechaRegistro time.Time `json:"fecha_registro"`
}

// UpdatePerfilRequest es la solicitud de actualización de perfil
type UpdatePerfilRequest struct {
	Nombre   string `json:"nombre" validate:"required,min=2"`
	Telefono string `json:"telefono" validate:"required"`
}

// ChangePasswordRequest es la solicitud de cambio de contraseña
type ChangePasswordRequest struct {
	PasswordActual string `json:"password_actual" validate:"required"`
	PasswordNueva  string `json:"password_nueva" validate:"required,min=6"`
}

// LoginRequest es la solicitud de login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse es la respuesta de login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // segundos
	Usuario      UsuarioPublicInfo `json:"usuario"`
}

// UsuarioPublicInfo es la información pública del usuario
type UsuarioPublicInfo struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	NombreCompleto string `json:"nombre_completo"`
	Rol            string `json:"rol"`
}

// RefreshTokenRequest es la solicitud para renovar token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse es la respuesta de renovación
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// JWTClaims son los claims del JWT
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Rol    string `json:"rol"`
}

// ProvisioningToken representa un token de provisioning de un solo uso
type ProvisioningToken struct {
	ID        int       `db:"id"`
	ESP32ID   string    `db:"esp32_id"`
	Token     string    `db:"token"` // hash del token original
	UsuarioID int       `db:"usuario_id"`
	UsedAt    *time.Time `db:"used_at"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

// LinkDeviceRequest es la solicitud para vincular un dispositivo vía QR
type LinkDeviceRequest struct {
	ESP32ID           string `json:"esp32_id" validate:"required"`
	ProvisioningToken string `json:"provisioning_token" validate:"required"`
}

// LinkDeviceResponse es la respuesta de vinculación (retorna el lote creado)
type LinkDeviceResponse struct {
	Lote    *LoteCafe `json:"lote"`
	Message string    `json:"message"`
}

// Sensor representa un ESP32
type Sensor struct {
	ID                int        `db:"id"`
	ESP32ID           string     `db:"esp32_id"`
	MacAddress        *string    `db:"mac_address"`
	LoteID            *int       `db:"lote_id"`
	ProvisioningToken *string    `db:"provisioning_token"`
	TokenUsado        bool       `db:"token_usado"`
	LinkedAt          *time.Time `db:"linked_at"`
	LastSeen          *time.Time `db:"last_seen"`
	Estado            string     `db:"estado"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}

// LoteCafe representa un lote de secado de café
type LoteCafe struct {
	ID                int        `db:"id" json:"id"`
	UsuarioID         int        `db:"usuario_id" json:"usuario_id"`
	NombreLote        string     `db:"nombre_lote" json:"nombre_lote"`
	Variedad          string     `db:"variedad" json:"variedad"`
	TipoProceso       string     `db:"tipo_proceso" json:"tipo_proceso"`
	PesoKg            float64    `db:"peso_kg" json:"peso_kg"`
	Ubicacion         string     `db:"ubicacion" json:"ubicacion"`
	IDSensor          *int       `db:"id_sensor" json:"id_sensor"`
	CodigoQR          string     `db:"codigo_qr" json:"codigo_qr"`
	Estado            string     `db:"estado" json:"estado"` // 'en_proceso', 'finalizado', 'cancelado'
	FechaInicioSecado time.Time  `db:"fecha_inicio_secado" json:"fecha_inicio_secado"`
	FechaFinSecado    *time.Time `db:"fecha_fin_secado" json:"fecha_fin_secado"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
}

// LoteListItem es un item en la lista paginada de lotes
type LoteListItem struct {
	ID                int        `json:"id"`
	NombreLote        string     `json:"nombre_lote"`
	Variedad          string     `json:"variedad"`
	TipoProceso       string     `json:"tipo_proceso"`
	PesoKg            float64    `json:"peso_kg"`
	Ubicacion         string     `json:"ubicacion"`
	IDSensor          *int       `json:"id_sensor"`
	CodigoQR          string     `json:"codigo_qr"`
	Estado            string     `json:"estado"`
	FechaInicioSecado time.Time  `json:"fecha_inicio_secado"`
	FechaFinSecado    *time.Time `json:"fecha_fin_secado"`
	CreatedAt         time.Time  `json:"created_at"`
}

// LoteDetalle es la respuesta del detalle de un lote con información adicional del sensor
type LoteDetalle struct {
	LoteCafe
	UltimaTemperatura *float64    `json:"ultima_temperatura"`
	UltimaHumedad     *float64    `json:"ultima_humedad"`
	AlertasActivas    int         `json:"alertas_activas"`
	UltimaPrediccion  *Prediccion `json:"ultima_prediccion"`
}

// LotesListResponse es la respuesta paginada de lotes
type LotesListResponse struct {
	Data   []LoteListItem `json:"data"`
	Total  int            `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}

// CreateLoteRequest es la solicitud de creación de un lote
type CreateLoteRequest struct {
	NombreLote  string  `json:"nombre_lote" validate:"required"`
	Variedad    string  `json:"variedad" validate:"required"`
	TipoProceso string  `json:"tipo_proceso" validate:"required,oneof=lavado honey natural"`
	PesoKg      float64 `json:"peso_kg" validate:"required,gt=0"`
	Ubicacion   string  `json:"ubicacion" validate:"required"`
	IDSensor    *int    `json:"id_sensor"`
}

// UpdateLoteRequest es la solicitud de actualización de un lote
type UpdateLoteRequest struct {
	NombreLote string  `json:"nombre_lote" validate:"required"`
	Variedad   string  `json:"variedad" validate:"required"`
	PesoKg     float64 `json:"peso_kg" validate:"required,gt=0"`
	Ubicacion  string  `json:"ubicacion" validate:"required"`
}

// LecturaAmbiental representa una lectura ambiental de sensor
type LecturaAmbiental struct {
	ID          int       `db:"id" json:"id_lectura"`
	LoteID      int       `db:"lote_id" json:"lote_id"`
	SensorID    int       `db:"sensor_id" json:"sensor_id"`
	Temperatura float64   `db:"temperatura" json:"temperatura"`
	Humedad     float64   `db:"humedad" json:"humedad"`
	Presion     float64   `db:"presion" json:"presion"`
	CreatedAt   time.Time `db:"created_at" json:"timestamp"`
}

// EstadisticasLote estadísticas resumidas de un lote
type EstadisticasLote struct {
	TemperaturaPromedio float64    `json:"temperatura_promedio"`
	TemperaturaMin      float64    `json:"temperatura_min"`
	TemperaturaMax      float64    `json:"temperatura_max"`
	HumedadPromedio     float64    `json:"humedad_promedio"`
	HumedadMin          float64    `json:"humedad_min"`
	HumedadMax          float64    `json:"humedad_max"`
	TotalLecturas       int        `json:"total_lecturas"`
	TotalAlertas        int        `json:"total_alertas"`
	AlertasCriticas     int        `json:"alertas_criticas"`
	AlertasSinAtender   int        `json:"alertas_sin_atender"`
	DiasSecado          int        `json:"dias_secado"`
	UltimaLectura       *time.Time `json:"ultima_lectura"`
}

// Alerta representa una alerta del sistema para un lote
type Alerta struct {
	ID            int        `db:"id" json:"id_alerta"`
	LoteID        int        `db:"lote_id" json:"lote_id"`
	Tipo          string     `db:"tipo" json:"tipo_alerta"`
	Mensaje       string     `db:"mensaje" json:"mensaje"`
	Nivel         string     `db:"nivel" json:"nivel_severidad"`
	Atendida      bool       `db:"atendida" json:"atendida"`
	FechaAtencion *time.Time `db:"fecha_atencion" json:"fecha_atencion"`
	CreatedAt     time.Time  `db:"created_at" json:"fecha_generada"`
	UpdatedAt     time.Time  `db:"updated_at" json:"-"`
}

// Prediccion representa una predicción del ml-service
type Prediccion struct {
	ID                  int       `db:"id" json:"id_prediccion"`
	LoteID              int       `db:"lote_id" json:"lote_id"`
	TiempoEstimadoHoras float64   `db:"tiempo_estimado_horas" json:"tiempo_estimado_horas"`
	CalidadEstimada     string    `db:"calidad_estimada" json:"calidad_estimada"`
	Confianza           float64   `db:"confianza" json:"confianza"`
	FechaPrediccion     time.Time `db:"fecha_prediccion" json:"fecha_prediccion"`
	ModeloVersion       string    `db:"modelo_version" json:"modelo_version"`
}

// Recomendacion representa una recomendación para un lote
type Recomendacion struct {
	ID            int       `db:"id" json:"id_recomendacion"`
	LoteID        int       `db:"lote_id" json:"lote_id"`
	Texto         string    `db:"texto" json:"texto"`
	Origen        string    `db:"origen" json:"origen"`
	FechaGenerada time.Time `db:"fecha_generada" json:"fecha_generada"`
}

// HistorialEvento representa un evento en el historial de un lote
type HistorialEvento struct {
	ID          int       `db:"id" json:"id_evento"`
	LoteID      int       `db:"lote_id" json:"lote_id"`
	Tipo        string    `db:"tipo" json:"tipo_evento"`
	Descripcion string    `db:"descripcion" json:"descripcion"`
	CreatedAt   time.Time `db:"created_at" json:"fecha_evento"`
}

// Reporte representa un reporte solicitado
type Reporte struct {
	ID          int       `db:"id" json:"id"`
	LoteID      int       `db:"lote_id" json:"id_lote"`
	UsuarioID   int       `db:"usuario_id" json:"id_usuario"`
	TipoReporte string    `db:"tipo_reporte" json:"tipo_reporte"`
	Formato     string    `db:"formato" json:"formato"`
	Estado      string    `db:"estado" json:"estado"`
	URLArchivo  *string   `db:"url_archivo" json:"url_archivo"`
	CreatedAt   time.Time `db:"created_at" json:"fecha_generacion"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

// SolicitudReporteRequest solicitud de generación de reporte
type SolicitudReporteRequest struct {
	IDLote      int    `json:"id_lote" validate:"required"`
	TipoReporte string `json:"tipo_reporte" validate:"required"`
	Formato     string `json:"formato" validate:"required,oneof=pdf csv"`
}

// DashboardLoteResumen resumen de un lote para el dashboard
type DashboardLoteResumen struct {
	IDLote            int      `json:"id_lote"`
	NombreLote        string   `json:"nombre_lote"`
	Estado            string   `json:"estado"`
	DiasSecado        int      `json:"dias_secado"`
	UltimaTemperatura *float64 `json:"ultima_temperatura"`
	UltimaHumedad     *float64 `json:"ultima_humedad"`
	AlertasActivas    int      `json:"alertas_activas"`
}

// UltimaPrediccionDashboard predicción más reciente del dashboard
type UltimaPrediccionDashboard struct {
	IDLote              int       `json:"id_lote"`
	NombreLote          string    `json:"nombre_lote"`
	TiempoEstimadoHoras float64   `json:"tiempo_estimado_horas"`
	CalidadEstimada     string    `json:"calidad_estimada"`
	FechaPrediccion     time.Time `json:"fecha_prediccion"`
}

// DashboardResponse respuesta del dashboard del productor
type DashboardResponse struct {
	TotalLotes                int                        `json:"total_lotes"`
	LotesActivos              int                        `json:"lotes_activos"`
	LotesFinalizados          int                        `json:"lotes_finalizados"`
	AlertasSinAtender         int                        `json:"alertas_sin_atender"`
	AlertasCriticasActivas    int                        `json:"alertas_criticas_activas"`
	TemperaturaPromedioActual float64                    `json:"temperatura_promedio_actual"`
	HumedadPromedioActual     float64                    `json:"humedad_promedio_actual"`
	UltimaPrediccion          *UltimaPrediccionDashboard `json:"ultima_prediccion"`
	LotesResumen              []DashboardLoteResumen     `json:"lotes_resumen"`
}

// ErrorResponse es la respuesta de error
type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}
