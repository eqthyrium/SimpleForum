package cookie

import "net/http"

func SetTokenCookie(w http.ResponseWriter, token string) {

	signedToken, err := CreateSignedToken()
	if err != nil {

	}
}
