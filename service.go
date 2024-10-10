package goeng

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

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
		s.respError(c, http.StatusBadRequest, "incorrect data")
		return
	}

	user_name := strings.ToLower(dtoUser.UserName)
	user_name = strings.TrimSpace(user_name)
	password := strings.TrimSpace(dtoUser.Password)

	if utf8.RuneCountInString(user_name) < 4 {
		s.respError(c, http.StatusBadRequest, "Username can't be less 4 characters")
		return
	}

	if utf8.RuneCountInString(password) < 6 {
		s.respError(c, http.StatusBadRequest, "Password can't be less 6 characters")
		return
	}

	_, success := s.app.Get(KeyUserMatch(user_name))
	if success {
		s.respError(c, http.StatusNotAcceptable, "User name alredy used")
		return
	}

	userId, _ := s.app.Increment(KeyUserIncrement())
	salt := randStringRunes(5)

	user := User{ID: userId, UserName: user_name, PasswordHash: makePasswordHash(password, salt), Salt: salt}
	user_json, _ := json.Marshal(user)
	s.app.Set(KeyUser(userId), string(user_json))
	s.app.Set(KeyUserMatch(user_name), strconv.FormatInt(userId, 10))

	token, _ := jwtEncrypt(s.config.SecretKey, userId)
	c.SetCookie("auth_token", token, 3e7, "/", s.config.CoockieHost, false, true)

	s.app.CloseConnection()

	dtoSuccess := DtoSuccess{true}
	c.IndentedJSON(http.StatusOK, dtoSuccess)
}

func (s *Service) signIn(c *gin.Context) {
	s.app.NewConnection()

	var dtoUser DtoUser
	if err := c.BindJSON(&dtoUser); err != nil {
		s.respError(c, http.StatusBadRequest, "incorrect data")
		return
	}

	user_name := strings.ToLower(dtoUser.UserName)
	user_name = strings.TrimSpace(user_name)
	password := strings.TrimSpace(dtoUser.Password)

	userIdDb, success := s.app.Get(KeyUserMatch(user_name))
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

	if user.PasswordHash != makePasswordHash(password, user.Salt) {
		s.respError(c, http.StatusNotAcceptable, "incorrect password")
		return
	}

	token, _ := jwtEncrypt(s.config.SecretKey, usrId)
	c.SetCookie("auth_token", token, 3600, "/", s.config.CoockieHost, false, true)

	s.app.CloseConnection()

	dtoSuccess := DtoSuccess{true}
	c.IndentedJSON(http.StatusOK, dtoSuccess)
}

func (s *Service) createDict(c *gin.Context) {
	var newDict Dict

	s.app.NewConnection()

	if err := c.BindJSON(&newDict); err != nil {
		s.respError(c, http.StatusBadRequest, "incorrect data")
		return
	}

	user_id, err := s.getUserId(c)
	if err != nil {
		s.respError(c, http.StatusUnauthorized, "Auth fail")
		return
	}

	newDict.Name = strings.TrimSpace(newDict.Name)
	if utf8.RuneCountInString(newDict.Name) == 0 {
		s.respError(c, http.StatusBadRequest, "Empty name")
		return
	}

	increment, _ := s.app.Increment(KeyDictIncrement(user_id))
	newDict.ID = increment

	dict_json, _ := json.Marshal(newDict)

	s.app.Set(KeyDict(user_id, increment), string(dict_json))
	s.app.Push(KeyDictList(user_id), KeyDict(user_id, increment))

	c.IndentedJSON(http.StatusOK, newDict)
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

func (s *Service) deleteDict(c *gin.Context) {
	s.app.NewConnection()

	dict_id, err := strconv.ParseInt(c.Param("dict_id"), 10, 64)
	if err != nil {
		fmt.Println(1)
		return
	}

	user_id, err := s.getUserId(c)
	if err != nil {
		s.respError(c, http.StatusUnauthorized, "Auth fail")
		return
	}

	items, _ := s.app.Pull(KeyWordList(user_id, dict_id))
	for key, _ := range items {
		word_id, _ := strconv.ParseInt(key, 10, 64)
		s.app.Delete(KeyWord(user_id, dict_id, word_id))
	}

	s.app.Remove(KeyUser(user_id), KeyDict(user_id, dict_id))
	s.app.Delete(KeyDict(user_id, dict_id))

	s.app.CloseConnection()
	dtoSuccess := DtoSuccess{true}
	c.IndentedJSON(http.StatusOK, dtoSuccess)
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

	newWord.First = strings.TrimSpace(newWord.First)
	newWord.Second = strings.TrimSpace(newWord.Second)
	if utf8.RuneCountInString(newWord.First) == 0 || utf8.RuneCountInString(newWord.Second) == 0 {
		s.respError(c, http.StatusBadRequest, "Empty word")
		return
	}

	increment, _ := s.app.Increment(KeyWordIncrement(user_id, dict_id))
	newWord.ID = increment
	wordJson, _ := json.Marshal(newWord)
	s.app.Set(KeyWord(user_id, dict_id, increment), string(wordJson))
	s.app.Push(KeyWordList(user_id, dict_id), KeyWord(user_id, dict_id, increment))

	s.app.CloseConnection()
	c.IndentedJSON(http.StatusOK, newWord)
}

func (s *Service) deleteWord(c *gin.Context) {
	s.app.NewConnection()

	dict_id, err := strconv.ParseInt(c.Param("dict_id"), 10, 64)
	if err != nil {
		fmt.Println(1)
		return
	}

	word_id, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		fmt.Println(1)
		return
	}

	user_id, err := s.getUserId(c)
	if err != nil {
		s.respError(c, http.StatusUnauthorized, "Auth fail")
		return
	}

	s.app.Remove(KeyDict(user_id, dict_id), KeyWord(user_id, dict_id, word_id))
	s.app.Delete(KeyWord(user_id, dict_id, word_id))

	s.app.CloseConnection()
	dtoSuccess := DtoSuccess{true}
	c.IndentedJSON(http.StatusOK, dtoSuccess)
}

func (s *Service) respError(c *gin.Context, code int, message string) {
	c.IndentedJSON(code, DtoError{Code: code, Message: message})
}
