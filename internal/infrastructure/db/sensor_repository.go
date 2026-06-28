package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kajve/api-mobile/internal/application/interfaces"
	"github.com/kajve/api-mobile/internal/domain/entities"
)

type SensorRepository struct {
	db *PostgresDB
}

func NewSensorRepository(db *PostgresDB) interfaces.SensorRepository {
	return &SensorRepository{db: db}
}

const sensorCols = `id, esp32_id, mac_address, lote_id, provisioning_token, token_usado, linked_at, last_seen, estado, created_at, updated_at`

func scanSensor(row interface{ Scan(...any) error }, s *entities.Sensor) error {
	return row.Scan(
		&s.ID, &s.ESP32ID, &s.MacAddress, &s.LoteID,
		&s.ProvisioningToken, &s.TokenUsado,
		&s.LinkedAt, &s.LastSeen, &s.Estado, &s.CreatedAt, &s.UpdatedAt,
	)
}

func (r *SensorRepository) GetByESP32ID(ctx context.Context, esp32ID string) (*entities.Sensor, error) {
	return r.GetByIdentifier(ctx, esp32ID)
}

func (r *SensorRepository) GetByIdentifier(ctx context.Context, identifier string) (*entities.Sensor, error) {
	s := &entities.Sensor{}
	err := scanSensor(r.db.GetPool().QueryRow(ctx, `
		SELECT `+sensorCols+`
		FROM sensores
		WHERE esp32_id = $1 OR mac_address = $1
		LIMIT 1
	`, identifier), s)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *SensorRepository) GetByID(ctx context.Context, id int) (*entities.Sensor, error) {
	s := &entities.Sensor{}
	err := scanSensor(r.db.GetPool().QueryRow(ctx,
		`SELECT `+sensorCols+` FROM sensores WHERE id = $1`, id), s)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *SensorRepository) Create(ctx context.Context, sensor *entities.Sensor) (int, error) {
	var id int
	err := r.db.GetPool().QueryRow(ctx, `
		INSERT INTO sensores (esp32_id, mac_address, lote_id, estado, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`, sensor.ESP32ID, sensor.MacAddress, sensor.LoteID, sensor.Estado).Scan(&id)
	return id, err
}

func (r *SensorRepository) LinkToLote(ctx context.Context, sensorID, loteID int) error {
	tag, err := r.db.GetPool().Exec(ctx, `
		UPDATE sensores
		SET lote_id = $1, linked_at = NOW(), estado = 'activo', updated_at = NOW()
		WHERE id = $2
	`, loteID, sensorID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sensor %d not found", sensorID)
	}
	return nil
}

func (r *SensorRepository) MarcarTokenUsado(ctx context.Context, sensorID int) error {
	_, err := r.db.GetPool().Exec(ctx, `
		UPDATE sensores SET token_usado = true, updated_at = NOW() WHERE id = $1
	`, sensorID)
	return err
}
