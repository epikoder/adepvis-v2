package validator

import (
	"unicode"
)

func isEmail(str string) bool {
	return regexEmail.MatchString(str)
}

func isNumeric(str string) bool {
	return regexNumeric.MatchString(str)
}

func isAlphaNumeric(s string) bool {
	return regexAlphaNumeric.MatchString(s)
}

func isAlpha(s string) bool {
	return regexAlpha.MatchString(s)
}

func isAlphaSpaceNoDash(s string) bool {
	return regexAlphaSpaceNoDash.MatchString(s)
}

func isAlphaNumericCase(s string) (b bool) {
	var (
		hasNumber = false
		hasUpper  = false
		hasLower  = false
		hasSpace  = false
	)
	for _, char := range s {
		switch {
		case unicode.IsNumber(char):
			{
				hasNumber = true
			}
		case unicode.IsUpper(char):
			{
				hasUpper = true
			}
		case unicode.IsLower(char):
			{
				hasLower = true
			}
		case unicode.IsSpace(char):
			{
				hasSpace = true
			}
		}
	}

	return hasLower && hasNumber && hasUpper && !hasSpace
}

func isUUID(s string) bool {
	return regexUUID.MatchString(s)
}

func isPhone(s string) bool {
	return regexPhone.MatchString(s)
}

func isString(s string) bool {
	return regexString.MatchString(s)
}

func isEmpty(v interface{}) bool {
	return v == nil || v == ""
}
