package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`                    // 用户ID（主键）
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`       // 邮箱（唯一索引）
	Password  string         `gorm:"not null" json:"-"`                       // 密码（bcrypt加密，不返回到JSON）
	Name      string         `gorm:"size:100" json:"name"`                    // 用户名
	CreatedAt time.Time      `json:"created_at"`                              // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                              // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                          // 软删除时间
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserResponse 用户响应结构（不包含敏感信息）
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse 将 User 转换为 UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

