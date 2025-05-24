package postgres

import (
	"context"
	"fmt"
)

func (d *Database) AddNewBilling(ctx context.Context, userID int, courseID int, billing_id string) error {
	query := `
	INSERT INTO PURCHACES (User_ID, Course_ID, Status, Billing_ID)
	VALUES ($1, $2, 'pending', $3)
	RETURNING ID;
`

	var insertedID int64
	err := d.conn.QueryRow(query, userID, courseID, billing_id).Scan(&insertedID)
	if err != nil {
		panic(err)
	}

	return nil
}

func (d *Database) UpdateBilling(ctx context.Context, billing_id string) error {
	query := `
	UPDATE PURCHACES
	SET Status = $1,
		Updated_at = CURRENT_TIMESTAMP
	WHERE Billing_ID = $2;
`

	_, err := d.conn.Exec(query, "success", billing_id)
	if err != nil {
		return fmt.Errorf("error executing update: %w", err)
	}

	return nil

}
