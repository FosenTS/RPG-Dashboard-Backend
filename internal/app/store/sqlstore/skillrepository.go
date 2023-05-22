package sqlstore

import (
	"home/fosen/Document/golang/RestAPI/internal/model"
)

type SkillRepository struct {
	store *Store
}

func (r *SkillRepository) Create(ms *model.Skill) error {
	return r.store.db.QueryRow(
		"INSERT INTO skills (user_email, group_skills, description) values($1, $2, $3) returning user_email",
		ms.User_email,
		ms.Group_skills,
		ms.Description,
	).Scan(&ms.User_email, &ms.Group_skills, &ms.Description)
}

func (r *SkillRepository) GetAllSkills() ([]model.Skill, error) {
	var array_s []model.Skill
	rows, err := r.store.db.Query(
		"SELECT user_email, group_skills, description FROM skills")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := &model.Skill{}
		if err := rows.Scan(
			&s.User_email,
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
		"SELECT user_email, group_skills, description FROM skills WHERE user_email = $1",
		user_email)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		s := &model.Skill{}
		if err := rows.Scan(
			&s.User_email,
			&s.Group_skills,
			&s.Description,
		); err != nil {
			return nil, err
		}
		array_s = append(array_s, *s)
	}
	return array_s, nil
}

func (r *SkillRepository) FindByEmail_Gs(user_email string, group_skills string) ([]model.Skill, error) {
	var array_s []model.Skill
	rows, err := r.store.db.Query(
		"SELECT user_email, group_skills, description FROM skills WHERE user_email = $1 AND WHERE group_skills = $2",
		user_email, group_skills)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		s := &model.Skill{}
		if err := rows.Scan(
			&s.User_email,
			&s.Group_skills,
			&s.Description,
		); err != nil {
			return nil, err
		}
		array_s = append(array_s, *s)
	}
	return array_s, nil
}
