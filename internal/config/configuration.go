package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Config *Configuration

type Configuration struct {
	Addr        *string
	Dsn         *string
	GithubOauth *Github
	GoogleOauth *Google
}

type Github struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type Google struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func NewConfiguration() *Configuration {
	configObject := &Configuration{
		GithubOauth: new(Github),
		GoogleOauth: new(Google),
	}
	configObject.Addr = flag.String("addr", ":8888", "customHttp listen port")
	configObject.Dsn = flag.String("dsn", "../mydb.db", "Sqllite3 data source name")
	flag.Parse()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file, error message is:", err)
	}

	configObject.GithubOauth.ClientID = os.Getenv("GithubClientID")
	configObject.GithubOauth.ClientSecret = os.Getenv("GithubClientSecret")
	configObject.GithubOauth.RedirectURI = os.Getenv("GithubRedirectURI")

	configObject.GoogleOauth.ClientID = os.Getenv("GoogleClientID")
	configObject.GoogleOauth.ClientSecret = os.Getenv("GoogleClientSecret")
	configObject.GoogleOauth.RedirectURI = os.Getenv("GoogleRedirectURI")

	return configObject
}
