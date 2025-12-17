package userrepo

import (
	"cashly/internal/entity"
	"context"
	"log/slog"
)

func (ur Repository) Save(ctx context.Context, user *entity.User) (*entity.User, error) {
	q := `INSERT INTO users (id, username, firstname)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE 
    		SET id = EXCLUDED.id
			RETURNING id, username, firstname, joined_at;`

	u := new(entity.User)

	err := ur.db.Pool().QueryRow(ctx, q, user.ID, user.Username, user.Firstname).Scan(&u.ID, &u.Username, &u.Firstname, &u.JoinedAt)
	if err != nil {
		ur.logger.Error("failed to save user", slog.String("err", err.Error()))
	}

	return u, err
}

func (ur Repository) SaveToFamily(ctx context.Context, familyID int, userID int64) error {
	q := `INSERT INTO users_to_families (user_id, family_id)
			VALUES ($1, $2)`

	_, err := ur.db.Pool().Exec(ctx, q, userID, familyID)
	if err != nil {
		ur.logger.Error("failed to create family",
			slog.String("err", err.Error()),
			slog.Int("family_id", familyID))
	}

	return err
}
