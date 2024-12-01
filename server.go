package goeng

import (
	"fmt"
	"net/http"
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

	router := gin.New()

	router.Static("/assets", "./static/assets")

	router.LoadHTMLFiles("./static/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/web/*url", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	api := router.Group("/api")
	{
		api.GET("/auth", service.auth)
		api.POST("/sign_in", service.signIn)
		api.POST("/sign_up", service.signUp)

		api.GET("/dict", service.getDictList)
		api.POST("/dict", service.createDict)
		api.GET("/dict/:id", service.getDict)
		api.DELETE("/dict/:id", service.deleteDict)

		api.POST("/word/:id", service.addWord)
		api.DELETE("/word/:dict_id/:word_id", service.deleteWord)

		api.GET("/knownword/:dict_id/:word_id", service.knownWord)
		api.GET("/unknownword/:dict_id/:word_id", service.unknownWord)
	}

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
