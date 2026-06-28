package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/kajve/api-mobile/internal/application/interfaces"
	"github.com/kajve/api-mobile/internal/domain/entities"
)

type HistorialService struct {
	historialRepository interfaces.HistorialRepository
	loteRepository      interfaces.LoteRepository
}

func NewHistorialService(
	historialRepository interfaces.HistorialRepository,
	loteRepository interfaces.LoteRepository,
) interfaces.HistorialService {
	return &HistorialService{
		historialRepository: historialRepository,
		loteRepository:      loteRepository,
	}
}

func (s *HistorialService) GetHistorial(ctx context.Context, loteID, usuarioID int, page, limit int) ([]entities.HistorialEvento, int, error) {
	lote, err := s.loteRepository.GetByID(ctx, loteID)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting lote: %w", err)
	}
	if lote == nil {
		return nil, 0, errors.New("lote not found")
	}
	if lote.UsuarioID != usuarioID {
		return nil, 0, errors.New("unauthorized")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	eventos, total, err := s.historialRepository.GetByLoteIDPaginated(ctx, loteID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting historial: %w", err)
	}
	return eventos, total, nil
}
