package goeng

import "strconv"

func KeyDictIncrement(user_id int) string {
	return "dict_increment_" + strconv.Itoa(user_id)
}

func KeyDict(user_id int, dict_id int64) string {
	return "dict_" + strconv.Itoa(user_id) + "_" + strconv.FormatInt(dict_id, 10)
}

func KeyDictList(user_id int) string {
	return "dict_list_" + strconv.Itoa(user_id)
}

func KeyWordIncrement(user_id int, dict_id int64) string {
	return "word_increment_" + strconv.Itoa(user_id) + "_" + strconv.FormatInt(dict_id, 10)
}

func KeyWord(user_id int, dict_id int64, word_id int64) string {
	return "dict_" + strconv.Itoa(user_id) + "_" + strconv.FormatInt(dict_id, 10) + "_" + strconv.FormatInt(word_id, 10)
}

func KeyWordList(user_id int, dict_id int64) string {
	return "word_list_" + strconv.Itoa(user_id) + "_" + strconv.FormatInt(dict_id, 10)
}
