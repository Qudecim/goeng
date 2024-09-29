package goeng

import "strconv"

func KeyUser(userId int64) string {
	return "user_" + strconv.FormatInt(userId, 10)
}

func KeyUserMatch(userName string) string {
	return "user_math_" + userName
}

func KeyUserIncrement() string {
	return "user_increment"
}

func KeyDictIncrement(user_id int64) string {
	return "dict_increment_" + strconv.FormatInt(user_id, 10)
}

func KeyDict(user_id int64, dict_id int64) string {
	return "dict_" + strconv.FormatInt(user_id, 10) + "_" + strconv.FormatInt(dict_id, 10)
}

func KeyDictList(user_id int64) string {
	return "dict_list_" + strconv.FormatInt(user_id, 10)
}

func KeyWordIncrement(user_id int64, dict_id int64) string {
	return "word_increment_" + strconv.FormatInt(user_id, 10) + "_" + strconv.FormatInt(dict_id, 10)
}

func KeyWord(user_id int64, dict_id int64, word_id int64) string {
	return "dict_" + strconv.FormatInt(user_id, 10) + "_" + strconv.FormatInt(dict_id, 10) + "_" + strconv.FormatInt(word_id, 10)
}

func KeyWordList(user_id int64, dict_id int64) string {
	return "word_list_" + strconv.FormatInt(user_id, 10) + "_" + strconv.FormatInt(dict_id, 10)
}
