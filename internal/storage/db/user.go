package db

import (
	"context"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
)

func (d *dbStorage) GetUserByLogin(ctx context.Context, login string) (user models.User, err error) {
	row := d.db.QueryRowContext(ctx, "SELECT login, password, balance from users WHERE login = $1", login)
	err = row.Scan(&user.Login, &user.Password, &user.Balance)

	return
}

func (d *dbStorage) CreateUser(ctx context.Context, login, password string) error {
	_, err := d.db.ExecContext(
		ctx,
		"INSERT INTO users (login, password, balance) VALUES($1, $2) ON CONFLICT (login) DO NOTHING",
		login, password, 0,
	)

	if err != nil {
		logger.Log.Debug("error while creating user ", err)
		return err
	}

	return nil
}

func (d *dbStorage) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	_, err := d.db.ExecContext(
		ctx,
		"UPDATE users SET login = $1, password = $2, balance = $3 WHERE login = $3 ON CONFLICT (login) DO UPDATE SET balance = $3",
		user.Login, user.Password, user.Balance,
	)

	if err != nil {
		logger.Log.Debug("error while updating user ", err)
		return models.User{}, err
	}

	return user, nil
}
