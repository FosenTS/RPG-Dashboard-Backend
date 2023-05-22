package sqlstore

import (
	"database/sql"
	"fmt"
	"home/fosen/Document/golang/RestAPI/internal/app/store"
	"home/fosen/Document/golang/RestAPI/internal/model"
)

type TaskRepository struct {
	store *Store
}

func (r *TaskRepository) Create(t *model.Task) error {
	t.Status = "false"
	return r.store.db.QueryRow(
		"INSERT INTO tasks(name_curator, email_curator, email_employee, description) values($1, $2, $3, $4) returning id",
		t.Name_curator,
		t.Email_curator,
		t.Email_employee,
		t.Description,
	).Scan(&t.ID)
}

func (r *TaskRepository) GetUserTask(email string) ([]model.Task, error) {
	var array_t []model.Task
	rows, err := r.store.db.Query(
		"Select id, name_curator, email_curator, email_employee, description, status from tasks where email_employee = $1",
		email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.Task{}
		if err := rows.Scan(
			&t.ID,
			&t.Name_curator,
			&t.Email_curator,
			&t.Email_employee,
			&t.Description,
			&t.Status,
		); err != nil {
			return nil, err
		}
		array_t = append(array_t, *t)
	}
	return array_t, nil
}

func (r *TaskRepository) SearchReward(id int) (*int, error) {
	t := &model.Task{}
	if err := r.store.db.QueryRow(
		"SELECT reward FROM tasks WHERE id = $1",
		id,
	).Scan(
		&t.Reward,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return &t.Reward, nil
}

func (r *TaskRepository) StatusUpdate(email string) error {

	if _, err := r.store.db.Exec("UPDATE tasks set status = $1 WHERE email_employee = $2",
		true, email,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		fmt.Println(err)
		return err
	}
	return nil
}
