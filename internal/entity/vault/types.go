package vault

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

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
