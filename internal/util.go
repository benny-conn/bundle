package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WriteResponse(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func NewCaseInsensitiveRegex(value string) primitive.Regex {
	return primitive.Regex{Pattern: value, Options: "i"}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/bundle/")
	viper.AddConfigPath("$HOME/.bundle")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
}
