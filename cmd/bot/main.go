package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/profectus200/contact-book-bot/internal/clients/tg"
	"github.com/profectus200/contact-book-bot/internal/config"
	"github.com/profectus200/contact-book-bot/internal/database"
	"github.com/profectus200/contact-book-bot/internal/model/callbacks"
	"github.com/profectus200/contact-book-bot/internal/model/messages"
	"github.com/profectus200/contact-book-bot/internal/worker"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	config, err := config.New()
	if err != nil {
		log.Fatal("Cannot create config")
	}

	db, err := database.New(config)
	if err != nil {
		log.Fatal("Cannot create database")
	}

	contactsDB := database.NewContactsDB(db)
	usersDB := database.NewUsersDB(db)

	tgClient, err := tg.New(config)
	if err != nil {
		log.Fatal("Cannot create new tg-client")
	}

	msgModel := messages.New(tgClient, contactsDB, usersDB)
	callbackModel := callbacks.New(tgClient, contactsDB, usersDB)

	updateListenerWorker := worker.NewUpdateListenerWorker(tgClient, msgModel, callbackModel)

	updateListenerWorker.Run(ctx)
}
