package entity

type Account struct {
	ID       string `gorm:"type:varchar(255);primaryKey" json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
	IsDelete bool   `json:"is_delete"`
	StatusID string `json:"status_id"`
	RoleID   string `json:"role_id"`
	UserID   string `json:"user_id"`
}
