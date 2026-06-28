package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/kajve/api-mobile/internal/application/interfaces"
	"github.com/kajve/api-mobile/internal/domain/entities"
)

type UsuarioRepository struct {
	db *PostgresDB
}

// NewUsuarioRepository crea una nueva instancia del repositorio
func NewUsuarioRepository(db *PostgresDB) interfaces.UsuarioRepository {
	return &UsuarioRepository{db: db}
}

// GetByEmail obtiene un usuario activo por email (sin RLS, uso en login)
func (r *UsuarioRepository) GetByEmail(ctx context.Context, email string) (*entities.Usuario, error) {
	usuario := &entities.Usuario{}

	err := r.db.GetPool().QueryRow(ctx, `
		SELECT id, email, password, nombre_completo, telefono, rol, estado, created_at, updated_at
		FROM usuarios
		WHERE email = $1 AND estado = 'activo'
	`, email).Scan(
		&usuario.ID,
		&usuario.Email,
		&usuario.Password,
		&usuario.NombreCompleto,
		&usuario.Telefono,
		&usuario.Rol,
		&usuario.Estado,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return usuario, nil
}

// GetByID obtiene un usuario por ID
func (r *UsuarioRepository) GetByID(ctx context.Context, id int) (*entities.Usuario, error) {
	usuario := &entities.Usuario{}

	err := r.db.GetPool().QueryRow(ctx, `
		SELECT id, email, password, nombre_completo, telefono, rol, estado, created_at, updated_at
		FROM usuarios
		WHERE id = $1
	`, id).Scan(
		&usuario.ID,
		&usuario.Email,
		&usuario.Password,
		&usuario.NombreCompleto,
		&usuario.Telefono,
		&usuario.Rol,
		&usuario.Estado,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return usuario, nil
}

// ExistsByEmail comprueba si ya existe algún usuario con ese email, independientemente del estado
func (r *UsuarioRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.GetPool().QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM usuarios WHERE email = $1)
	`, email).Scan(&exists)
	return exists, err
}

// Create crea un nuevo usuario
func (r *UsuarioRepository) Create(ctx context.Context, usuario *entities.Usuario) error {
	err := r.db.GetPool().QueryRow(ctx, `
		INSERT INTO usuarios (email, password, nombre_completo, telefono, rol, estado, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`, usuario.Email, usuario.Password, usuario.NombreCompleto, usuario.Telefono, usuario.Rol, usuario.Estado).Scan(
		&usuario.ID,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)

	return err
}

// Update actualiza nombre y telefono del usuario; usa transacción con RLS
func (r *UsuarioRepository) Update(ctx context.Context, id int, nombre, telefono string) (*entities.Usuario, error) {
	tx, err := r.db.BeginTx(ctx, id)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	usuario := &entities.Usuario{}
	err = tx.QueryRow(ctx, `
		UPDATE usuarios
		SET nombre_completo = $1, telefono = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, email, nombre_completo, telefono, rol, estado, created_at, updated_at
	`, nombre, telefono, id).Scan(
		&usuario.ID,
		&usuario.Email,
		&usuario.NombreCompleto,
		&usuario.Telefono,
		&usuario.Rol,
		&usuario.Estado,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
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
	return usuario, nil
}

// UpdatePassword actualiza la contraseña hasheada del usuario con RLS
func (r *UsuarioRepository) UpdatePassword(ctx context.Context, id int, hashedPassword string) error {
	_, err := r.db.Exec(ctx, id, `
		UPDATE usuarios SET password = $1, updated_at = NOW() WHERE id = $2
	`, hashedPassword, id)
	return err
}
