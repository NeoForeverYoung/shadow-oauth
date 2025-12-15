package models

import (
	"time"

	"gorm.io/gorm"
)

// AuthorizationCode 授权码模型
// 临时授权码，用于交换 Access Token
type AuthorizationCode struct {
	ID            uint           `gorm:"primarykey" json:"id"`                    // 主键
	Code          string         `gorm:"uniqueIndex;not null;size:100" json:"code"` // 授权码（唯一）
	ClientID      string         `gorm:"not null;size:100" json:"client_id"`      // 客户端ID
	UserID        uint           `gorm:"not null" json:"user_id"`                 // 用户ID
	RedirectURI   string         `gorm:"not null" json:"redirect_uri"`            // 重定向URI
	ExpiresAt     time.Time      `gorm:"not null" json:"expires_at"`              // 过期时间（通常10分钟）
	Used          bool           `gorm:"default:false" json:"used"`               // 是否已使用（授权码只能使用一次）
	CreatedAt     time.Time      `json:"created_at"`                               // 创建时间
	UpdatedAt     time.Time      `json:"updated_at"`                               // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`                         // 软删除时间
}

// TableName 指定表名
func (AuthorizationCode) TableName() string {
	return "authorization_codes"
}

// IsExpired 检查授权码是否过期
func (ac *AuthorizationCode) IsExpired() bool {
	return time.Now().After(ac.ExpiresAt)
}

// IsValid 检查授权码是否有效（未使用且未过期）
func (ac *AuthorizationCode) IsValid() bool {
	return !ac.Used && !ac.IsExpired()
}

