package domain

type Patient struct {
	Id        int    `json:"id"`
	Surname   string `json:"surname" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Document  string `json:"document" binding:"required"`
	CreatedAt string `json:"created_at" binding:"required"`
}
