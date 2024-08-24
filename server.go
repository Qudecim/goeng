package goeng

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ssl            bool   `yaml:"ssl"`
	host           string `yaml:"host"`
	coockieHost    string `yaml:"coockie_host"`
	certificate    string `yaml:"certificate"`
	certificateKey string `yaml:"certificate_key"`
	secretKey      string `yaml:"secret_key"`
}

func Main() {
	config, err := readConfig()
	if err != nil {
		fmt.Println("Config error")
	}

	service := newService(config)

	router := gin.Default()

	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	router.GET("/api/sign_in", service.signIn)
	router.GET("/api/sign_up", service.signUp)

	router.GET("/api/dict", service.getDictList)
	router.GET("/api/dict/:id", service.getDict)
	router.POST("/api/dict", service.createDict)

	router.POST("/api/word/:id", service.addWord)

	if config.ssl {
		router.RunTLS(config.host, config.certificate, config.certificateKey)
	} else {
		router.Run(config.host)
	}
}

func readConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
