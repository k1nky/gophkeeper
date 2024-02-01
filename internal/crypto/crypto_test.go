package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

func TestT1(t *testing.T) {
	buf := bytes.NewBufferString("not such long secret text but very noble and intresting")
	k := sha256.Sum256([]byte("secret"))
	fmt.Println(hex.EncodeToString(k[:]))
	enc, _ := NewEncryptReader("secret", buf)
	out := bytes.NewBuffer(nil)
	n, err := out.ReadFrom(enc)
	os.WriteFile("/tmp/txt.enc", out.Bytes(), 0666)
	fmt.Println(n, err, out)
	f, _ := os.ReadFile("/tmp/txt.enc")
	buf = bytes.NewBuffer(f)
	dec, _ := NewDecryptReader("secret", buf)
	out = bytes.NewBuffer(nil)
	out.ReadFrom(dec)
	fmt.Println(out.String())
}
