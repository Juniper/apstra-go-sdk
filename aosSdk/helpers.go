package aosSdk

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func intSliceContains(in []int, t int) bool {
	for _, i := range in {
		if i == t {
			return true
		}
	}
	return false
}

func intSliceToStringSlice(in []int) []string {
	var result []string
	for _, i := range in {
		result = append(result, strconv.Itoa(i))
	}
	return result
}

func pp(in interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(in); err != nil {
		return err
	}
	return nil
}

type jwt struct {
	raw       string
	header    jwtHeader
	payload   jwtPayload
	signature []byte
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtPayload struct {
	Username    string `json:"username"`
	CreatedAt   string `json:"created_at"`
	UserSession string `json:"user_session"`
	Exp         int64  `json:"exp"`
}

func (o *jwt) decode() error {
	parts := strings.Split(o.raw, ".")
	if len(parts) != 3 {
		return fmt.Errorf("error decoding jwt, expected 3 string fields, got %d", len(parts))
	}

	headerJson, err := decodePart(parts[0])
	if err != nil {
		return fmt.Errorf("error decoding jwt header - %v", err)
	}
	err = json.Unmarshal(headerJson, &o.header)
	if err != nil {
		return fmt.Errorf("error unmarshaling jwt header - %v", err)
	}

	payloadJson, err := decodePart(parts[1])
	if err != nil {
		return fmt.Errorf("error decoding jwt payload - %v", err)
	}
	err = json.Unmarshal(payloadJson, &o.payload)
	if err != nil {
		return fmt.Errorf("error unmarshaling jwt payload - %v", err)
	}

	o.signature, err = decodePart(parts[2])
	return err
}

func decodePart(in string) ([]byte, error) {
	if l1 := len(in) % 4; l1 > 0 {
		in += strings.Repeat("=", 4-l1)
	}
	return base64.URLEncoding.DecodeString(in)
}

func (o jwt) Raw() string {
	return o.raw
}

func newJwt(in string) (*jwt, error) {
	token := jwt{raw: in}
	return &token, token.decode()
}

func (o jwt) expires() time.Time {
	return time.Unix(o.payload.Exp, 0)
}
