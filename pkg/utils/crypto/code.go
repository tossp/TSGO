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
	tmp, err := GzipEncode(input)
	if err == nil {
		input = tmp
	}
	return base64.StdEncoding.EncodeToString(input)
}

//Base64Decode 解码
func Base64Decode(input string) []byte {
	r, _ := base64.StdEncoding.DecodeString(input)
	tmp, err := GzipDecode(r)
	if err == nil {
		r = tmp
	}
	return r
}

func HexEncode(input []byte) string {
	return hex.EncodeToString(input)
}
func HexDecode(input string) []byte {
	r, _ := hex.DecodeString(input)
	return r
}

func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer = new(bytes.Buffer)
		out    []byte
		err    error
	)
	writer, _ := gzip.NewWriterLevel(buffer, gzip.BestCompression)
	_, err = writer.Write(in)
	if err != nil {
		_ = writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}
	return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
