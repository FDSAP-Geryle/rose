package models

import "time"

type User struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName            string    `gorm:"unique;not null" json:"username"`
	Password            string    `json:"password"`
	Email               string    `json:"email"`
	RoleId              uint      `json:"roleid"`
	FullName            string    `json:"fullname"`
	MobileNo            string    `json:"mobilenumber"`
	PwdExpiredDate      time.Time `json:"password_expired_at"`
	LastLoginDate       time.Time `json:"last_login_date"`
	LastLoginFrom       string    `json:"last_login_from"`
	LastLogoutDate      time.Time `json:"last_logout_date"`
	Logged              int
	IsLock              int
	NumberOfFailedLogin int
}

func (User) TableName() string {
	return "app_user"
}

type PasswordHistory struct {
	ID        uint
	UserId    uint
	Password  string
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (PasswordHistory) TableName() string {
	return "user_password_history"
}
