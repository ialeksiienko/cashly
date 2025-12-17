package family

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

func (fr Repository) Delete(ctx context.Context, tx pgx.Tx, id int) error {
	_, err := tx.Exec(ctx, `DELETE FROM users_to_families WHERE family_id = $1`, id)
	if err != nil {
		fr.logger.Error("unable to delete from users_to_families", slog.Int("family_id", id), slog.String("err", err.Error()))
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM user_bank_tokens WHERE family_id = $1`, id)
	if err != nil {
		fr.logger.Error("unable to delete from user_bank_tokens", slog.Int("family_id", id), slog.String("err", err.Error()))
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM family_invite_codes WHERE family_id = $1`, id)
	if err != nil {
		fr.logger.Error("unable to delete from family_invite_codes", slog.Int("family_id", id), slog.String("err", err.Error()))
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM families WHERE id = $1`, id)
	if err != nil {
		fr.logger.Error("unable to delete from families", slog.Int("family_id", id), slog.String("err", err.Error()))
	}

	return err
}
