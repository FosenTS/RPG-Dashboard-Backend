package apiserver

import (
	"encoding/json"
	"errors"
	"home/fosen/Document/golang/RestAPI/internal/app/store"
	"home/fosen/Document/golang/RestAPI/internal/model"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionName = "RestAPI_Auth"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

type server struct {
	router       *mux.Router
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users/create", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/tasks/create", s.handlerTaskCreate()).Methods("POST")
	s.router.HandleFunc("/tasks", s.handleTaskGetUser()).Methods("POST")
	s.router.HandleFunc("/tasks/complete", s.handleTaskComplete()).Methods("POST")
	s.router.HandleFunc("/users/get", s.handleGetUser()).Methods("POST")
	s.router.HandleFunc("/users", s.handleGetAllUsers()).Methods("POST")
	s.router.HandleFunc("/skills/add", s.handleAddSkill()).Methods("POST")
	s.router.HandleFunc("/skills/get", s.handleGetSkill()).Methods("POST")
	s.router.HandleFunc("/skills/get/email_gs", s.handleGetSkillEmail_Gs()).Methods("POST")
}

func (s *server) handleUsersCreate() http.HandlerFunc {

	type request struct {
		Email    string `json:"email"`
		UserName string `json:"user_name"`
		Password string `json:"password"`
		Role     string `json:"user_role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			UserName: req.UserName,
			Password: req.Password,
			UserRole: req.Role,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByMail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handlerTaskCreate() http.HandlerFunc {

	type request struct {
		Name_curator   string `json:"name_curator"`
		Email_curator  string `json:"email_curator"`
		Email_employee string `json:"email_employee"`
		Description    string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByMail(req.Email_curator)
		if err != nil || u.UserRole != "curator" {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		t := &model.Task{
			Name_curator:   req.Name_curator,
			Email_curator:  req.Email_curator,
			Email_employee: req.Email_employee,
			Description:    req.Description,
		}

		if err := s.store.Task().Create(t); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, t)

	}
}

func (s *server) handleTaskGetUser() http.HandlerFunc {

	type request struct {
		Name_curator string `json:"name_curator"`
		Email        string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		array_t, err := s.store.Task().GetUserTask(req.Email)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, array_t)
	}

}

func (s *server) handleGetUser() http.HandlerFunc {

	type request struct {
		Id int `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindById(req.Id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusAccepted, u)
	}
}

func (s *server) handleGetAllUsers() http.HandlerFunc {

	type request struct {
		Role string `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.Role != "" {
			list_u, err := s.store.User().GetAllUser_filter(req.Role)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			s.respond(w, r, http.StatusAccepted, list_u)
			return
		}

		list_u, err := s.store.User().GetAllUser()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusAccepted, list_u)
	}
}

func (s *server) handleTaskComplete() http.HandlerFunc {

	type request struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err := s.store.Task().StatusUpdate(req.Email)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusAccepted, "Updated")
	}
}
func (s *server) handleAddSkill() http.HandlerFunc {

	type request struct {
		Group_skills string `json:"Group_skills"`
		Description  string `json:"Description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		ms := &model.Skill{
			Group_skills: req.Group_skills,
			Description:  req.Description,
		}

		if req.Group_skills == "" || req.Description == "" {
			w.WriteHeader(http.StatusBadRequest)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Skill().Create(ms); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, ms)

	}
}

func (s *server) handleGetSkill() http.HandlerFunc {

	type request struct {
		User_email string `json:"User_email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.User_email != "" {
			list_s, err := s.store.Skill().FindByEmail(req.User_email)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			s.respond(w, r, http.StatusAccepted, list_s)
			return
		}

		list_s, err := s.store.Skill().GetAllSkills()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusAccepted, list_s)
	}
}

func (s *server) handleGetSkillEmail_Gs() http.HandlerFunc {

	type request struct {
		User_email   string `json:"User_email"`
		Group_skills string `json:"Group_skills"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.User_email != "" && req.Group_skills != "" {
			list_s, err := s.store.Skill().FindByEmail_Gs(req.User_email, req.Group_skills)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			s.respond(w, r, http.StatusAccepted, list_s)
			return
		}

		list_s, err := s.store.Skill().GetAllSkills()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		type skill_ext struct {
			Group_skills string
			array        []model.Skill
			num          int
		}

		var sk skill_ext = skill_ext{req.User_email, list_s, len(list_s)}

		s.respond(w, r, http.StatusAccepted, sk)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
