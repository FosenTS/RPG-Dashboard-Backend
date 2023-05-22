package store

import (
	"home/fosen/Document/golang/RestAPI/internal/model"
)

type UserRepository interface {
	Create(*model.User) error
	FindByMail(string) (*model.User, error)
	FindById(int) (*model.User, error)
	GetAllUser() ([]model.User, error)
	GetAllUser_filter(string) ([]model.User, error)
	LevelUpdate(string, int) error
}

type TaskRepository interface {
	Create(*model.Task) error
	StatusUpdate(string) error
	GetUserTask(string) ([]model.Task, error)
	SearchReward(int) (*int, error)
}

type SkillRepository interface {
	Create(*model.Skill) error
	GetAllSkills() ([]model.Skill, error)
	FindByEmail(string) ([]model.Skill, error)
	FindByEmail_Gs(string, string) ([]model.Skill, error)
}
