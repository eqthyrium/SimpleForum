package app

import (
	"SimpleForum/internal/service/usecase"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"SimpleForum/internal/config"
	"SimpleForum/internal/repository/sqllite"
	"SimpleForum/internal/transport/customHttp"
	"SimpleForum/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

func RunApplication() {
	conf := config.NewConfiguration()

	LoggerObjectHttp := logger.NewLogger().GetLoggerObject("../logging/info.log", "../logging/error.log", "../logging/debug.log", "HTTP")
	db, err := openDb(*conf.Dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	repositoryObject := sqllite.NewRepository(db)
	serviceObject := usecase.NewUseCase(repositoryObject)
	httpTransport := customHttp.NewTransportHttpHandler(serviceObject, LoggerObjectHttp)

	router := httpTransport.Routering()
	message := fmt.Sprintf("The server is running at: http://localhost%s/\n", *conf.Addr)
	log.Print(message)
	httpTransport.InfoLog.Println(message)
	log.Fatalln(http.ListenAndServe(*conf.Addr, router))
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
