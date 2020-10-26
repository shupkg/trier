package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func New() (*Service, error) {
	var cfg Config
	method := jwt.GetSigningMethod(cfg.Method)
	if method == nil {
		return nil, fmt.Errorf("jwt method %s is not valid", cfg.Method)
	}
	return &Service{
		method:    method,
		secretKey: []byte(cfg.SecretKey),
	}, nil
}

type Service struct {
	method    jwt.SigningMethod
	secretKey []byte
}

func (s *Service) Create(claims AuthClaims) (string, error) {
	token := jwt.NewWithClaims(s.method, claims)
	return token.SignedString(s.secretKey)
}

func (s *Service) Validate(tokenString string, claims *AuthClaims) error {
	_, err := jwt.ParseWithClaims(tokenString, claims, s.provideKey)
	return err
}

func (s *Service) provideKey(*jwt.Token) (interface{}, error) {
	return s.secretKey, nil
}

func (s *Service) FromContext(ctx context.Context) *AuthClaims {
	v := ctx.Value("user-auth")
	if v != nil {
		if c, ok := v.(*AuthClaims); ok {
			return c
		}
	}
	return nil
}

func (s *Service) SetContext(ctx context.Context, c *AuthClaims) context.Context {
	return context.WithValue(ctx, "user-auth", c)
}

type AuthClaims struct {
	ID    uint64 `json:"i,omitempty"`
	Uid   uint64 `json:"u,omitempty"`
	Iss   string `json:"s,omitempty"`
	R     string `json:"r,omitempty"`
	Exp   int64  `json:"e,omitempty"`
	Nonce string `json:"j,omitempty"`
}

func (a AuthClaims) Valid() error {
	if a.Exp <= time.Now().Unix() {
		return errors.New("token is expired")
	}
	if a.ID == 0 || a.Uid == 0 {
		return errors.New("identity is not set")
	}
	if a.Nonce == "" {
		return errors.New("token not safe")
	}
	return nil
}
