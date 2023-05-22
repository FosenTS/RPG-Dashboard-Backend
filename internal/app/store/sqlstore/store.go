package sqlstore

import (
	"database/sql"
	"home/fosen/Document/golang/RestAPI/internal/app/store"

	_ "github.com/lib/pq"
)

type Store struct {
	db              *sql.DB
	UserRepository  *UserRepository
	TaskRepository  *TaskRepository
	SkillRepository *SkillRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
	}

	return s.UserRepository
}

func (s *Store) Task() store.TaskRepository {
	if s.TaskRepository != nil {
		return s.TaskRepository
	}

	s.TaskRepository = &TaskRepository{
		store: s,
	}

	return s.TaskRepository
}

func (s *Store) Skill() store.SkillRepository {
	if s.SkillRepository != nil {
		return s.SkillRepository
	}

	s.SkillRepository = &SkillRepository{
		store: s,
	}

	return s.SkillRepository
}
