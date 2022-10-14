package migration

import (
	"uwwolf/app/enum"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration02 = &gormigrate.Migration{
	ID: "2",
	Migrate: func(tx *gorm.DB) error {
		return tx.Exec(`
            CREATE OR REPLACE FUNCTION update_user_status_to_in_game()
            RETURNS TRIGGER AS
            $$
            BEGIN
                UPDATE users SET status_id = ? WHERE id = NEW.user_id;

                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;

            CREATE trigger player_play_new_game
            AFTER INSERT ON role_assignments
            FOR EACH ROW
            EXECUTE PROCEDURE update_user_status_to_in_game();
        `, enum.InGameStatus).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Exec(`
            DROP TRIGGER IF EXISTS player_play_new_game ON role_assignments;
            DROP FUNCTION IF EXISTS update_user_status_to_in_game;
        `).Error
	},
}
