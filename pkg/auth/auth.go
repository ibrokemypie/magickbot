package auth

func getFediAccessToken(instanceURL string) string {
	app := getApp(instanceURL)
	token := authorizeUser(instanceURL, app)
	accessToken := oauthToken(instanceURL, app, token)

	return accessToken
}

// Authorize - start oauth flow
func Authorize() (string, string) {
	instanceURL := getFediInstance()
	accessToken := getFediAccessToken(instanceURL)

	return instanceURL, accessToken
}
