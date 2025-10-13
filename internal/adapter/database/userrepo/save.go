package userrepo

import (
	"context"
	"log/slog"
	"monofamily/internal/entity"
)

func (ur *UserRepository) Save(ctx context.Context, user *entity.User) (*entity.User, error) {
	q := `INSERT INTO users (id, username, firstname)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE 
    		SET id = EXCLUDED.id
			RETURNING id, username, firstname, joined_at;`

	u := new(entity.User)

	err := ur.db.QueryRow(ctx, q, user.ID, user.Username, user.Firstname).Scan(&u.ID, &u.Username, &u.Firstname, &u.JoinedAt)
	if err != nil {
		ur.sl.Error("failed to save user", slog.String("error", err.Error()))
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) SaveUserToFamily(ctx context.Context, familyID int, userID int64) error {
	q := `INSERT INTO users_to_families (user_id, family_id)
			VALUES ($1, $2)`

	_, err := ur.db.Exec(ctx, q, userID, familyID)
	if err != nil {
		ur.sl.Error("failed to create family",
			slog.String("err", err.Error()),
			slog.Int("family_id", familyID))
		return err
	}

	return nil
}
