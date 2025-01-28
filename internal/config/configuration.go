package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	MaxImageSize = 20 * 1024 * 1024 // 20 MB
	UploadDir    = "../uploads/"    // Directory to store images
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}
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

	configObject.GithubOauth.ClientID = os.Getenv("GithubClient_ID")
	configObject.GithubOauth.ClientSecret = os.Getenv("GithubSecretKey")
	configObject.GithubOauth.RedirectURI = os.Getenv("GithubRedirectURI")

	configObject.GoogleOauth.ClientID = os.Getenv("GoogleClientID")
	configObject.GoogleOauth.ClientSecret = os.Getenv("GoogleClientSecret")
	configObject.GoogleOauth.RedirectURI = os.Getenv("GoogleRedirectURI")

	return configObject
}
