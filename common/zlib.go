/*
@Time : 3/2/2021 公元 02:45
@Author : philiphu
@File : zap
@Software: GoLand
*/
package common

import (
	"bytes"
	"compress/zlib"
)

func ZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_, err := w.Write(src)
	if err != nil {
		return src
	}
	_ = w.Close()
	return in.Bytes()
}
