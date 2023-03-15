package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Contribution-DAO/cdao-ton-token-gate-core/model"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/nft"
	"gorm.io/gorm/clause"
)

// func (h *ServiceHandler) ListProjects(uid string) ([]model.Project, error) {
// 	var projects []model.Project

// 	if result := h.db.Where("user_id = ?", uid).Where("is_system_project IS FALSE").Preload("Deployments").Preload("Subgraphs").Find(&projects); result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return projects, nil
// }

func (h *ServiceHandler) GetWalletSimple(address string) (*model.Wallet, error) {
	var wallet model.Wallet

	field := "id"

	if !strings.Contains(address, ":") {
		// Human address
		field = "human_address"
	}

	if result := h.db.First(&wallet, field+" = ?", address); result.Error != nil {
		return nil, result.Error
	}

	return &wallet, nil
}

func (h *ServiceHandler) GetWallet(address string) (*model.Wallet, error) {
	var wallet model.Wallet

	field := "id"

	if !strings.Contains(address, ":") {
		// Human address
		field = "human_address"
	}

	if result := h.db.Preload("Sbts").Preload("TelegramGroups").First(&wallet, field+" = ?", address); result.Error != nil {
		return nil, result.Error
	}

	return &wallet, nil
}

func (h *ServiceHandler) CreateWallet(address string, humanAddress string) (*model.Wallet, error) {
	wallet := &model.Wallet{
		ID:           address,
		HumanAddress: humanAddress,
	}

	if result := h.db.Clauses(clause.OnConflict{DoNothing: true}).Create(wallet); result.Error != nil {
		return nil, result.Error
	}

	return wallet, nil
}

func (h *ServiceHandler) GenerateWalletSignPayload() (string, error) {
	// create the request
	req, err := http.NewRequest("POST", os.Getenv("TON_SIGNATURE_HOST")+"/ton-proof/generatePayload", nil)
	if err != nil {
		return "", err
	}

	// set the content type header
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// decode the response body
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response["payload"].(string), nil
}

// Return token
func (h *ServiceHandler) ValidateWalletSignature(proof model.TonWalletProofDTO) (string, error) {
	// encode the payload as JSON
	payloadBytes, err := json.Marshal(proof)
	if err != nil {
		return "", err
	}

	// create the request
	req, err := http.NewRequest("POST", os.Getenv("TON_SIGNATURE_HOST")+"/ton-proof/checkProof", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	// set the content type header
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// decode the response body
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	xxx, _ := json.Marshal(response)

	println(string(payloadBytes))
	println(string(xxx))

	return response["token"].(string), nil
}

func (h *ServiceHandler) GetTonAddressInfo(token string, network string) (*model.TonAccountInfo, error) {
	endpoint := os.Getenv("TON_SIGNATURE_HOST") + "/dapp/getAccountInfo"
	params := "?network=" + network

	req, err := http.NewRequest("GET", endpoint+params, nil)
	if err != nil {
		return nil, err
	}

	// Add Authorization Bearer header
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data model.TonAccountInfo
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (h *ServiceHandler) ValidateNft(nftContractAddress string, nftOwner string) (bool, error) {
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		return false, err
	}

	api := ton.NewAPIClient(client)

	nftAddr := address.MustParseAddr(nftContractAddress)
	item := nft.NewItemClient(api, nftAddr)

	nftData, err := item.GetNFTData(context.Background())
	if err != nil {
		return false, err
	}

	if !nftData.Initialized {
		return false, err
	}

	collectionAddress := nftData.CollectionAddress.String()

	fmt.Println(collectionAddress)

	ownerAddress := nftData.OwnerAddress.String()

	fmt.Println(ownerAddress)

	return collectionAddress == os.Getenv("TON_COLLECTION_ADDRESS") && ownerAddress == nftOwner, nil
}
