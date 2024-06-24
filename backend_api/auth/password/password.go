package password

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hashedpassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
	return err
}
