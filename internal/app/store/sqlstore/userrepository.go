package sqlstore

import (
	"database/sql"
	"home/fosen/Document/golang/RestAPI/internal/app/store"
	"home/fosen/Document/golang/RestAPI/internal/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users(email, name_user, encrypted_password, role) VALUES($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.UserName,
		u.EncryptedPassword,
		u.UserRole,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByMail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, name_user, encrypted_password, role, user_level FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.UserName,
		&u.EncryptedPassword,
		&u.UserRole,
		&u.UserLevel,
	); err != nil {

		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindById(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, name_user, role, user_level FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.UserName,
		&u.UserRole,
		&u.UserLevel,
	); err != nil {

		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetAllUser() ([]model.User, error) {
	var array_u []model.User
	rows, err := r.store.db.Query(
		"Select id, email, name_user, role, user_level from users order by user_level desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.UserName,
			&u.UserRole,
			&u.UserLevel,
		); err != nil {
			return nil, err
		}
		array_u = append(array_u, *u)
	}
	return array_u, nil
}

func (r *UserRepository) GetAllUser_filter(role_line string) ([]model.User, error) {
	var array_u []model.User
	rows, err := r.store.db.Query(
		"Select id, email, name_user, role, user_level from users where role = $1",
		role_line)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.UserName,
			&u.UserRole,
			&u.UserLevel,
		); err != nil {
			return nil, err
		}
		array_u = append(array_u, *u)
	}
	return array_u, nil
}

func (r *UserRepository) LevelUpdate(email string, id int) error {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, name_user, encrypted_password, role, user_level FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.UserName,
		&u.EncryptedPassword,
		&u.UserRole,
		&u.UserLevel,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}

		return err
	}

	reward, err := r.store.TaskRepository.SearchReward(id)
	if err != nil {
		return err
	}

	updateLevel := *reward + u.UserLevel

	if _, err := r.store.db.Exec("UPDATE users set user_level = $1 WHERE email = $2",
		updateLevel, email,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}
	return nil
}
