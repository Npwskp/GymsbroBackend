package auth

import (
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/Npwskp/GymsbroBackend/api/v1/config"
	"github.com/Npwskp/GymsbroBackend/api/v1/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

	// Validate password strength
	if err := validatePasswordStrength(register.Password); err != nil {
		return nil, err
	}

	// Hash password with higher cost
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(register.Password),
		config.BcryptCost,
	)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Clear plain text password immediately
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
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTSecret()), nil
	})
	if err != nil {
		return nil, err
	}
	userService := user.UserService{DB: as.DB}
	user, err := userService.GetUser(claims["sub"].(string))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func createJWTToken(user *user.User) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	now := time.Now()
	exp := now.Add(config.JWTExpirationTime).Unix()

	claims["sub"] = user.ID.Hex()
	claims["iat"] = now.Unix()
	claims["exp"] = exp
	claims["jti"] = uuid.New().String()
	claims["username"] = user.Username
	claims["email"] = user.Email

	tokenString, err := token.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		return "", 0, err
	}

	return tokenString, exp, nil
}

func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must contain at least one uppercase letter, lowercase letter, number, and special character")
	}

	return nil
}
