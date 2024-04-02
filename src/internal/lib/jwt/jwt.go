package jwt

//
//func NewToken(user *models.User, key string, duration time.Duration) (string, error) {
//	token := jwt.New(jwt.SigningMethodHS256)
//
//	claims := token.Claims.(jwt.MapClaims)
//	claims["uid"] = user.Id
//	claims["role"] = user.Role
//	claims["exp"] = time.Now().Add(duration).Unix()
//
//	tokenString, err := token.SignedString([]byte(key))
//	if err != nil {
//		return "", errors.Wrap(err, "lib.jwt.NewToken error in sign")
//	}
//
//	return tokenString, nil
//}

//func ParseToken() {
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		// Don't forget to validate the alg is what you expect:
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
//		}
//
//		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
//		return hmacSampleSecret, nil
//	})
//
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		fmt.Println(claims["foo"], claims["nbf"])
//	} else {
//		fmt.Println(err)
//	}

//}
