package goeng

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Qudecim/ipmc"
	"github.com/gin-gonic/gin"
)

type service struct {
	app *ipmc.App
}

func newService() *service {
	config := ipmc.NewConfig("binlog/", 1e5, "snapshot/")
	app := ipmc.NewApp(config)
	app.Init()

	return &service{app}
}

func (s *service) getUserId(c *gin.Context) int {
	return 1
}

func (s *service) createDict(c *gin.Context) {
	var newDict Dict

	s.app.NewConnection()

	if err := c.BindJSON(&newDict); err != nil {
		return
	}

	user_id := s.getUserId(c)

	increment, _ := s.app.Increment(KeyDictIncrement(user_id))
	newDict.ID = increment

	dict_json, _ := json.Marshal(newDict)

	s.app.Set(KeyDict(user_id, increment), string(dict_json))
	s.app.Push(KeyDictList(user_id), KeyDict(user_id, increment))

	c.IndentedJSON(http.StatusOK, "okey")
	s.app.CloseConnection()
}

func (s *service) getDict(c *gin.Context) {
	user_id := s.getUserId(c)
	s.app.NewConnection()

	dict_id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return
	}

	var words []Word
	items, _ := s.app.Pull(KeyWordList(user_id, dict_id))
	for _, item := range items {
		word := Word{}
		if err := json.Unmarshal([]byte(item), &word); err != nil {
			panic(err)
		}
		words = append(words, word)
	}

	s.app.CloseConnection()
	c.IndentedJSON(http.StatusOK, words)
}

func (s *service) getDictList(c *gin.Context) {
	user_id := s.getUserId(c)
	s.app.NewConnection()

	var dicts []Dict
	items, _ := s.app.Pull(KeyDictList(user_id))
	for _, item := range items {
		dict := Dict{}
		if err := json.Unmarshal([]byte(item), &dict); err != nil {
			panic(err)
		}
		dicts = append(dicts, dict)
	}

	s.app.CloseConnection()
	c.IndentedJSON(http.StatusOK, dicts)
}

func (s *service) addWord(c *gin.Context) {
	s.app.NewConnection()

	dict_id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		fmt.Println(1)
		return
	}

	user_id := s.getUserId(c)

	var newWord Word
	if err := c.BindJSON(&newWord); err != nil {
		fmt.Println(2)
		return
	}
	increment, _ := s.app.Increment(KeyWordIncrement(user_id, dict_id))
	newWord.ID = increment
	wordJson, _ := json.Marshal(newWord)
	s.app.Set(KeyWord(user_id, dict_id, increment), string(wordJson))
	s.app.Push(KeyWordList(user_id, dict_id), KeyWord(user_id, dict_id, increment))

	s.app.CloseConnection()
	c.IndentedJSON(http.StatusOK, "okey")
}
