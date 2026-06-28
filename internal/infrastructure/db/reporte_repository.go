package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/kajve/api-mobile/internal/application/interfaces"
	"github.com/kajve/api-mobile/internal/domain/entities"
)

type ReporteRepository struct {
	db *PostgresDB
}

func NewReporteRepository(db *PostgresDB) interfaces.ReporteRepository {
	return &ReporteRepository{db: db}
}

const reporteCols = `id, lote_id, usuario_id, tipo_reporte, formato, estado, url_archivo, created_at, updated_at`

func scanReporte(row interface{ Scan(...any) error }, rep *entities.Reporte) error {
	return row.Scan(
		&rep.ID, &rep.LoteID, &rep.UsuarioID, &rep.TipoReporte,
		&rep.Formato, &rep.Estado, &rep.URLArchivo, &rep.CreatedAt, &rep.UpdatedAt,
	)
}

func (r *ReporteRepository) GetByID(ctx context.Context, id int) (*entities.Reporte, error) {
	rep := &entities.Reporte{}
	err := scanReporte(r.db.GetPool().QueryRow(ctx,
		`SELECT `+reporteCols+` FROM reportes WHERE id = $1`, id), rep)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return rep, nil
}

func (r *ReporteRepository) GetByUsuarioID(ctx context.Context, usuarioID int) ([]entities.Reporte, error) {
	rows, err := r.db.GetPool().Query(ctx, `
		SELECT `+reporteCols+`
		FROM reportes
		WHERE usuario_id = $1
		ORDER BY created_at DESC
	`, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reportes []entities.Reporte
	for rows.Next() {
		var rep entities.Reporte
		if err := scanReporte(rows, &rep); err != nil {
			return nil, err
		}
		reportes = append(reportes, rep)
	}
	return reportes, rows.Err()
}

func (r *ReporteRepository) Create(ctx context.Context, rep *entities.Reporte) error {
	return r.db.GetPool().QueryRow(ctx, `
		INSERT INTO reportes (lote_id, usuario_id, tipo_reporte, formato, estado, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`, rep.LoteID, rep.UsuarioID, rep.TipoReporte, rep.Formato, rep.Estado,
	).Scan(&rep.ID, &rep.CreatedAt, &rep.UpdatedAt)
}

func (r *ReporteRepository) Update(ctx context.Context, rep *entities.Reporte) error {
	_, err := r.db.GetPool().Exec(ctx, `
		UPDATE reportes
		SET estado = $1, url_archivo = $2, updated_at = NOW()
		WHERE id = $3
	`, rep.Estado, rep.URLArchivo, rep.ID)
	return err
}
