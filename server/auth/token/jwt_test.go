package token

import (
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	pkFile, err := os.Open("../private.key")
	if err != nil {
		t.Fatalf("cannot open private key:%v", err)
	}
	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		t.Fatalf("cannout read private key:%v", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		t.Fatalf("cannot parse private key:%v", err)
	}
	g := NewJWTTokenGen("server/auth", key)
	g.nowFunc = func() time.Time {
		return time.Unix(1516239022, 0)
	}
	tkn, err := g.GenerateToken("o01eX5CyKGkgrzaswbKGBLQUWIQg", 2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate token: %v", err)
	}
	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoic2VydmVyL2F1dGgiLCJzdWIiOiJvMDFlWDVDeUtHa2dyemFzd2JLR0JMUVVXSVFnIn0.V9L6U5rZR1h_viKfN_r29F6PrRra5CycwZAjHVBh9ix3DbrXfGhu47wr7PZx05trouc3TAPqI1pagBDsDtcMbZyRcmRrFPodPNTelx74mJagcog07OEHTOCmqp1SohS7UxjAATJ2IdjE5S9S0ScO4zue6x1ExK5bsUK6QUCtdVrcDhc8LQ_YikkLL8PI11YTygIn7KI4L-r90eCN0JKgFBEZXmqd3oLCYNpk9pjWbztK70ECa7HSgopkTQx08ilH3Kk2ko10n4NjEB4OOxibHpY_zW1Xk6l9tCJco5PZnL6d6euXon72jYiYaH8lLWZpAH3VI5SMUHivefpotHARBQ"
	if tkn != want {
		t.Errorf("wrong token generated. want: %q, got: %q", want, tkn)
	}
}
