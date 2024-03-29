package vault

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// SecretType тип секрета.
type SecretType int

const (
	// Простой текст
	TypeText SecretType = iota
	// Логин - пароль
	TypeLoginPassword
	// Банковская карта
	TypeCreditCard
	// Файл
	TypeFile
)

// LoginPassword секрет как "логин-пароль"
type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CreditCard секрет как "банковская карта"
type CreditCard struct {
	Number     string `json:"number"`
	Holder     string `json:"holder"`
	CVV        string `json:"cvv"`
	Expiration string `json:"expiration"`
}

// StringPrompt запращивает ввод данных как строку от пользователя
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

// Запрос от пользователя логина-пароля для секрета.
func (lp *LoginPassword) Prompt() error {
	lp.Login = StringPrompt("login")
	lp.Password = StringPrompt("password")
	return nil
}

func (lp *LoginPassword) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(lp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Запрос от пользователя данных банковской карты для секрета.
func (cc *CreditCard) Prompt() error {
	cc.Number = StringPrompt("Number")
	cc.Holder = StringPrompt("Holder")
	cc.CVV = StringPrompt("CVV")
	cc.Expiration = StringPrompt("Expiration")
	return nil
}

func (cc *CreditCard) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(cc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t SecretType) String() string {
	switch t {
	case TypeText:
		return "TEXT"
	case TypeLoginPassword:
		return "LOGIN_PASSWORD"
	case TypeCreditCard:
		return "CREDIT_CARD"
	case TypeFile:
		return "FILE"
	}
	return "UNKNOWN"
}
