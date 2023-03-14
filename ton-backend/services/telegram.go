package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Contribution-DAO/cdao-ton-token-gate-core/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

// Start link telegram group

// id, first_name, last_name, username, photo_url, auth_date and hash
func CheckTelegramAuthorization(params map[string][]string) bool {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	keyHash := sha256.New()
	keyHash.Write([]byte(token))
	secretkey := keyHash.Sum(nil)

	var checkparams []string
	for k, v := range params {
		if k != "hash" {
			checkparams = append(checkparams, fmt.Sprintf("%s=%s", k, v[0]))
		}
	}
	sort.Strings(checkparams)
	checkString := strings.Join(checkparams, "\n")
	hash := hmac.New(sha256.New, secretkey)
	hash.Write([]byte(checkString))
	hashstr := hex.EncodeToString(hash.Sum(nil))
	fmt.Println(hashstr)
	if hashstr == params["hash"][0] {
		return true
	}
	return false
}

func (h *ServiceHandler) LinkTelegram(address string, telegramUserId string, telegramName string, telegramUsername string, telegramAvatar string) (*model.Wallet, error) {
	wallet := model.Wallet{
		ID: address,
	}

	if result := h.db.Model(&wallet).Updates(map[string]interface{}{
		"telegram_user_id":  telegramUserId,
		"telegram_name":     telegramName,
		"telegram_username": telegramUsername,
		"telegram_avatar":   telegramAvatar,
	}); result.Error != nil {
		return nil, result.Error
	}

	return &wallet, nil
}

// End link telegram group

func (h *ServiceHandler) ConvertToTelegramGroupResponseDTO(group *model.TelegramGroup, address string) *model.TelegramGroupResponseDTO {
	if address == "" {
		return &model.TelegramGroupResponseDTO{
			TelegramGroup: *group,
			IsOwner:       false,
			IsApproved:    false,
			IsJoined:      false,
		}
	} else {
		var sbt *model.Sbt

		for _, x := range group.Sbts {
			if x.WalletID == address {
				sbt = &x
				break
			}
		}

		return &model.TelegramGroupResponseDTO{
			TelegramGroup: *group,
			IsOwner:       group.WalletID == address,
			IsApproved:    sbt != nil,
			IsJoined:      sbt != nil && sbt.IsJoined,
		}
	}
}

func (h *ServiceHandler) ListTelegramGroups(address string) ([]model.TelegramGroupResponseDTO, error) {
	var groups []model.TelegramGroup

	if result := h.db.Preload("Wallet").Preload("Sbts").Find(&groups); result.Error != nil {
		return nil, result.Error
	}

	var responseDto []model.TelegramGroupResponseDTO

	for _, group := range groups {
		groupDto := h.ConvertToTelegramGroupResponseDTO(&group, address)
		if !groupDto.IsSecret || groupDto.IsApproved {
			responseDto = append(responseDto, *groupDto)
		}
	}

	return responseDto, nil
}

func (h *ServiceHandler) ListOwnedTelegramGroups(address string) ([]model.TelegramGroup, error) {
	var groups []model.TelegramGroup

	if result := h.db.Preload("Wallet").Where("wallet_id = ?", address).Find(&groups); result.Error != nil {
		return nil, result.Error
	}

	return groups, nil
}

func (h *ServiceHandler) GetTelegramGroupSimple(id string) (*model.TelegramGroup, error) {
	var group model.TelegramGroup

	if result := h.db.First(&group, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	return &group, nil
}

func (h *ServiceHandler) GetTelegramGroup(id string, address string) (*model.TelegramGroupResponseDTO, error) {
	var group model.TelegramGroup

	if result := h.db.Preload("Sbts").Preload("Wallet").First(&group, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	return h.ConvertToTelegramGroupResponseDTO(&group, address), nil
}

func (h *ServiceHandler) CreateTelegramGroup(id string, address string, twitterUsername string, isSecret bool) (*model.TelegramGroup, error) {
	twitterUsername = strings.TrimPrefix(twitterUsername, "@")

	// Replace with your bot token
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return nil, err
	}

	// Replace with the chat ID you want to get information for
	chatID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	// Get chat information
	chatConfig := tgbotapi.ChatConfig{ChatID: chatID}
	chat, err := bot.GetChat(chatConfig)
	if err != nil {
		return nil, err
	}

	// Get chat avatar URL
	var avatarURL string
	if chat.Photo != nil {
		photoConfig := tgbotapi.FileConfig{FileID: chat.Photo.BigFileID}
		photo, err := bot.GetFile(photoConfig)
		if err != nil {
			return nil, err
		}
		avatarURL = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, photo.FilePath)
	}

	// Print chat information
	// fmt.Println("Chat ID:", chat.ID)
	// fmt.Println("Chat Type:", chat.Type)
	// fmt.Println("Chat Title:", chat.Title)
	// fmt.Println("Chat Username:", chat.UserName)
	// fmt.Println("Chat First Name:", chat.FirstName)
	// fmt.Println("Chat Last Name:", chat.LastName)
	// fmt.Println("Chat Avatar URL:", avatarURL)

	// Get existing group
	existingGroup, err := h.GetTelegramGroup(id, address)

	group := &model.TelegramGroup{
		ID:              id,
		WalletID:        address,
		TwitterUsername: twitterUsername,
		Name:            chat.Title,
		Avatar:          avatarURL,
		IsSecret:        isSecret,
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a new group
			if result := h.db.Create(group); result.Error != nil {
				return nil, result.Error
			}

			return group, nil
		} else {
			return nil, err
		}
	}

	if existingGroup.IsOwner {
		// Update group
		if result := h.db.Model(&group).Updates(group); result.Error != nil {
			return nil, result.Error
		}

		return group, nil
	} else {
		return nil, errors.New("not owner")
	}
}
