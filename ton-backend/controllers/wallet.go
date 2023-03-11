package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Contribution-DAO/cdao-ton-token-gate-core/model"
	"github.com/gin-gonic/gin"
)

func (h *ControllerHandler) GenerateWalletSignPayload(c *gin.Context) {
	payload, err := h.s.GenerateWalletSignPayload()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"payload": payload,
		})
	}
}

func (h *ControllerHandler) ValidateWalletSignature(c *gin.Context) {
	var proof model.TonWalletProofDTO

	if err := c.Bind(&proof); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	payloadBytes, _ := json.Marshal(proof)
	println(string(payloadBytes))

	token, err := h.s.ValidateWalletSignature(proof)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	account, err := h.s.GetTonAddressInfo(token, proof.Network)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	wallet, err := h.s.CreateWallet(proof.Address, account.Address.Bounceable)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"wallet": wallet,
		})
	}
}

func (h *ControllerHandler) GetTonAddressInfo(c *gin.Context) {
	network := c.Query("network")
	s := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(s, "Bearer ")

	info, err := h.s.GetTonAddressInfo(token, network)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, info)
	}
}

func (h *ControllerHandler) GetWallet(c *gin.Context) {
	walletId := c.Param("id")

	wallet, err := h.s.GetWallet(walletId)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, wallet)
	}
}
