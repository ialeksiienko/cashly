package userrepo

import (
	"cashly/internal/entity"
	"context"
	"log/slog"
)

func (ur *UserRepository) GetAllUsersInFamily(ctx context.Context, familyID int) ([]entity.User, error) {
	q := `SELECT u.id, u.username, u.firstname, u.joined_at
	FROM users_to_families utf
	JOIN users u ON u.id = utf.user_id
	WHERE utf.family_id = $1`

	rows, err := ur.db.Query(ctx, q, familyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Firstname, &u.JoinedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		ur.sl.Error("failed to get family by user id", slog.String("err", err.Error()))
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	q := `SELECT id, username, firstname, joined_at
		FROM users WHERE id = $1 `

	u := new(entity.User)

	err := ur.db.QueryRow(ctx, q, id).Scan(&u.ID, &u.Username, &u.Firstname, &u.JoinedAt)
	if err != nil {
		ur.sl.Error("failed to get family", slog.String("err", err.Error()))
		return nil, err
	}

	return u, nil
}
