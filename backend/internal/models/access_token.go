package models

import (
	"time"

	"gorm.io/gorm"
)

// AccessToken OAuth Access Token 模型
// 用于访问受保护资源的令牌
type AccessToken struct {
	ID          uint           `gorm:"primarykey" json:"id"`                    // 主键
	Token       string         `gorm:"uniqueIndex;not null;size:500" json:"token"` // Token 字符串（JWT）
	ClientID    string         `gorm:"not null;size:100" json:"client_id"`      // 客户端ID
	UserID      uint           `gorm:"not null" json:"user_id"`                // 用户ID
	ExpiresAt   time.Time      `gorm:"not null" json:"expires_at"`               // 过期时间
	CreatedAt   time.Time      `json:"created_at"`                               // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                              // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                          // 软删除时间
}

// TableName 指定表名
func (AccessToken) TableName() string {
	return "access_tokens"
}

// IsExpired 检查 Token 是否过期
func (at *AccessToken) IsExpired() bool {
	return time.Now().After(at.ExpiresAt)
}

