package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	private    *rsa.PrivateKey
	public     *rsa.PublicKey
	expiration time.Duration
}

type Claims struct {
	UserID string `json:"user_id"`
	Roles  string `json:"roles"`
	jwt.RegisteredClaims
}

func NewManager(private *rsa.PrivateKey, public *rsa.PublicKey, expirationMinutes int) *Manager {
	return &Manager{
		private:    private,
		public:     public,
		expiration: time.Duration(expirationMinutes) * time.Minute,
	}
}

func (m *Manager) GenerateToken(userID string, roles string) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(m.expiration)

	claims := Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(m.private)
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt.Unix(), nil
}

func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.public, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
