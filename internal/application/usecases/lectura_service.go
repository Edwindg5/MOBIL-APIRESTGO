package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kajve/api-mobile/internal/application/interfaces"
	"github.com/kajve/api-mobile/internal/domain/entities"
)

type LecturaService struct {
	lecturaRepository interfaces.LecturaRepository
	loteRepository    interfaces.LoteRepository
}

func NewLecturaService(
	lecturaRepository interfaces.LecturaRepository,
	loteRepository interfaces.LoteRepository,
) interfaces.LecturaService {
	return &LecturaService{
		lecturaRepository: lecturaRepository,
		loteRepository:    loteRepository,
	}
}

func (s *LecturaService) verifyOwnership(ctx context.Context, loteID, usuarioID int) error {
	lote, err := s.loteRepository.GetByID(ctx, loteID)
	if err != nil {
		return fmt.Errorf("error getting lote: %w", err)
	}
	if lote == nil {
		return errors.New("lote not found")
	}
	if lote.UsuarioID != usuarioID {
		return errors.New("unauthorized")
	}
	return nil
}

func (s *LecturaService) GetLatestLecturas(ctx context.Context, loteID, usuarioID int, limit int) ([]entities.LecturaAmbiental, error) {
	if err := s.verifyOwnership(ctx, loteID, usuarioID); err != nil {
		return nil, err
	}
	lecturas, err := s.lecturaRepository.GetLatestByLoteID(ctx, loteID, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting lecturas: %w", err)
	}
	return lecturas, nil
}

func (s *LecturaService) GetLecturasFiltradas(ctx context.Context, loteID, usuarioID int, limit int, desde time.Time) ([]entities.LecturaAmbiental, error) {
	if err := s.verifyOwnership(ctx, loteID, usuarioID); err != nil {
		return nil, err
	}
	lecturas, err := s.lecturaRepository.GetByLoteIDFiltered(ctx, loteID, limit, desde)
	if err != nil {
		return nil, fmt.Errorf("error getting lecturas: %w", err)
	}
	return lecturas, nil
}

func (s *LecturaService) GetEstadisticas(ctx context.Context, loteID, usuarioID int) (*entities.EstadisticasLote, error) {
	if err := s.verifyOwnership(ctx, loteID, usuarioID); err != nil {
		return nil, err
	}
	stats, err := s.lecturaRepository.GetEstadisticas(ctx, loteID)
	if err != nil {
		return nil, fmt.Errorf("error getting estadisticas: %w", err)
	}
	return stats, nil
}
