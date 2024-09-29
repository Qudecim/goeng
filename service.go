package goeng

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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

func (s *Service) getUserId(c *gin.Context) (int64, error) {

	cookie, err := c.Cookie("auth_token")
	if err != nil {
		return 0, err
	}

	user_id, err := jwtDecrypt(s.config.SecretKey, cookie)
	if err != nil {
		return 0, err
	}

	return user_id, nil
}

func (s *Service) auth(c *gin.Context) {
	s.app.NewConnection()

	token, err := c.Cookie("auth_token")

	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, nil)
		return
	}

	_, err = jwtDecrypt(s.config.SecretKey, token)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, nil)
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
}

func (s *Service) signUp(c *gin.Context) {
	s.app.NewConnection()

	var dtoUser DtoUser
	if err := c.BindJSON(&dtoUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	_, success := s.app.Get(KeyUserMatch(dtoUser.UserName))
	if success {
		c.IndentedJSON(http.StatusNotAcceptable, nil)
		return
	}

	userId, _ := s.app.Increment(KeyUserIncrement())
	salt := randStringRunes(5)
	user := User{ID: userId, UserName: dtoUser.UserName, PasswordHash: makePasswordHash(dtoUser.Password, salt), Salt: salt}
	user_json, _ := json.Marshal(user)
	s.app.Set(KeyUser(userId), string(user_json))
	s.app.Set(KeyUserMatch(dtoUser.UserName), strconv.FormatInt(userId, 10))

	token, _ := jwtEncrypt(s.config.SecretKey, userId)
	c.SetCookie("auth_token", token, 3600, "/", s.config.CoockieHost, false, true)

	s.app.CloseConnection()
	c.IndentedJSON(http.StatusOK, nil)
}

func (s *Service) signIn(c *gin.Context) {
	s.app.NewConnection()

	var dtoUser DtoUser
	if err := c.BindJSON(&dtoUser); err != nil {
		s.respError(c, http.StatusBadRequest, "incorrect data")
		return
	}

	userIdDb, success := s.app.Get(KeyUserMatch(dtoUser.UserName))
	if !success {
		s.respError(c, http.StatusNotAcceptable, "incorrect username or password")
		return
	}

	usrId, _ := strconv.ParseInt(userIdDb, 10, 64)
	userJson, success := s.app.Get(KeyUser(usrId))
	if !success {
		s.respError(c, http.StatusBadRequest, "internal error")
		return
	}

	var user User
	err := json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		s.respError(c, http.StatusBadRequest, "internal error")
		return
	}

	if user.PasswordHash != makePasswordHash(dtoUser.Password, user.Salt) {
		s.respError(c, http.StatusNotAcceptable, "incorrect password")
		return
	}

	token, _ := jwtEncrypt(s.config.SecretKey, usrId)
	c.SetCookie("auth_token", token, 3600, "/", s.config.CoockieHost, false, true)

	s.app.CloseConnection()
	c.IndentedJSON(http.StatusOK, nil)
}

func (s *Service) createDict(c *gin.Context) {
	var newDict Dict

	s.app.NewConnection()

	if err := c.BindJSON(&newDict); err != nil {
		return
	}

	user_id, err := s.getUserId(c)
	if err != nil {
		s.respError(c, http.StatusUnauthorized, "Auth fail")
		return
	}

	increment, _ := s.app.Increment(KeyDictIncrement(user_id))
	newDict.ID = increment

	dict_json, _ := json.Marshal(newDict)

	s.app.Set(KeyDict(user_id, increment), string(dict_json))
	s.app.Push(KeyDictList(user_id), KeyDict(user_id, increment))

	c.IndentedJSON(http.StatusOK, "okey")
	s.app.CloseConnection()
}

func (s *Service) getDict(c *gin.Context) {
	user_id, err := s.getUserId(c)
	if err != nil {
		s.respError(c, http.StatusUnauthorized, "Auth fail")
		return
	}

	s.app.NewConnection()

	dict_id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return
	}

	words := []Word{}
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
	user_id, err := s.getUserId(c)
	if err != nil {
		s.respError(c, http.StatusUnauthorized, "Auth fail")
		return
	}

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

	sort.Slice(dicts, func(i, j int) bool {
		return dicts[i].ID < dicts[j].ID
	})

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

	user_id, err := s.getUserId(c)
	if err != nil {
		s.respError(c, http.StatusUnauthorized, "Auth fail")
		return
	}

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

func (s *Service) respError(c *gin.Context, code int, message string) {
	c.IndentedJSON(code, DtoError{Code: code, Message: message})
}
