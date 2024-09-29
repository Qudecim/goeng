package goeng

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Ssl            bool   `yaml:"ssl"`
	Host           string `yaml:"host"`
	CoockieHost    string `yaml:"coockie_host"`
	Certificate    string `yaml:"certificate"`
	CertificateKey string `yaml:"certificate_key"`
	SecretKey      string `yaml:"secret_key"`
}

func Main() {
	config, err := readConfig()
	if err != nil {
		fmt.Println("Config error")
	}

	service := newService(config)

	router := gin.Default()

	router.Static("/assets", "./static/assets")
	router.StaticFile("/", "./static/index.html")

	router.GET("/api/auth", service.auth)
	router.POST("/api/sign_in", service.signIn)
	router.POST("/api/sign_up", service.signUp)

	router.GET("/api/dict", service.getDictList)
	router.POST("/api/dict", service.createDict)
	router.GET("/api/dict/:id", service.getDict)

	router.POST("/api/word/:id", service.addWord)

	if config.Ssl {
		router.RunTLS(config.Host, config.Certificate, config.CertificateKey)
	} else {
		router.Run(config.Host)
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
	fmt.Println(config)

	return &config, nil
}
