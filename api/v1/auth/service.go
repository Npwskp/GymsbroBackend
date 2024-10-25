package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/user"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DB *mongo.Database
}

type IAuthService interface {
	Login(login *LoginDto) (string, int64, error)
	Register(register *RegisterDto) (*user.User, error)
	Me(token string) (*user.User, error)
}

func (as *AuthService) Login(login *LoginDto) (string, int64, error) {
	if login.Email == "" || login.Password == "" {
		return "", 0, errors.New("plese enter your email and password")
	}
	userService := user.UserService{DB: as.DB}
	user, err := userService.GetUserByEmail(login.Email)
	if err != nil {
		return "", 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return "", 0, err
	}

	token, exp, err := createJWTToken(user)
	if err != nil {
		return "", 0, err
	}

	return token, exp, nil
}

func (as *AuthService) Register(register *RegisterDto) (*user.User, error) {
	if register.Email == "" || register.Password == "" {
		return nil, errors.New("plese enter your email and password")
	}
	userService := user.UserService{DB: as.DB}
	check, err := userService.GetUserByEmail(register.Email)
	fmt.Println("hello2")
	if check != nil && err == nil {
		return nil, errors.New("email have been used")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	register.Password = string(hashedPassword)

	user := user.CreateUserDto{
		Username: register.Username,
		Email:    register.Email,
		Password: register.Password,
		Age:      register.Age,
		Gender:   register.Gender,
	}

	fmt.Println(user)
	result, err := userService.CreateUser(&user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("hello3")
	fmt.Println(result)
	return result, nil
}

func (as *AuthService) Me(token string) (*user.User, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	userService := user.UserService{DB: as.DB}
	user, err := userService.GetUser(claims["id"].(string))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func createJWTToken(user *user.User) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	exp := time.Now().Add(time.Hour * 72).Unix()
	claims["id"] = user.ID.Hex()
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["password"] = user.Password
	claims["exp"] = exp
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", 0, err
	}
	return tokenString, exp, nil
}
