package model

type Task struct {
	ID             int    `json:"id"`
	Name_curator   string `json:"name_curator"`
	Email_curator  string `json:"email_curator"`
	Email_employee string `json:"email_employee"`
	Description    string `json:"description"`
	Status         string `json:"false"`
}
