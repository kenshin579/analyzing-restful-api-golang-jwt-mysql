package auth

import (
	"encoding/json"
	"fmt"
	"github.com/kenshin579/analyzing-restful-api-golang-jwt-mysql/app/models"
	"github.com/kenshin579/analyzing-restful-api-golang-jwt-mysql/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

type JwtToken struct {
	AccessToken string `json:"access-token"`
}

var jwt_secret = os.Getenv("jwt_secret")

func Login(w http.ResponseWriter, req *http.Request) {
	user := &models.User{}
	fmt.Println("req.Body", req.Body)

	err := json.NewDecoder(req.Body).Decode(user)

	//todo : 여기서 오류가 발생함
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	defer req.Body.Close()
	username := user.Username
	password := user.Password

	/* Another way to grab the form inputs from the request
	req.ParseForm()
	username := req.FormValue("Username")
	password := req.FormValue("Password")
	*/
	result := models.GetUsername(username)
	if result == nil {
		utils.Respond(w, utils.Message(false, "Your credentials do not match our records"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))

	if err != nil {
		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Your credentials do not match our records"))
		return
	}
	// access token ttl
	ttl := 2 * time.Minute
	accessTokenExpire := os.Getenv("access_token_expire")
	min, err := strconv.Atoi(accessTokenExpire)
	if err != nil {
		log.Println(err)
	}
	if accessTokenExpire != "" {
		ttl = time.Duration(min) * time.Minute
	}
	CreateToken(w, username, password, ttl)
}

func CreateToken(w http.ResponseWriter, username string, password string, ttl time.Duration) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      ttl,
	})

	tokenString, error := token.SignedString([]byte(jwt_secret))
	if error != nil {
		fmt.Println(error)
	}
	resp := utils.Message(true, "success")
	resp["data"] = JwtToken{AccessToken: tokenString}
	utils.Respond(w, resp)
	return
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(jwt_secret), nil
				})
				if error != nil {
					utils.Respond(w, utils.Message(false, error.Error()))
					return
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					utils.Respond(w, utils.Message(false, "Invalid authorization token"))
					return
				}
			}
		} else {
			utils.Respond(w, utils.Message(false, "An authorization header is required"))
			return
		}
	})
}
