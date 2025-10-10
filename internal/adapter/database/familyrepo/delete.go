package familyrepo

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

func (fr *FamilyRepository) DeleteFamily(ctx context.Context, tx pgx.Tx, familyID int) error {
	_, err := tx.Exec(ctx, `DELETE FROM users_to_families WHERE family_id = $1`, familyID)
	if err != nil {
		fr.sl.Error("unable to delete from users_to_families", slog.Int("family_id", familyID), slog.String("err", err.Error()))
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM user_bank_tokens WHERE family_id = $1`, familyID)
	if err != nil {
		fr.sl.Error("unable to delete from user_bank_tokens", slog.Int("family_id", familyID), slog.String("err", err.Error()))
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM family_invite_codes WHERE family_id = $1`, familyID)
	if err != nil {
		fr.sl.Error("unable to delete from family_invite_codes", slog.Int("family_id", familyID), slog.String("err", err.Error()))
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM families WHERE id = $1`, familyID)
	if err != nil {
		fr.sl.Error("unable to delete from families", slog.Int("family_id", familyID), slog.String("err", err.Error()))
		return err
	}
	return nil
}
