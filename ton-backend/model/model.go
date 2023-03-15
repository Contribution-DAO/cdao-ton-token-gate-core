package model

import (
	"time"
)

type TelegramGroup struct {
	ID              string `json:"id" gorm:"type:text;primarykey"`
	WalletID        string `json:"walletAddress" gorm:"index"`
	TwitterUsername string `json:"twitterUsername" gorm:"index"`
	Name            string `json:"name" gorm:"index"`
	Avatar          string `json:"avatar"`
	InvitationLink  string `json:"invitationLink"`
	IsSecret        bool   `json:"isSecret"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Sbts              []Sbt              `json:"sbts"`
	Wallet            Wallet             `json:"wallet"`
	TelegramApprovals []TelegramApproval `json:"telegramApprovals"`
}

type Wallet struct {
	ID string `json:"walletAddress" gorm:"type:text;primarykey"`

	HumanAddress string `json:"humanAddress" gorm:"index"`

	TelegramUserId   string `json:"telegramUserId" gorm:"index"`
	TelegramName     string `json:"telegramName" gorm:"index"`
	TelegramUsername string `json:"telegramUsername" gorm:"index"`
	TelegramAvatar   string `json:"telegramAvatar"`

	TwitterUserId            string `json:"twitterUserId" gorm:"index"`
	TwitterUsername          string `json:"twitterUsername" gorm:"index"`
	TwitterName              string `json:"twitterName" gorm:"index"`
	TwitterAvatar            string `json:"twitterAvatar"`
	TwitterAccessToken       string `json:"twitterAccessToken"`
	TwitterAccessTokenSecret string `json:"twitterAccessTokenSecret"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Sbts              []Sbt              `json:"sbts"`
	TelegramGroups    []TelegramGroup    `json:"telegramGroups"`
	TelegramApprovals []TelegramApproval `json:"telegramApprovals"`
}

type Sbt struct {
	ID              string `json:"id" gorm:"type:text;primarykey"`
	WalletID        string `json:"walletAddress" gorm:"index"`
	TelegramGroupID string `json:"telegramGroupId" gorm:"index"`
	ContractAddress string `json:"contractAddress" gorm:"index"`
	IsJoined        bool   `json:"isJoined"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Wallet        Wallet        `json:"wallet"`
	TelegramGroup TelegramGroup `json:"telegramGroup"`
}

type TelegramApproval struct {
	ID              string `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	WalletID        string `json:"walletAddress" gorm:"index"`
	TelegramGroupID string `json:"telegramGroupId" gorm:"index"`

	Wallet        Wallet        `json:"wallet"`
	TelegramGroup TelegramGroup `json:"telegramGroup"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
