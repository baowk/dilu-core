package cryptos

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestRsaSignWithhash(t *testing.T) {
	priK := `-----BEGIN PRIVATE KEY-----
...
-----END PRIVATE KEY-----`
	pk, err := ParsePriKeyPkcs8([]byte(priK))
	if err != nil {
		t.Error(err)
		return
	}
	data, err := RsaSignWithHash(pk, []byte(`{}`), 1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(base64.StdEncoding.EncodeToString(data))
}
