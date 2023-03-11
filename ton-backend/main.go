package main

import (
	"flag"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
)

var isAutomigrate bool

func main() {
	flag.BoolVar(&isAutomigrate, "automigrate", false, "Auto Migrate")

	flag.Parse()

	if godotenv.Load() != nil {
		log.Fatalf("Error loading .env file")
		return
	}

	callbackURL := "http://127.0.0.1:8080/social-auth/twitter/callback"

	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	goth.UseProviders(
		twitter.New(os.Getenv("TWITTER_API_KEY"), os.Getenv("TWITTER_API_SECRET"), callbackURL),
		// If you'd like to use authenticate instead of authorize in Twitter provider, use this instead.
		// twitter.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),
	)

	db := InitDb()

	if isAutomigrate {
		AutoMigrate(db)
	} else {
		InitRouter(db)
	}
}
