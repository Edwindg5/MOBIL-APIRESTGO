package db

import (
	"context"

	"github.com/kajve/api-mobile/internal/application/interfaces"
	"github.com/kajve/api-mobile/internal/domain/entities"
)

type HistorialRepository struct {
	db *PostgresDB
}

func NewHistorialRepository(db *PostgresDB) interfaces.HistorialRepository {
	return &HistorialRepository{db: db}
}

func (r *HistorialRepository) GetByLoteID(ctx context.Context, loteID int) ([]entities.HistorialEvento, error) {
	eventos, _, err := r.GetByLoteIDPaginated(ctx, loteID, 0, 0)
	return eventos, err
}

func (r *HistorialRepository) GetByLoteIDPaginated(ctx context.Context, loteID int, limit, offset int) ([]entities.HistorialEvento, int, error) {
	var total int
	if err := r.db.GetPool().QueryRow(ctx,
		`SELECT COUNT(*) FROM historial_eventos WHERE lote_id = $1`, loteID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	var rows interface {
		Next() bool
		Scan(...any) error
		Close()
		Err() error
	}
	var err error

	if limit > 0 {
		rows, err = r.db.GetPool().Query(ctx, `
			SELECT id, lote_id, tipo, descripcion, created_at
			FROM historial_eventos
			WHERE lote_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, loteID, limit, offset)
	} else {
		rows, err = r.db.GetPool().Query(ctx, `
			SELECT id, lote_id, tipo, descripcion, created_at
			FROM historial_eventos
			WHERE lote_id = $1
			ORDER BY created_at DESC
		`, loteID)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var eventos []entities.HistorialEvento
	for rows.Next() {
		var e entities.HistorialEvento
		if err := rows.Scan(&e.ID, &e.LoteID, &e.Tipo, &e.Descripcion, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		eventos = append(eventos, e)
	}
	return eventos, total, rows.Err()
}

func (r *HistorialRepository) Create(ctx context.Context, evento *entities.HistorialEvento) error {
	return r.db.GetPool().QueryRow(ctx, `
		INSERT INTO historial_eventos (lote_id, tipo, descripcion, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`, evento.LoteID, evento.Tipo, evento.Descripcion,
	).Scan(&evento.ID, &evento.CreatedAt)
}
