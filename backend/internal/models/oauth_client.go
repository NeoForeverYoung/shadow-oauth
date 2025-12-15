package models

import (
	"time"

	"gorm.io/gorm"
)

// OAuthClient OAuth 客户端模型
// 代表一个第三方应用，想要访问用户资源
type OAuthClient struct {
	ID          uint           `gorm:"primarykey" json:"id"`                    // 客户端ID（主键）
	ClientID    string         `gorm:"uniqueIndex;not null;size:100" json:"client_id"` // 客户端标识符（公开）
	ClientSecret string        `gorm:"not null;size:255" json:"-"`              // 客户端密钥（保密，不返回JSON）
	Name        string         `gorm:"not null;size:100" json:"name"`            // 客户端名称
	RedirectURI string         `gorm:"not null" json:"redirect_uri"`             // 重定向URI（授权后跳转的地址）
	CreatedAt   time.Time      `json:"created_at"`                               // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                               // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                            // 软删除时间
}

// TableName 指定表名
func (OAuthClient) TableName() string {
	return "oauth_clients"
}

// OAuthClientResponse 客户端响应结构（不包含密钥）
type OAuthClientResponse struct {
	ID          uint      `json:"id"`
	ClientID    string    `json:"client_id"`
	Name        string    `json:"name"`
	RedirectURI string    `json:"redirect_uri"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToResponse 转换为响应结构
func (c *OAuthClient) ToResponse() OAuthClientResponse {
	return OAuthClientResponse{
		ID:          c.ID,
		ClientID:    c.ClientID,
		Name:        c.Name,
		RedirectURI: c.RedirectURI,
		CreatedAt:   c.CreatedAt,
	}
}

