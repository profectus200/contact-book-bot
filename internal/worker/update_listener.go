package worker

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/profectus200/contact-book-bot/internal/model/callbacks"
	"github.com/profectus200/contact-book-bot/internal/model/messages"
	"log"
)

type updateFetcher interface {
	Start() tgbotapi.UpdatesChannel
	Request(callback tgbotapi.CallbackConfig) error
	Stop()
}

type MessageHandler interface {
	IncomingMessage(ctx context.Context, msg *messages.Message) error
}

type CallbackHandler interface {
	IncomingCallback(ctx context.Context, callback *callbacks.CallbackData) error
}

type UpdateListenerWorker struct {
	updateFetcher   updateFetcher
	messageHandler  MessageHandler
	callbackHandler CallbackHandler
}

func NewUpdateListenerWorker(updateFetcher updateFetcher,
	messageHandler MessageHandler, callbackHandler CallbackHandler) *UpdateListenerWorker {
	return &UpdateListenerWorker{
		updateFetcher:   updateFetcher,
		messageHandler:  messageHandler,
		callbackHandler: callbackHandler,
	}
}

func (w *UpdateListenerWorker) Run(ctx context.Context) {
	updates := w.updateFetcher.Start()

	for {
		select {
		case <-ctx.Done():
			w.updateFetcher.Stop()
			return
		case update, ok := <-updates:
			if !ok {
				w.updateFetcher.Stop()
				return
			}
			err := w.HandleUpdate(ctx, update)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (w *UpdateListenerWorker) HandleUpdate(ctx context.Context, update tgbotapi.Update) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"HandleUpdate",
	)
	defer span.Finish()

	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		err := w.messageHandler.IncomingMessage(ctx, &messages.Message{
			Text:      update.Message.Text,
			UserID:    update.Message.From.ID,
			MessageID: update.Message.MessageID,
		})

		if err != nil {
			return errors.Wrap(err, "cannot IncomingMessage")
		}
	} else if update.CallbackQuery != nil {
		log.Printf("[%s] data: %s",
			update.CallbackQuery.Message.From.UserName,
			update.CallbackQuery.Data,
		)

		err := w.callbackHandler.IncomingCallback(ctx, &callbacks.CallbackData{
			Data:       update.CallbackData(),
			FromID:     update.CallbackQuery.From.ID,
			MessageID:  update.CallbackQuery.Message.MessageID,
			CallbackID: update.CallbackQuery.ID,
		})

		if err != nil {
			return errors.Wrap(err, "cannot IncomingCallback")
		}
	}

	return nil
}
