package models

type BanUserReq struct {
	ID    string `json:"id"`    // Username of the profile to retrieve
	Email string `json:"email"` // Username of the profile to retrieve
}

type UnbanUserReq struct {
	ID    string `json:"id"`    // Username of the profile to retrieve
	Email string `json:"email"` // Username of the profile to retrieve
}

type AddCourierReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DeleteCourierReq struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AddProductManagerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DeleteProductManagerReq struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
