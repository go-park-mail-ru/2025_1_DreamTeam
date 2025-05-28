package postgres

import (
	"context"
	"fmt"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"
)

func (d *Database) GetUserById(ctx context.Context, userId int) (*usermodels.User, error) {
	var user usermodels.User
	err := d.conn.QueryRow("SELECT email, name, hide_email FROM usertable WHERE id = $1", userId).Scan(&user.Email, &user.Name, &user.HideEmail)
	if err != nil {
		logs.PrintLog(ctx, "GetUserById", fmt.Sprintf("%+v", err))
		return nil, err
	}
	return &user, nil
}
