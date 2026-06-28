package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kajve/api-mobile/internal/application/interfaces"
	"github.com/kajve/api-mobile/internal/domain/entities"
)

type LoteRepository struct {
	db *PostgresDB
}

// NewLoteRepository crea una nueva instancia del repositorio
func NewLoteRepository(db *PostgresDB) interfaces.LoteRepository {
	return &LoteRepository{db: db}
}

const loteColumns = `id, usuario_id, nombre_lote, variedad, tipo_proceso, peso_kg, ubicacion,
	id_sensor, codigo_qr, estado, fecha_inicio_secado, fecha_fin_secado, created_at, updated_at`

// GetByID obtiene un lote por ID
func (r *LoteRepository) GetByID(ctx context.Context, id int) (*entities.LoteCafe, error) {
	lote := &entities.LoteCafe{}
	err := r.db.GetPool().QueryRow(ctx, `
		SELECT `+loteColumns+`
		FROM lotes_cafe WHERE id = $1
	`, id).Scan(
		&lote.ID, &lote.UsuarioID, &lote.NombreLote, &lote.Variedad, &lote.TipoProceso,
		&lote.PesoKg, &lote.Ubicacion, &lote.IDSensor, &lote.CodigoQR, &lote.Estado,
		&lote.FechaInicioSecado, &lote.FechaFinSecado, &lote.CreatedAt, &lote.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return lote, nil
}

// GetByUsuarioID obtiene lotes de un usuario con filtro de estado y paginación opcional.
// limit=0 devuelve todos los registros sin paginación.
func (r *LoteRepository) GetByUsuarioID(ctx context.Context, usuarioID int, estado string, limit, offset int) ([]entities.LoteCafe, int, error) {
	var total int
	err := r.db.GetPool().QueryRow(ctx, `
		SELECT COUNT(*) FROM lotes_cafe
		WHERE usuario_id = $1 AND (estado = $2 OR $2 = '')
	`, usuarioID, estado).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	var (
		rows interface {
			Next() bool
			Scan(...any) error
			Close()
			Err() error
		}
		queryErr error
	)

	if limit > 0 {
		rows, queryErr = r.db.GetPool().Query(ctx, `
			SELECT `+loteColumns+`
			FROM lotes_cafe
			WHERE usuario_id = $1 AND (estado = $2 OR $2 = '')
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`, usuarioID, estado, limit, offset)
	} else {
		rows, queryErr = r.db.GetPool().Query(ctx, `
			SELECT `+loteColumns+`
			FROM lotes_cafe
			WHERE usuario_id = $1 AND (estado = $2 OR $2 = '')
			ORDER BY created_at DESC
		`, usuarioID, estado)
	}
	if queryErr != nil {
		return nil, 0, queryErr
	}
	defer rows.Close()

	var lotes []entities.LoteCafe
	for rows.Next() {
		var l entities.LoteCafe
		if err := rows.Scan(
			&l.ID, &l.UsuarioID, &l.NombreLote, &l.Variedad, &l.TipoProceso,
			&l.PesoKg, &l.Ubicacion, &l.IDSensor, &l.CodigoQR, &l.Estado,
			&l.FechaInicioSecado, &l.FechaFinSecado, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		lotes = append(lotes, l)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return lotes, total, nil
}

// Create inserta un nuevo lote; el codigo_qr se genera con gen_random_uuid() en la BD
func (r *LoteRepository) Create(ctx context.Context, lote *entities.LoteCafe) (*entities.LoteCafe, error) {
	tx, err := r.db.BeginTx(ctx, lote.UsuarioID)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	created := &entities.LoteCafe{}
	err = tx.QueryRow(ctx, `
		INSERT INTO lotes_cafe
			(usuario_id, nombre_lote, variedad, tipo_proceso, peso_kg, ubicacion,
			 id_sensor, codigo_qr, estado, fecha_inicio_secado, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, gen_random_uuid()::text, $8, NOW(), NOW(), NOW())
		RETURNING `+loteColumns,
		lote.UsuarioID, lote.NombreLote, lote.Variedad, lote.TipoProceso,
		lote.PesoKg, lote.Ubicacion, lote.IDSensor, lote.Estado,
	).Scan(
		&created.ID, &created.UsuarioID, &created.NombreLote, &created.Variedad, &created.TipoProceso,
		&created.PesoKg, &created.Ubicacion, &created.IDSensor, &created.CodigoQR, &created.Estado,
		&created.FechaInicioSecado, &created.FechaFinSecado, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return created, nil
}

// Update modifica nombre_lote, variedad, peso_kg y ubicacion de un lote en estado 'en_proceso'
func (r *LoteRepository) Update(ctx context.Context, loteID, usuarioID int, nombre, variedad string, pesoKg float64, ubicacion string) (*entities.LoteCafe, error) {
	tx, err := r.db.BeginTx(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	lote := &entities.LoteCafe{}
	err = tx.QueryRow(ctx, `
		UPDATE lotes_cafe
		SET nombre_lote = $1, variedad = $2, peso_kg = $3, ubicacion = $4, updated_at = NOW()
		WHERE id = $5 AND usuario_id = $6 AND estado = 'en_proceso'
		RETURNING `+loteColumns,
		nombre, variedad, pesoKg, ubicacion, loteID, usuarioID,
	).Scan(
		&lote.ID, &lote.UsuarioID, &lote.NombreLote, &lote.Variedad, &lote.TipoProceso,
		&lote.PesoKg, &lote.Ubicacion, &lote.IDSensor, &lote.CodigoQR, &lote.Estado,
		&lote.FechaInicioSecado, &lote.FechaFinSecado, &lote.CreatedAt, &lote.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return lote, nil
}

// UpdateEstado cambia el estado de un lote en 'en_proceso'; si fechaFin != nil asigna fecha_fin_secado
func (r *LoteRepository) UpdateEstado(ctx context.Context, loteID, usuarioID int, estado string, fechaFin *time.Time) (*entities.LoteCafe, error) {
	tx, err := r.db.BeginTx(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	lote := &entities.LoteCafe{}
	err = tx.QueryRow(ctx, `
		UPDATE lotes_cafe
		SET estado = $1, fecha_fin_secado = $2, updated_at = NOW()
		WHERE id = $3 AND usuario_id = $4 AND estado = 'en_proceso'
		RETURNING `+loteColumns,
		estado, fechaFin, loteID, usuarioID,
	).Scan(
		&lote.ID, &lote.UsuarioID, &lote.NombreLote, &lote.Variedad, &lote.TipoProceso,
		&lote.PesoKg, &lote.Ubicacion, &lote.IDSensor, &lote.CodigoQR, &lote.Estado,
		&lote.FechaInicioSecado, &lote.FechaFinSecado, &lote.CreatedAt, &lote.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return lote, nil
}
