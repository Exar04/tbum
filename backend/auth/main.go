package main

import (
	"auth/helper"
	"auth/initilizer"
	"auth/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func init() {
	initilizer.ConnectToDB()
	initilizer.SyncDb()
}

func main() {
	port := os.Getenv("AUTH_SERVER_PORT")

	if port == "" {
		port = "9000"
	}
	router := mux.NewRouter()

	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/signup", Signup).Methods("POST")
	router.HandleFunc("/validate", validate).Methods("POST")

	log.Println("Auth service listining on port", port)
	http.ListenAndServe(":"+port, router)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &user)

	if user.Username == "" {
		fmt.Println("username field is empty!")
	}

	var gotPass = user.Password
	result := initilizer.DB.Where("username = ?", user.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			w.Write([]byte("user doesn't exist"))
			return
		} else {
			fmt.Println("Error querying database:", result.Error)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if user.Password == gotPass {
		expTime := time.Now().Add(time.Minute * 5)
		claims := &models.Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(helper.JwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		w.Write([]byte(tokenString))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &user)
	if user.Username == "" {
		fmt.Println("username field is empty!")
	}
	result := initilizer.DB.Where("username = ?", user.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			newUser := models.User{
				Username: user.Username,
				Password: user.Password,
			}
			initilizer.DB.Create(&newUser)
		} else {
			fmt.Println("Error querying database:", result.Error)
		}
	} else {
		fmt.Println("Username already exists.")
	}
}

func validate(w http.ResponseWriter, r *http.Request) {
	var user models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &user)

	if user.Username == "" {
		fmt.Println("username field is empty!")
	}
	if user.Token == "" {
		fmt.Println("token field is empty!")
	}

	tokenStr := user.Token
	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) { return helper.JwtKey, nil },
	)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte("token was valid"))
}

func refresh(w http.ResponseWriter, r *http.Request) {
	var user models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &user)

	// if user.Username == "" {
	// 	fmt.Println("username field is empty!")
	// }
	if user.Token == "" {
		fmt.Println("token field is empty!")
	}

	tokenStr := user.Token
	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) { return helper.JwtKey, nil },
	)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expTime := time.Now().Add(time.Minute * 5)
	claims.ExpiresAt = expTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(helper.JwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Write([]byte(tokenString))
}
