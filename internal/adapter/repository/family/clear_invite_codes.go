package family

import (
	"context"
	"log/slog"
)

func (fr Repository) ClearInviteCodes(ctx context.Context) error {
	_, err := fr.db.Pool().Exec(ctx, `
        DELETE FROM family_invite_codes
        WHERE expires_at < NOW()
    `)
	if err != nil {
		fr.logger.Error("failed to delete expired invite codes",
			slog.String("err", err.Error()))
	}
	return err
}
