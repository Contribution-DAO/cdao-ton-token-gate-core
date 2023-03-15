package controllers

import (
	"net/http"

	"github.com/Contribution-DAO/cdao-ton-token-gate-core/model"
	"github.com/Contribution-DAO/cdao-ton-token-gate-core/services"
	"github.com/gin-gonic/gin"
)

func (h *ControllerHandler) HandleTelegramCallback(c *gin.Context) {
	address, err := GetUidFromHeader(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	params := c.Request.URL.Query()
	ok := services.CheckTelegramAuthorization(params)
	if !ok {
		c.String(http.StatusUnauthorized, "Invalid authorization")
		return
	}

	user, err := h.s.LinkTelegram(address, params["id"][0], params["first_name"][0], params["username"][0], params["photo_url"][0])

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func (h *ControllerHandler) ListTelegramGroups(c *gin.Context) {
	address, err := GetUidFromHeaderOptional(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	groups, err := h.s.ListTelegramGroups(address)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, groups)
	}
}

func (h *ControllerHandler) GetTelegramGroup(c *gin.Context) {
	groupId := c.Param("id")

	address, err := GetUidFromHeaderOptional(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	group, err := h.s.GetTelegramGroup(groupId, address)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, group)
	}
}

func (h *ControllerHandler) GetTelegramGroupRoot(c *gin.Context) {
	groupId := c.Param("id")
	address := c.Param("address")

	group, err := h.s.GetTelegramGroup(groupId, address)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, group)
	}
}

func (h *ControllerHandler) CreateTelegramGroup(c *gin.Context) {
	address, err := GetUidFromHeader(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var dto model.CreateTelegramGroupDTO

	if err := c.Bind(&dto); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	group, err := h.s.CreateTelegramGroup(dto.Id, address, dto.TwitterUsername, dto.InvitationLink, dto.IsSecret)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, group)
	}
}
