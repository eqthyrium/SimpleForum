package app

import (
	"SimpleForum/internal/config"
	"SimpleForum/internal/repository/sqllite"
	"SimpleForum/internal/service/repository"
	"SimpleForum/internal/service/usecase"
	"SimpleForum/internal/transport/customHttp"
	"SimpleForum/pkg/logger"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func RunApplication() {

	conf := config.NewConfiguration()
	LoggerObjectHttp := logger.NewLogger().GetLoggerObject("../../logging/info.log", "../../logging/error.log", "Middleware")

	db, err := openDb(*conf.Dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	serviceObject := &usecase.Application{
		ServiceDB: repository.NewServiceRepository(sqllite.NewRepository(db)),
	}

	handlerObject := &customHttp.HandlerHttp{
		Service:  customHttp.NewTransportHttpHandler(serviceObject),
		ErrorLog: LoggerObjectHttp.ErrorLogger,
		InfoLog:  LoggerObjectHttp.InfoLogger,
	}

	router := handlerObject.Routering()
	handlerObject.InfoLog.Println("The server is running at: http://localhost%s\n", *conf.Addr)
	log.Print(http.ListenAndServe(*conf.Addr, router))

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
