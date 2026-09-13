package dto

type RegisterRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponseDTO struct {
	Success bool   `json:"success"`
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
