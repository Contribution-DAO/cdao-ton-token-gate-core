package services

import (
	"github.com/Contribution-DAO/cdao-ton-token-gate-core/model"
	"gorm.io/gorm/clause"
)

func (h *ServiceHandler) CreateSbt(approvalId string, address string, telegramGroupId string, contractAddress string) (*model.Sbt, error) {
	sbt := &model.Sbt{
		ID:              approvalId,
		WalletID:        address,
		TelegramGroupID: telegramGroupId,
		ContractAddress: contractAddress,
		IsJoined:        false,
	}

	if result := h.db.Clauses(clause.OnConflict{DoNothing: true}).Create(sbt); result.Error != nil {
		return nil, result.Error
	}

	return sbt, nil
}

func (h *ServiceHandler) MarkJoined(address string, telegramGroupId string) error {
	if result := h.db.Model(&model.Sbt{}).Where("wallet_id = ? and telegram_group_id = ?", address, telegramGroupId).Update("is_joined", true); result.Error != nil {
		return result.Error
	}

	return nil
}
