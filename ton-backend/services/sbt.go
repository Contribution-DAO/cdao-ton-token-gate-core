package services

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/Contribution-DAO/cdao-ton-token-gate-core/model"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/nft"
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

func (h *ServiceHandler) ValidateNft(nftContractAddress string, nftOwner string) (bool, string, error) {
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		return false, "", err
	}

	api := ton.NewAPIClient(client)

	nftAddr := address.MustParseAddr(nftContractAddress)
	item := nft.NewItemClient(api, nftAddr)

	nftData, err := item.GetNFTData(context.Background())
	if err != nil {
		return false, "", err
	}

	if !nftData.Initialized {
		return false, "", err
	}

	collectionAddress := nftData.CollectionAddress.String()

	fmt.Println(collectionAddress)

	ownerAddress := nftData.OwnerAddress.String()

	fmt.Println(ownerAddress)

	return collectionAddress == os.Getenv("TON_COLLECTION_ADDRESS") && ownerAddress == nftOwner, nftData.Content.(*nft.ContentOffchain).URI, nil
}

func (h *ServiceHandler) ScanSbt(approvalId string, nftOwner string) (string, error) {
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		return "", err
	}

	api := ton.NewAPIClient(client)

	// get info about our nft's collection
	collection := nft.NewCollectionClient(api, address.MustParseAddr(os.Getenv("TON_COLLECTION_ADDRESS")))
	collectionData, err := collection.GetCollectionData(context.Background())
	if err != nil {
		return "", err
	}

	nextIndex := collectionData.NextItemIndex.Int64()
	finalIndex := nextIndex - 10

	if finalIndex < 0 {
		finalIndex = 0
	}

	for i := nextIndex; i >= finalIndex; i-- {
		nftAddress, err := collection.GetNFTAddressByIndex(context.Background(), big.NewInt(i))
		if err != nil {
			continue
		}

		valid, nftApprovalId, err := h.ValidateNft(nftAddress.String(), nftOwner)

		fmt.Println(nftApprovalId, approvalId, nftApprovalId == approvalId)

		if !valid || err != nil {
			continue
		}

		if nftApprovalId == approvalId {
			return nftAddress.String(), nil
		}
	}

	return "", nil
}
