package google

import (
	"net/http"

	"google.golang.org/api/oauth2/v2"
)

func VerifyToken(idToken string, httpClient *http.Client) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(httpClient)
	if err != nil {
		return nil, err
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}
