package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type jwtUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

func (j *Auth) GenerateTokenPair(user *jwtUser) (TokenPairs, error) {
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"
	claims["isAdmin"] = user.IsAdmin

	// Set the expiry for JWT
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create a refresh token and set claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["username"] = user.Username
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()
	refreshTokenClaims["isAdmin"] = user.IsAdmin

	// Set the expiry for the refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	// Create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create TokenPairs and populate with signed tokens
	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}
	log.Println("Signed access token: ", tokenPairs.Token)
	log.Println("Signed refresh token: ", tokenPairs.RefreshToken)

	// Return TokenPairs
	return tokenPairs, nil
}

func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    refreshToken,
		Expires:  time.Now().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteNoneMode, // Set SameSite to None because front and back domain are different
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}

	// log.Println("Cookie Name: ", cookie.Name)
	// log.Println("Cookie Path: ", cookie.Path)
	// log.Println("Cookie Value: ", cookie.Value)
	// log.Println("Cookie Expires: ", cookie.Expires)
	// log.Println("Cookie MaxAge: ", cookie.MaxAge)
	// log.Println("Cookie SameSite: ", cookie.SameSite)
	// log.Println("Cookie Domain: ", cookie.Domain)
	// log.Println("Cookie HttpOnly: ", cookie.HttpOnly)
	// log.Println("Cookie Secure: ", cookie.Secure)

	return cookie
}

func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode, // Set SameSite to None because front and back domain are different
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

// for auth middleware, extracts auth header from request and validates token
func (j *Auth) GetAndVerifyHeaderToken(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	// indicates that the server's response will vary based on the value of the Authorization header.
	// good practice to include this if the response of endpoint can be different depending on the Authorization header.
	w.Header().Add("Vary", "Authorization")

	// the `Authorization` header from the HTTP request is extracted.
	authHeader := r.Header.Get("Authorization")

	// checks whether the Authorization header exists and is in the correct format.
	if authHeader == "" {
		log.Println("No auth header")
		return "", nil, errors.New("no auth header")
	}

	// Authorization header should have the format Bearer JWTtoken, so it's split by spaces and checks if the first part is Bearer.
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		log.Println("Auth header length incorrect")
		log.Println("Header parts: ", headerParts)
		return "", nil, errors.New("auth header length incorrect")
	}

	// checks if the first part is Bearer.
	if headerParts[0] != "Bearer" {
		log.Println("Invalid auth header")
		log.Println("Header: ", headerParts[0])
		return "", nil, errors.New("invalid auth header")
	}

	// JWT token is extracted from the Authorization header.
	token := headerParts[1]
	claims := &Claims{}

	// parse the JWT token and fill the claims struct with the claims in the token
	// and check correct signing method
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Println("Unexpected signing mathod")
			log.Println("Method: ", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	// checks if the token is expired by examining the error returned from ParseWithClaims.
	if err != nil {
		log.Println("Token error: ", err)
		if strings.HasPrefix(err.Error(), "token is expired by") {
			log.Println("Expired token")
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}

	// checks whether the issuer (iss) claim in the token matches the expected issuer.
	// if the issuer is not what's expected, it returns an error.
	if claims.Issuer != j.Issuer {
		log.Println("Issuer: ", claims.Issuer)
		log.Println("Expected issuer: ", j.Issuer)
		return "", nil, errors.New("invalid issuer")
	}

	return token, claims, nil
}
