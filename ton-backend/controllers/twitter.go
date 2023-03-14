package controllers

import (
	"context"
	"net/http"
	"os"

	"github.com/Contribution-DAO/cdao-ton-token-gate-core/services"
	"github.com/gin-gonic/gin"
	session "github.com/go-session/session/v3"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func (h *ControllerHandler) HandleUserReceived(address string, user goth.User) error {
	_, err := h.s.LinkTwitter(address, user.Name, user.AvatarURL, user.AccessToken, user.AccessTokenSecret)
	return err
}

func (h *ControllerHandler) HandleTwitterLogin(c *gin.Context) {
	store, err := session.Start(context.Background(), c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Add provider to query
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	token := c.Query("token")

	address, err := services.DecodeJwtToken(token)

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// try to get the user without re-authenticating
	if user, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		err := h.HandleUserReceived(address, user)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, user)
		}
	} else {
		// Store token in the session
		store.Set("address", address)
		err = store.Save()

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func (h *ControllerHandler) HandleTwitterCallback(c *gin.Context) {
	store, err := session.Start(context.Background(), c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Add provider to query
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	address, ok := store.Get("address")
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = h.HandleUserReceived(address.(string), user)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		// c.JSON(http.StatusOK, user)
		c.Redirect(301, os.Getenv("CONNECT_FRONTEND_HOST"))
	}
}

func (h *ControllerHandler) VerifyTwitterFollow(c *gin.Context) {
	groupId := c.Param("groupId")

	address, err := GetUidFromHeader(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	wallet, err := h.s.GetWalletSimple(address)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	group, err := h.s.GetTelegramGroup(groupId, address)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	followed, err := h.s.VerifyTwitterFollow(group.TwitterUsername, wallet.TwitterAccessToken, wallet.TwitterAccessTokenSecret)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followed": followed,
	})
}
