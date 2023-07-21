package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Make(i interface{}) (r string, err error) {
	s_ := fmt.Sprintf("%v", i)
	b, err := bcrypt.GenerateFromPassword([]byte(s_), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Check(h string, s interface{}) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(fmt.Sprintf("%v", s)))
	return err == nil
}
