package controllers

import (
	"net/http"

	"github.com/Contribution-DAO/cdao-ton-token-gate-core/model"
	"github.com/gin-gonic/gin"
)

func (h *ControllerHandler) LinkSbt(c *gin.Context) {
	address, err := GetUidFromHeader(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var sbt model.LinkSbtDTO

	if err := c.Bind(&sbt); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	wallet, err := h.s.GetWalletSimple(address)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	validated, validatedApprovalId, err := h.s.ValidateNft(sbt.ContractAddress, wallet.HumanAddress)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if validated {
		if validatedApprovalId != sbt.ApprovalId {
			c.JSON(http.StatusBadRequest, gin.H{
				"validated": false,
			})
			return
		}

		approval, err := h.s.GetTelegramApproval(sbt.ApprovalId)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		sbtOut, err := h.s.CreateSbt(sbt.ApprovalId, address, approval.TelegramGroupID, sbt.ContractAddress)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, sbtOut)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"validated": false,
		})
	}
}

func (h *ControllerHandler) ScanSbt(c *gin.Context) {
	approvalId := c.Param("approvalId")
	address, err := GetUidFromHeader(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	contractAddress, err := h.s.ScanSbt(approvalId, address)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if contractAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"validated": false,
		})
		return
	}

	approval, err := h.s.GetTelegramApproval(approvalId)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	sbtOut, err := h.s.CreateSbt(approvalId, address, approval.TelegramGroupID, contractAddress)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, sbtOut)
}

func (h *ControllerHandler) GetNftMetadata(c *gin.Context) {
	approvalId := c.Param("approvalId")

	approval, err := h.s.GetTelegramApproval(approvalId)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	group, err := h.s.GetTelegramGroup(approval.TelegramGroupID, "")

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	image := group.Avatar

	if image == "" {
		image = "https://ton-connect.contributiondao.com/connect/img/telegram-whitebg.png"
	}

	c.JSON(http.StatusOK, gin.H{
		"name":        group.Name,
		"description": "You have been approved to join \"" + group.Name + "\" telegram group through CDAO SBT Token Gate",
		"image":       image,
		"attributes": []gin.H{
			{
				"trait_type": "Group",
				"value":      group.Name,
			},
		},
	})
}
