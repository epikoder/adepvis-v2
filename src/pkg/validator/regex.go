package validator

import (
	"regexp"
)

const (
	// alpha represents regular expression for alpha characters
	alpha string = "^[a-zA-Z]+$"

	alphaSpaceNoDash string = "^[a-zA-Z ]+$"
	// alphaDash represents regular expression for alpha characters with underscore and dash
	alphaDash string = "^[a-zA-Z0-9_-]+$"
	// alphaSpace represents regular expression for alpha characters with underscore, space and dash
	alphaSpace string = "^[-a-zA-Z0-9_ ]+$"
	// alphaNumeric represents regular expression for alpha numeric characters
	alphaNumeric string = "^[a-zA-Z0-9]+$"
	// creditCard represents regular expression for credit cards like (Visa, MasterCard, American Express, Diners Club, Discover, and JCB cards). Ref: https://stackoverflow.com/questions/9315647/regex-credit-card-number-tests
	creditCard string = "^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$"
	// coordinate represents latitude and longitude regular expression
	coordinate string = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?),\\s*[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$" // Ref: https://stackoverflow.com/questions/3518504/regular-expression-for-matching-latitude-longitude-coordinates
	// cSSColor represents css valid color code with hex, rgb, rgba, hsl, hsla etc. Ref: http://www.regexpal.com/97509
	cSSColor string = "^(#([\\da-f]{3}){1,2}|(rgb|hsl)a\\((\\d{1,3}%?,\\s?){3}(1|0?\\.\\d+)\\)|(rgb|hsl)\\(\\d{1,3}%?(,\\s?\\d{1,3}%?){2}\\))$"
	// date represents regular expression for valid date like: yyyy-mm-dd
	date string = "^(((19|20)([2468][048]|[13579][26]|0[48])|2000)[/-]02[/-]29|((19|20)[0-9]{2}[/-](0[469]|11)[/-](0[1-9]|[12][0-9]|30)|(19|20)[0-9]{2}[/-](0[13578]|1[02])[/-](0[1-9]|[12][0-9]|3[01])|(19|20)[0-9]{2}[/-]02[/-](0[1-9]|1[0-9]|2[0-8])))$"
	// dateDDMMYY represents regular expression for valid date of format dd/mm/yyyy , dd-mm-yyyy etc.Ref: http://regexr.com/346hf
	dateDDMMYY string = "^(0?[1-9]|[12][0-9]|3[01])[\\/\\-](0?[1-9]|1[012])[\\/\\-]\\d{4}$"
	// digits represents regular epxression for validating digits
	digits string = "^[+-]?([0-9]*\\.?[0-9]+|[0-9]+\\.?[0-9]*)([eE][+-]?[0-9]+)?$"
	// email represents regular expression for email
	email string = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$"
	// float represents regular expression for finding float number
	float string = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	// iP represents regular expression for ip address
	iP string = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"
	// iPV4 represents regular expression for ip address version 4
	iPV4 string = "^([0-9]{1,3}\\.){3}[0-9]{1,3}(\\/([0-9]|[1-2][0-9]|3[0-2]))?$"
	// iPV6 represents regular expression for ip address version 6
	iPV6 string = `^s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:)))(%.+)?s*(\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?$`
	// lLatitude represents latitude regular expression
	lLatitude string = "^(\\+|-)?(?:90(?:(?:\\.0{1,6})?)|(?:[0-9]|[1-8][0-9])(?:(?:\\.[0-9]{1,6})?))$"
	// longitude represents longitude regular expression
	longitude string = "^(\\+|-)?(?:180(?:(?:\\.0{1,6})?)|(?:[0-9]|[1-9][0-9]|1[0-7][0-9])(?:(?:\\.[0-9]{1,6})?))$"
	// macAddress represents regular expression for mac address
	macAddress string = "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
	// numeric represents regular expression for numeric
	numeric string = "^-?[0-9]+$"
	// uRL represents regular expression for url
	uRL string = "^(?:http(s)?:\\/\\/)?[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:/?#[\\]@!\\$&'\\(\\)\\*\\+,;=.]+$" // Ref: https://stackoverflow.com/questions/136505/searching-for-uuids-in-text-with-regex
	// uUID represents regular expression for uUID
	uUID string = "^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}$"
	// uUID3 represents regular expression for UUID version 3
	uUID3 string = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	// uUID4 represents regular expression for UUID version 4
	uUID4 string = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	// uUID5 represents regular expression for UUID version 5
	uUID5   string = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	pHone   string = "^(0|234)([0-9]{10})$"
	string_ string = "^[-a-zA-Z0-9_ -/.\\,]+$"
)

var (
	regexAlpha            = regexp.MustCompile(alpha)
	regexAlphaSpaceNoDash = regexp.MustCompile(alphaSpaceNoDash)
	regexAlphaDash        = regexp.MustCompile(alphaDash)
	regexAlphaSpace       = regexp.MustCompile(alphaSpace)
	regexAlphaNumeric     = regexp.MustCompile(alphaNumeric)
	regexCreditCard       = regexp.MustCompile(creditCard)
	regexCoordinate       = regexp.MustCompile(coordinate)
	regexCSSColor         = regexp.MustCompile(cSSColor)
	regexDate             = regexp.MustCompile(date)
	regexDateDDMMYY       = regexp.MustCompile(dateDDMMYY)
	regexDigits           = regexp.MustCompile(digits)
	regexEmail            = regexp.MustCompile(email)
	regexFloat            = regexp.MustCompile(float)
	regexMacAddress       = regexp.MustCompile(macAddress)
	regexNumeric          = regexp.MustCompile(numeric)
	regexLatitude         = regexp.MustCompile(lLatitude)
	regexLongitude        = regexp.MustCompile(longitude)
	regexIP               = regexp.MustCompile(iP)
	regexIPV4             = regexp.MustCompile(iPV4)
	regexIPV6             = regexp.MustCompile(iPV6)
	regexURL              = regexp.MustCompile(uRL)
	regexUUID             = regexp.MustCompile(uUID)
	regexUUID3            = regexp.MustCompile(uUID3)
	regexUUID4            = regexp.MustCompile(uUID4)
	regexUUID5            = regexp.MustCompile(uUID5)
	regexPhone            = regexp.MustCompile(pHone)
	regexString           = regexp.MustCompile(string_)
)
