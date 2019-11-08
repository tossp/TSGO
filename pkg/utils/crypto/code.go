package crypto

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
)

//Base64Encode 编码
func Base64Encode(input []byte) string {
	//tmp, err := GzipEncode(input)
	//if err == nil {
	//    input = tmp
	//}
	return base64.StdEncoding.EncodeToString(input)
}

//Base64Decode 解码
func Base64Decode(input string) (data []byte, err error) {
	if data, err = base64.StdEncoding.DecodeString(input); err != nil {
		return
	}
	//tmp, err := GzipDecode(data)
	//if err == nil {
	//    data = tmp
	//}
	return
}

//Base64Encode 编码
func Base64UrlEncode(input []byte) string {
	//tmp, err := GzipEncode(input)
	//if err == nil {
	//    input = tmp
	//}
	return base64.URLEncoding.EncodeToString(input)
}

//Base64Decode 解码
func Base64UrlDecode(input string) (data []byte, err error) {
	if data, err = base64.URLEncoding.DecodeString(input); err != nil {
		return
	}
	//tmp, err := GzipDecode(data)
	//if err == nil {
	//    data = tmp
	//}
	return
}

func HexEncode(input []byte) string {
	return hex.EncodeToString(input)
}
func HexDecode(input string) []byte {
	r, _ := hex.DecodeString(input)
	return r
}

func GzipEncode(in []byte) (out []byte, err error) {
	buffer := new(bytes.Buffer)
	writer, _ := gzip.NewWriterLevel(buffer, gzip.BestCompression)
	defer func() {
		err = writer.Close()
		if err != nil {
			_ = writer.Close()
		}
	}()
	_, err = writer.Write(in)
	if err != nil {
		return
	}
	out = buffer.Bytes()
	return
}

func GzipDecode(in []byte) (out []byte, err error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return
	}
	defer func() {
		_ = reader.Close()
	}()
	out, err = ioutil.ReadAll(reader)
	return
}
