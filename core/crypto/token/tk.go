package token

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"

	"github.com/doublemo/foxchat/core/crypto/aes"
)

type TK struct {
	key []byte
}

func (tk *TK) Encrypt(data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		return "", err
	}

	hasher := hmac.New(sha256.New, tk.key)
	if _, err := hasher.Write(buf.Bytes()); err != nil {
		return "", err
	}

	if err := binary.Write(buf, binary.LittleEndian, hasher.Sum(nil)); err != nil {
		return "", err
	}

    frame, err := aes.Encrypt(buf.Bytes(), tk.key)
    if err != nil {
    	return "", err
    }

	return base64.StdEncoding.EncodeToString(frame), nil
}

func (tk *TK) Decrypt(s string, o interface{}) error {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	dataSize := binary.Size(o)
	if len(data) < dataSize+sha256.Size {
		return errors.New("Token is too short")
	}

	data, err = aes.Decrypt(data, tk.key)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.LittleEndian, o); err != nil {
		return err
	}

	hbuf := new(bytes.Buffer)
	binary.Write(hbuf, binary.LittleEndian, o)

	hasher := hmac.New(sha256.New, tk.key)
	hasher.Write(hbuf.Bytes())

	if !hmac.Equal(data[dataSize:dataSize+sha256.Size], hasher.Sum(nil)) {
		return errors.New("ErrFailed")
	}

	return nil
}

func NewTK(k []byte) *TK {
	return &TK{key: k}
}