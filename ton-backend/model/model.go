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
	IsSecret        bool   `json:"isSecret"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Sbts   []Sbt  `json:"sbts"`
	Wallet Wallet `json:"wallet"`
}

type Wallet struct {
	ID string `json:"walletAddress" gorm:"type:text;primarykey"`

	HumanAddress string `json:"humanAddress" gorm:"index"`

	TelegramUserId   string `json:"telegramUserId" gorm:"index"`
	TelegramName     string `json:"telegramName" gorm:"index"`
	TelegramUsername string `json:"telegramUsername" gorm:"index"`
	TelegramAvatar   string `json:"telegramAvatar"`

	TwitterUsername          string `json:"twitterUsername" gorm:"index"`
	TwitterAvatar            string `json:"twitterAvatar"`
	TwitterAccessToken       string `json:"twitterAccessToken"`
	TwitterAccessTokenSecret string `json:"twitterAccessTokenSecret"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Sbts           []Sbt           `json:"sbts"`
	TelegramGroups []TelegramGroup `json:"telegramGroups"`
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

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
