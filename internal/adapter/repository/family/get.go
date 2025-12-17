package family

import (
	"cashly/internal/entity"
	"context"
	"log/slog"
	"time"
)

func (fr Repository) GetFamiliesByUserID(ctx context.Context, userID int64) ([]entity.Family, error) {
	q := `SELECT f.id, f.name 
	FROM users_to_families utf
	JOIN families f ON f.id = utf.family_id
	WHERE utf.user_id = $1`

	rows, err := fr.db.Pool().Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var families []entity.Family
	for rows.Next() {
		var f entity.Family
		if err := rows.Scan(&f.ID, &f.Name); err != nil {
			return nil, err
		}
		families = append(families, f)
	}

	if err := rows.Err(); err != nil {
		fr.logger.Error("failed to get family by user id", slog.String("err", err.Error()))
	}

	return families, err
}

func (fr Repository) GetByCode(ctx context.Context, code string) (*entity.Family, time.Time, error) {
	q := `SELECT f.id, f.created_by, f.name, fi.expires_at
		FROM family_invite_codes fi
		JOIN families f ON f.id = fi.family_id
		WHERE fi.code = $1`

	f := new(entity.Family)

	var expiresAt time.Time

	err := fr.db.Pool().QueryRow(ctx, q, code).Scan(&f.ID, &f.CreatedBy, &f.Name, &expiresAt)
	if err != nil {
		fr.logger.Error("failed to get family", slog.String("err", err.Error()))
	}

	return f, expiresAt, err
}

func (fr Repository) GetByID(ctx context.Context, id int) (*entity.Family, error) {
	q := `SELECT id, created_by, name FROM families WHERE id = $1`

	f := new(entity.Family)
	err := fr.db.Pool().QueryRow(ctx, q, id).Scan(&f.ID, &f.CreatedBy, &f.Name)
	if err != nil {
		fr.logger.Error("failed to get family", slog.String("err", err.Error()))
	}

	return f, err
}
