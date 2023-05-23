package sqlstore

import (
	"fmt"
	"home/fosen/Document/golang/RestAPI/internal/model"
)

type SkillRepository struct {
	store *Store
}

func (r *SkillRepository) Create(ms *model.Skill) error {
	return r.store.db.QueryRow(
		"INSERT INTO skills(email, group_skills, description) values($1, $2, $3) returning id",
		ms.Email,
		ms.Group_skills,
		ms.Description,
	).Scan(&ms.ID)
}

func (r *SkillRepository) GetAllSkills() ([]model.Skill, error) {
	var array_s []model.Skill
	rows, err := r.store.db.Query(
		"SELECT email, group_skills, description FROM skills")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := &model.Skill{}
		if err := rows.Scan(
			&s.Email,
			&s.Group_skills,
			&s.Description,
		); err != nil {
			return nil, err
		}
		array_s = append(array_s, *s)
	}
	return array_s, nil
}

func (r *SkillRepository) FindByEmail(user_email string) ([]model.Skill, error) {
	var array_s []model.Skill
	rows, err := r.store.db.Query(
		"SELECT email, group_skills, description FROM skills WHERE user_email = $1",
		user_email)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		s := &model.Skill{}
		if err := rows.Scan(
			&s.Email,
			&s.Group_skills,
			&s.Description,
		); err != nil {
			return nil, err
		}
		array_s = append(array_s, *s)
	}
	return array_s, nil
}

func (r *SkillRepository) FindByEmail_gs(user_email string, group string) ([]model.Skill, error) {
	var array_s []model.Skill
	fmt.Println(user_email, group)
	rows, err := r.store.db.Query(
		"SELECT id, email, group_skills, description FROM skills WHERE email = $1 AND group_skills = $2",
		user_email, group)
	if err != nil {
		return nil, err
	}
	fmt.Println(rows)
	defer rows.Close()
	for rows.Next() {
		u := &model.Skill{}
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Group_skills,
			&u.Description,
		); err != nil {
			return nil, err
		}
		array_s = append(array_s, *u)
	}
	return array_s, nil
}
