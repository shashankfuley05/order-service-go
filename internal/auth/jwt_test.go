package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {

	token, err := GenerateToken("USER-1")

	if err != nil {
		t.Fatalf("BC token nahi gnerate hua karan %v", err.Error())
	}

	if token == "" {
		t.Fatalf("BC ye empty token kaise aagya")
	}

}
func TestTokenParsing(t *testing.T) {

	tests := []struct {
		name        string
		tokenString string

		isErrorExpected bool
		expectedUser    string
	}{
		{
			name:            "Invalid Signing Method",
			tokenString:     createTestToken(jwt.SigningMethodHS384, time.Hour, jwt.MapClaims{}),
			isErrorExpected: true,
		},
		{
			name:            "Expired or Invalid Token",
			tokenString:     createTestToken(jwt.SigningMethodHS256, -time.Hour, jwt.MapClaims{}),
			isErrorExpected: true,
		},
		{
			name:            "Token Valid but User not found",
			tokenString:     createTestToken(jwt.SigningMethodHS256, time.Hour, jwt.MapClaims{}),
			isErrorExpected: true,
		},
		{
			name: "Token Valid user found",
			tokenString: createTestToken(jwt.SigningMethodHS256, time.Hour, jwt.MapClaims{
				"userId": "TEST_USER",
			}),
			isErrorExpected: false,
			expectedUser:    "TEST_USER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userId, err := ParseToken(tt.tokenString)

			if tt.isErrorExpected {
				if err == nil {
					t.Fatalf("Error was expected, but error is nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Error was not expected , but error received %v", err.Error())
			}

			if tt.expectedUser != userId {
				t.Errorf("Unexpected user received, expected %s received %s", tt.expectedUser, userId)
			}
		})
	}
}

func createTestToken(
	method jwt.SigningMethod,
	expiry time.Duration,
	claims jwt.MapClaims,
) string {

	claims["exp"] = time.Now().Add(expiry).Unix()

	token := jwt.NewWithClaims(method, claims)

	tokenString, _ := token.SignedString([]byte("AbraKadabra"))

	return tokenString
}

func TestGenerateTokenWithClaimsForSuccess(t *testing.T) {

	testUser := "TEST_USER"
	token, err := GenerateTokenWithClaims(testUser, "ADMIN", "ORDER_MANAGER")

	if err != nil {
		t.Fatalf("Token couldn't be generate.")
	}

	userClaims, err := ParseTokenWithClaims(token)

	if err != nil {
		t.Fatalf("Error while parsing token %v", err.Error())
	}

	if testUser != userClaims.UserId {
		t.Errorf("User expected was %v , but was received %v", testUser, userClaims.UserId)
	}

	roles := userClaims.Roles

	if len(roles) != 2 {
		t.Fatalf("Expected roles were 2, but recieved %d", len(roles))
	}

	if roles[0] != "ADMIN" {
		t.Errorf("Expected roles was ADMIN but received %v", roles[0])
	}

	if roles[1] != "ORDER_MANAGER" {
		t.Errorf("Expected role was ORDER_MANAGER but received %v", roles[0])
	}
}
