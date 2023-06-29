package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/profectus200/contact-book-bot/internal/clients/tg"
	"github.com/profectus200/contact-book-bot/internal/config"
	"github.com/profectus200/contact-book-bot/internal/database"
	"github.com/profectus200/contact-book-bot/internal/model/callbacks"
	"github.com/profectus200/contact-book-bot/internal/model/messages"
	"github.com/profectus200/contact-book-bot/internal/redis"
	"github.com/profectus200/contact-book-bot/internal/worker"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	config, err := config.New()
	if err != nil {

	}

	db, err := database.New(config)
	if err != nil {

	}

	cache, err := redis.New(config)
	if err != nil {

	}

	contactsDB := database.NewContactsDB(db, cache)
	usersDB := database.NewUsersDB(db)

	tgClient, err := tg.New(config)
	if err != nil {
	}

	msgModel := messages.New(tgClient, contactsDB, usersDB)
	callbackModel := callbacks.New(tgClient, contactsDB, usersDB)

	updateListenerWorker := worker.NewUpdateListenerWorker(tgClient, msgModel, callbackModel, cache)

	updateListenerWorker.Run(ctx)
}
