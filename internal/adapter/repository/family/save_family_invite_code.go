package family

import (
	"context"
	"log/slog"
	"time"
)

func (fr Repository) SaveInviteCode(ctx context.Context, userID int64, familyID int, code string) (time.Time, error) {
	q := `INSERT INTO family_invite_codes 
    	(family_id, code, created_by, expires_at)
    	VALUES ($1, $2, $3, $4)
    	RETURNING expires_at`

	var expiresAt time.Time

	err := fr.db.Pool().QueryRow(ctx, q, familyID, code, userID, time.Now().UTC().Add(48*time.Hour)).Scan(&expiresAt)
	if err != nil {
		fr.logger.Error("failed to save family invite code", slog.String("err", err.Error()),
			slog.Int("user_id", int(userID)), slog.Int("family_id", familyID))
	}

	return expiresAt, err
}
