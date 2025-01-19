package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"pastepanda/internal/models"
)

type application struct {
	logger *slog.Logger
	pastes *models.PasteModel
}

func main() {
	addr := flag.String("addr", ":4000", "The HTTP network address")
	dsn := flag.String("dsn", "web:password@/pastepanda?parseTime=true", "MySQL Data Sourse Name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger: logger,
		pastes: &models.PasteModel{DB: db},
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
