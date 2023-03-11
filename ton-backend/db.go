package main

import (
	"os"

	"github.com/Contribution-DAO/cdao-ton-token-gate-core/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb() *gorm.DB {
	db, dbErr := gorm.Open(postgres.Open(os.Getenv("CONNECTION_STRING")), &gorm.Config{})

	if dbErr != nil {
		panic(dbErr)
	}

	// db.AutoMigrate(&model.User{})
	// db.AutoMigrate(&model.Project{})
	// db.AutoMigrate(&model.Deployment{})
	// db.AutoMigrate(&model.Subgraph{})
	// db.AutoMigrate(&model.SubgraphUser{})
	// db.AutoMigrate(&model.Node{})
	// db.AutoMigrate(&model.NodeUser{})

	// db.Exec("CREATE INDEX IF NOT EXISTS \"idx_nodes_required_nodes\" ON \"nodes\" USING GIN (\"required_nodes\")")

	return db
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.TelegramGroup{})
	db.AutoMigrate(&model.Wallet{})
	db.AutoMigrate(&model.Sbt{})
	db.AutoMigrate(&model.TelegramApproval{})
}
