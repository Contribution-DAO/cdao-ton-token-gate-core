package controllers

import (
	"net/http"
	"strings"

	"github.com/Contribution-DAO/cdao-ton-token-gate-core/services"
	"github.com/gin-gonic/gin"
)

func GetUidFromHeaderOptional(c *gin.Context) (string, error) {
	s := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(s, "Bearer ")

	if token == "" {
		return "", nil
	}

	uid, err := services.DecodeJwtToken(token)

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return "", err
	}

	return uid, nil
}

func GetUidFromHeader(c *gin.Context) (string, error) {
	s := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(s, "Bearer ")

	uid, err := services.DecodeJwtToken(token)

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return "", err
	}

	return uid, nil
}

func AuthorizationMiddleware(c *gin.Context) {
	_, _ = GetUidFromHeader(c)
}

func (h *ControllerHandler) MeHandler(c *gin.Context) {
	uid, err := GetUidFromHeader(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, err := h.s.GetWallet(uid)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
