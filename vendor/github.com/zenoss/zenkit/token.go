package zenkit

import (
	"encoding/json"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var (
	ErrorInvalidToken = errors.New("invalid JWT")
)

func ParseUnverified(token string, claims interface{}) (err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.Wrap(ErrorInvalidToken, "wrong number of segments")
	}
	var claimBytes []byte
	if claimBytes, err = jwt.DecodeSegment(parts[1]); err != nil {
		return errors.Wrap(err, "unable to decode claims segment of JWT")
	}
	if err := json.Unmarshal(claimBytes, claims); err != nil {
		return errors.Wrap(err, "unable to deserialize claims")
	}
	return nil
}
