package familyrepo

import (
	"log/slog"

	"context"

	"github.com/jackc/pgx/v4"

	"time"
)

func (fr *FamilyRepository) GetInviteCode(ctx context.Context, familyID int) (string, time.Time, error) {
	q := `
	SELECT code, expires_at 
	FROM family_invite_codes
	WHERE family_id = $1 AND expires_at > NOW()`

	var code string
	var expiresAt time.Time

	err := fr.db.QueryRow(ctx, q, familyID).Scan(&code, &expiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", time.Time{}, nil
		}
		fr.sl.Error("unable to get family invite code", slog.Int("famliy_id", familyID))
		return "", time.Time{}, err
	}

	return code, expiresAt, nil
}
