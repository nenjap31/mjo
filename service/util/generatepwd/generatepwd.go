package generatepwd

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s", password)), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(result), nil
}