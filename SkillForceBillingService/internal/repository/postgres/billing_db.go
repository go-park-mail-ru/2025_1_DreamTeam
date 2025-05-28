package postgres

import (
	"context"
	"database/sql"
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

	var userID, courseID int64

	query = `SELECT User_ID, Course_ID FROM PURCHACES WHERE Billing_ID = $1 AND Status = 'success'`
	err = d.conn.QueryRow(query, billing_id).Scan(&userID, &courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no successful purchase found with Billing_ID: %s", billing_id)
		}
		return err
	}

	insertQuery := `
		INSERT INTO SIGNUPS (User_ID, Course_ID)
		VALUES ($1, $2)
		ON CONFLICT (Course_ID, User_ID) DO NOTHING
	`
	_, err = d.conn.Exec(insertQuery, userID, courseID)
	if err != nil {
		return fmt.Errorf("failed to insert into SIGNUPS: %w", err)
	}

	return nil

}

func (d *Database) GetBillingInfo(ctx context.Context, courseID int) (string, int, error) {
	var title string
	var price int

	query := `SELECT Title, Price FROM COURSE WHERE ID = $1`
	err := d.conn.QueryRow(query, courseID).Scan(&title, &price)
	if err != nil {
		return "", 0, fmt.Errorf("error executing select: %w", err)
	}
	return title, price, nil
}
