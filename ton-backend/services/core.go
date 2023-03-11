package services

import (
	"gorm.io/gorm"
)

type ServiceHandler struct {
	db *gorm.DB
}

func NewServiceHandler(db *gorm.DB) *ServiceHandler {
	return &ServiceHandler{db}
}

func TxServiceHandler(h ServiceHandler, tx *gorm.DB) *ServiceHandler {
	h.db = tx
	return &h
}

func (h *ServiceHandler) Transaction(fc func(tx *ServiceHandler) error) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		return fc(TxServiceHandler(*h, tx))
	})
}
