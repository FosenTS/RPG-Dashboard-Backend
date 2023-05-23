package model

type Skill struct {
	ID           int    `json:"id"`
	Email        string `json:"user_email"`
	Group_skills string `json:"group_skills"`
	Description  string `json:"description"`
}
