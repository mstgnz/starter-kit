package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mstgnz/starter-kit/internal/auth"
	"github.com/mstgnz/starter-kit/internal/config"
	"github.com/mstgnz/starter-kit/model"
)

type UserRepository struct {
}

func (r *UserRepository) Count(ctx context.Context) int {
	rowCount := 0

	// prepare count
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USERS_COUNT"])
	if err != nil {
		return rowCount
	}

	// query
	rows, err := stmt.Query()
	if err != nil {
		return rowCount
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&rowCount); err != nil {
			return rowCount
		}
	}

	return rowCount
}

func (r *UserRepository) Get(ctx context.Context, offset, limit int, search string) []*model.User {
	users := []*model.User{}

	// prepare users paginate
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USERS_PAGINATE"])
	if err != nil {
		return users
	}

	// query
	rows, err := stmt.Query("%"+search+"%", offset, limit)
	if err != nil {
		return users
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Fullname, &user.Email, &user.Password, &user.Phone, &user.IsAdmin, &user.Active, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
			return users
		}
		users = append(users, user)
	}

	return users
}

func (r *UserRepository) Create(ctx context.Context, register *model.Register) (*model.User, error) {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_INSERT"])
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	hashPass := auth.HashAndSalt(register.Password)
	err = stmt.QueryRow(register.Fullname, register.Email, hashPass, register.Phone).Scan(&user.ID, &user.Fullname, &user.Email, &user.Phone)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_EXISTS_WITH_EMAIL"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(email)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return false, err
		}
	}
	return exists > 0, nil
}

func (r *UserRepository) IDExists(ctx context.Context, id int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_EXISTS_WITH_ID"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(id)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return false, err
		}
	}
	return exists > 0, nil
}

func (r *UserRepository) GetWithId(ctx context.Context, id int) (*model.User, error) {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_GET_WITH_ID"])
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	found := false
	user := &model.User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Fullname, &user.Email, &user.IsAdmin, &user.Password); err != nil {
			return nil, err
		}
		found = true
	}

	if !found {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *UserRepository) GetWithMail(ctx context.Context, email string) (*model.User, error) {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_GET_WITH_EMAIL"])
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(email)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	found := false
	user := &model.User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Fullname, &user.Email, &user.IsAdmin, &user.Password); err != nil {
			return nil, err
		}
		found = true
	}

	if !found {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *UserRepository) ProfileUpdate(ctx context.Context, query string, params []any) error {

	stmt, err := config.App().DB.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(params...)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("user not updated")
	}

	return nil
}

func (r *UserRepository) PasswordUpdate(ctx context.Context, password string, userId int) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_UPDATE_PASS"])
	if err != nil {
		return err
	}

	updateAt := time.Now().Format("2006-01-02 15:04:05")
	hashPass := auth.HashAndSalt(password)
	result, err := stmt.Exec(hashPass, updateAt, userId)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("user password not updated")
	}

	return nil
}

func (r *UserRepository) LastLoginUpdate(ctx context.Context, userId int) error {
	lastLogin := time.Now().Format("2006-01-02 15:04:05")

	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_LAST_LOGIN"])
	if err != nil {
		return err
	}

	result, err := stmt.Exec(lastLogin, userId)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("user last login not updated")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userID int) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["USER_DELETE"])
	if err != nil {
		return err
	}

	deleteAndUpdate := time.Now().Format("2006-01-02 15:04:05")

	result, err := stmt.Exec(false, deleteAndUpdate, deleteAndUpdate, userID)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("user not deleted")
	}

	return nil
}
