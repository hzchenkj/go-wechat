package test

import (
	"bytes"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
)

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func FormatXML(data []byte) ([]byte, error) {
	b := &bytes.Buffer{}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	encoder := xml.NewEncoder(b)
	encoder.Indent("", "	")
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			encoder.Flush()
			return b.Bytes(), nil
		}
		if err != nil {
			return nil, err
		}
		err = encoder.EncodeToken(token)
		if err != nil {
			return nil, err
		}
	}
}

func StrToCDATA(str string) CDATAText {
	return CDATAText{"<![CDATA[" + str + "]]>"}
}

func Signature(timestamp string, nonce string, token string) string {
	strS := sort.StringSlice{token, timestamp, nonce}
	sort.Strings(strS)
	str := ""
	for _, s := range strS {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Send(url string, message string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/xml")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}