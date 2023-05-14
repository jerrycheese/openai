package wechat

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type MsgEncrypted struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
	Content      string   `xml:"Content"`
	Recognition  string   `xml:"Recognition"`

	MsgId int64 `xml:"MsgId,omitempty"`

	AESKey string `xml:"-"`
}

func NewMsgEncrypted(aesKey string, dataWithEncryptedContent []byte) (*MsgEncrypted, error) {
	var msg MsgEncrypted
	if err := xml.Unmarshal(dataWithEncryptedContent, &msg); err != nil {
		return nil, err
	}
	if c, err := AESDecrypt(aesKey, msg.Content); err != nil {
		return nil, err
	} else {
		msg.Content = c
	}

	msg.AESKey = aesKey
	return &msg, nil
}

func (msg *MsgEncrypted) GenerateEchoData(s string) ([]byte, error) {
	encrypted, err := AESEncrypt(msg.AESKey, s)
	if err != nil {
		return nil, err
	}
	data := Msg{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      encrypted,
	}
	bs, _ := xml.Marshal(&data)
	return bs, nil
}

func AESEncrypt(key, text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AESDecrypt(key, ciphertext string) (string, error) {
	c, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(c) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := c[:aes.BlockSize]
	c = c[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(c, c)

	return string(c), nil
}
