package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/profectus200/contact-book-bot/internal/config"
)

func New(service *config.Service) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s sslmode=%s",
		service.Config.Host,
		service.Config.Port,
		service.Config.User,
		service.Config.Password,
		service.Config.Database,
		service.Config.SslMode,
	)

	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return nil, errors.Wrap(err, "cannot Open")
	}

	return db, nil
}
