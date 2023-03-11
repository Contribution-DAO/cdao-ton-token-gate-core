package controllers

import (
	"github.com/Contribution-DAO/cdao-ton-token-gate-core/services"
	"gorm.io/gorm"
)

type ControllerHandler struct {
	s  *services.ServiceHandler
	db *gorm.DB
}

func NewControllerHandler(s *services.ServiceHandler, db *gorm.DB) *ControllerHandler {
	return &ControllerHandler{s, db}
}
