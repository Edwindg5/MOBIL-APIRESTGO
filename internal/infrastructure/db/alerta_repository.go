package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kajve/api-mobile/internal/application/interfaces"
	"github.com/kajve/api-mobile/internal/domain/entities"
)

type AlertaRepository struct {
	db *PostgresDB
}

func NewAlertaRepository(db *PostgresDB) interfaces.AlertaRepository {
	return &AlertaRepository{db: db}
}

const alertaCols = `id, lote_id, tipo, mensaje, nivel, atendida, fecha_atencion, created_at, updated_at`

func scanAlertaRow(row interface{ Scan(...any) error }, a *entities.Alerta) error {
	return row.Scan(
		&a.ID, &a.LoteID, &a.Tipo, &a.Mensaje, &a.Nivel,
		&a.Atendida, &a.FechaAtencion, &a.CreatedAt, &a.UpdatedAt,
	)
}

func (r *AlertaRepository) GetByID(ctx context.Context, id int) (*entities.Alerta, error) {
	a := &entities.Alerta{}
	err := scanAlertaRow(
		r.db.GetPool().QueryRow(ctx, `SELECT `+alertaCols+` FROM alertas WHERE id = $1`, id), a,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

func (r *AlertaRepository) GetByLoteID(ctx context.Context, loteID int) ([]entities.Alerta, error) {
	return r.GetByLoteIDFiltered(ctx, loteID, nil, "")
}

func (r *AlertaRepository) GetByLoteIDFiltered(ctx context.Context, loteID int, atendida *bool, nivel string) ([]entities.Alerta, error) {
	rows, err := r.db.GetPool().Query(ctx, `
		SELECT `+alertaCols+`
		FROM alertas
		WHERE lote_id = $1
		  AND ($2::boolean IS NULL OR atendida = $2)
		  AND ($3 = '' OR nivel = $3)
		ORDER BY created_at DESC
	`, loteID, atendida, nivel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alertas []entities.Alerta
	for rows.Next() {
		var a entities.Alerta
		if err := scanAlertaRow(rows, &a); err != nil {
			return nil, err
		}
		alertas = append(alertas, a)
	}
	return alertas, rows.Err()
}

func (r *AlertaRepository) MarcarAtendida(ctx context.Context, alertaID int) (*entities.Alerta, error) {
	now := time.Now()
	a := &entities.Alerta{}
	err := scanAlertaRow(r.db.GetPool().QueryRow(ctx, `
		UPDATE alertas
		SET atendida = true, fecha_atencion = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING `+alertaCols, now, alertaID), a)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

func (r *AlertaRepository) Create(ctx context.Context, alerta *entities.Alerta) error {
	return r.db.GetPool().QueryRow(ctx, `
		INSERT INTO alertas (lote_id, tipo, mensaje, nivel, atendida, created_at, updated_at)
		VALUES ($1, $2, $3, $4, false, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`, alerta.LoteID, alerta.Tipo, alerta.Mensaje, alerta.Nivel,
	).Scan(&alerta.ID, &alerta.CreatedAt, &alerta.UpdatedAt)
}
