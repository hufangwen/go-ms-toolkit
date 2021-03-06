/*
@Time : 1/2/2021 公元 10:28
@Author : philiphu
@File : json
@Software: GoLand
*/
package common

import (
   jsoniter "github.com/json-iterator/go"
)

var jsonite = jsoniter.ConfigCompatibleWithStandardLibrary

func Unmarshal(data []byte, v interface{}) error {
	return jsonite.Unmarshal(data, v)
}

func UnmarshalFromString(str string, v interface{}) error {
	return jsonite.UnmarshalFromString(str, v)
}

func MarshalToString(v interface{}) (string, error) {
	return jsonite.MarshalToString(v)
}

func Marshal(v interface{}) ([]byte, error) {
	return jsonite.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return jsonite.MarshalIndent(v, prefix, indent)
}

func Get(data []byte, path ...interface{}) jsoniter.Any {
	return jsonite.Get(data, path...)
}

func Valid(data []byte) bool {
	return jsonite.Valid(data)
}



// ParseJson parsing json strings
func ParseJson(data string, result interface{}) error {
	return ParseJsonFromBytes([]byte(data), result)
}

// StringifyJson json to string
func StringifyJson(data interface{}) string {
	return string(StringifyJsonToBytes(data))
}

// ParseJsonFromBytes parsing json bytes
func ParseJsonFromBytes(data []byte, result interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, result)
}

// json bytes to string
func StringifyJsonToBytes(data interface{}) []byte {
	b, _ := StringifyJsonToBytesWithErr(data)
	return b
}

func StringifyJsonToBytesWithErr(data interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(&data)
	return b, err
}