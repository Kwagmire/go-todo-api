package models

type ListCurator struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

type TodoItem struct {
	ID        int    `json:"id"`
	CuratorID int    `json:"curator_id"`
	Title     string `json:"title"`
	Desc      string `json:"description"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json: "password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json: "password"`
}

type CreateRequest struct {
	Title string `json:"title"`
	Desc  string `json:"description"`
}
