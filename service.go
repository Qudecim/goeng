package goeng

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Qudecim/ipmc"
	"github.com/gin-gonic/gin"
)

type Service struct {
	app    *ipmc.App
	config *Config
}

func newService(config *Config) *Service {

	app := ipmc.NewApp(ipmc.NewConfig("binlog/", 1e5, "snapshot/"))
	app.Init()

	return &Service{app, config}
}

func (s *Service) getUserId(c *gin.Context) int {
	return 1
}

func (s *Service) signUp(c *gin.Context) {
	s.app.NewConnection()

	var dtoUser DtoUser
	if err := c.BindJSON(&dtoUser); err != nil {
		return
	}

	_, success := s.app.Get(KeyUserMatch(dtoUser.userName))
	if success {
		// exist
		return
	}

	userId, _ := s.app.Increment(KeyUserIncrement())
	salt := randStringRunes(5)
	user := User{ID: userId, userName: dtoUser.userName, passwordHash: makePasswordHash(dtoUser.password, salt), salt: salt}
	user_json, _ := json.Marshal(user)
	s.app.Set(KeyUser(userId), string(user_json))

	token, _ := jwtEncrypt(s.config.secretKey, userId)
	c.SetCookie("auth_token", token, 3600, "/", s.config.coockieHost, false, true)

	s.app.CloseConnection()
}

func (s *Service) signIn(c *gin.Context) {
	s.app.NewConnection()

	var dtoUser DtoUser
	if err := c.BindJSON(&dtoUser); err != nil {
		return
	}

	userIdDb, success := s.app.Get(KeyUserMatch(dtoUser.userName))
	if !success {
		// not exist
		return
	}

	usrId, _ := strconv.ParseInt(userIdDb, 10, 64)
	userJson, success := s.app.Get(KeyUser(usrId))
	if !success {
		// not exist
		return
	}

	var user User
	err := json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		// todo
		return
	}

	if user.passwordHash != makePasswordHash(dtoUser.password, user.salt) {
		// todo
		return
	}

	token, _ := jwtEncrypt(s.config.secretKey, usrId)
	c.SetCookie("auth_token", token, 3600, "/", s.config.coockieHost, false, true)

	s.app.CloseConnection()
}

func (s *Service) createDict(c *gin.Context) {
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

func (s *Service) getDict(c *gin.Context) {
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

func (s *Service) getDictList(c *gin.Context) {
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

func (s *Service) addWord(c *gin.Context) {
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
