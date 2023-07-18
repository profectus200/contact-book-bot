package main

import (
	"context"
	"github.com/profectus200/contact-book-bot/cmd/logging"
	"github.com/profectus200/contact-book-bot/cmd/tracing"
	"go.uber.org/zap"
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

	logger := logging.InitLogger()
	tracing.InitTracing("actions_handler", logger)

	logger.Info("Initializing config")
	config, err := config.New()
	if err != nil {
		logger.Fatal("Cannot create config", zap.Error(err))
	}

	logger.Info("Initializing database")
	db, err := database.New(config)
	if err != nil {
		logger.Fatal("Cannot create database", zap.Error(err))
	}

	contactsDB := database.NewContactsDB(db)
	usersDB := database.NewUsersDB(db)

	logger.Info("Initializing tg client")
	tgClient, err := tg.New(config)
	if err != nil {
		logger.Fatal("Cannot create new tg client", zap.Error(err))
	}

	msgModel := messages.New(tgClient, contactsDB, usersDB)
	callbackModel := callbacks.New(tgClient, contactsDB, usersDB)

	updateListenerWorker := worker.NewUpdateListenerWorker(tgClient, msgModel, callbackModel)

	updateListenerWorker.Run(ctx)
}
