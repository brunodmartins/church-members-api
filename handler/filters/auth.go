package filters

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type AuthFilter struct {
	jwtmiddleware *jwtmiddleware.JWTMiddleware
}

func NewAuthFilter() *AuthFilter {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := viper.GetString("auth.aud")
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := viper.GetString("auth.iss")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	return &AuthFilter{
		jwtmiddleware: jwtMiddleware,
	}
}

func (auth *AuthFilter) checkAccessToken(token string) error {
	accessToken := viper.GetString("auth.token")
	switch {
	case accessToken == "":
		return errors.New("Accces token empty")
	case token == accessToken:
		return nil
	default:
		return errors.New("Access token invalid")
	}
}

func (auth *AuthFilter) Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.EscapedPath() == "/ping" {
			return
		}
		var err error
		if c.GetHeader("X-Token") != "" {
			err = auth.checkAccessToken(c.GetHeader("X-Token"))
		} else {
			err = auth.jwtmiddleware.CheckJWT(c.Writer, c.Request)
		}
		if err != nil {
			fmt.Println(err)
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
			return
		}
	}
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""

	resp, err := http.Get(viper.GetString("auth.jwk"))

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}
