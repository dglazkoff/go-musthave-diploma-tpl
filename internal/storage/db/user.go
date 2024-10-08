package db

import (
	"context"
	"database/sql"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
)

func (s *dbStorage) GetUserByLogin(ctx context.Context, login string) (user models.User, err error) {
	row := s.db.QueryRowContext(ctx, "SELECT login, password, balance from users WHERE login = $1", login)
	err = row.Scan(&user.Login, &user.Password, &user.Balance)

	return
}

func (s *dbStorage) GetUserByLoginTx(ctx context.Context, tx *sql.Tx, login string) (user models.User, err error) {
	row := tx.QueryRowContext(ctx, "SELECT login, password, balance from users WHERE login = $1", login)
	err = row.Scan(&user.Login, &user.Password, &user.Balance)

	return
}

// как в таких случаях поступать? дублировать функцию просто чтобы добавить FOR UPDATE?
// или делать прям в storage слое методы UpdateBalance?
// или еще как то?
func (s *dbStorage) GetUserByLoginForUpdate(ctx context.Context, tx *sql.Tx, login string) (user models.User, err error) {
	row := tx.QueryRowContext(ctx, "SELECT login, password, balance from users WHERE login = $1 FOR UPDATE", login)
	err = row.Scan(&user.Login, &user.Password, &user.Balance)

	return
}

func (s *dbStorage) CreateUser(ctx context.Context, tx *sql.Tx, login, password string) error {
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO users (login, password, balance) VALUES($1, $2, $3) ON CONFLICT (login) DO NOTHING",
		login, password, 0,
	)

	if err != nil {
		logger.Log.Debug("error while creating user ", err)
		return err
	}

	return nil
}

func (s *dbStorage) UpdateUser(ctx context.Context, tx *sql.Tx, user models.User) (models.User, error) {
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO users (login, password, balance) VALUES($1, $2, $3) ON CONFLICT (login) DO UPDATE SET balance = $3",
		user.Login, user.Password, user.Balance,
	)

	if err != nil {
		logger.Log.Debug("error while updating user ", err)
		return models.User{}, err
	}

	return user, nil
}
