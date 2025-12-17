package family

import (
	"cashly/internal/entity"
	"context"
	"log/slog"
)

func (fr Repository) Create(ctx context.Context, inp *entity.Family) (*entity.Family, error) {
	q := `INSERT INTO families (created_by, name) 
			VALUES ($1, $2) RETURNING id, created_by, name`

	f := new(entity.Family)

	err := fr.db.Pool().QueryRow(ctx, q, inp.CreatedBy, inp.Name).Scan(&f.ID, &f.CreatedBy, &f.Name)
	if err != nil {
		fr.logger.Error("failed to create family", slog.String("err", err.Error()), slog.String("family", inp.Name))
	}

	return f, err
}
