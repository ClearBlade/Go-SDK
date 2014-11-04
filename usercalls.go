package GoSDK

const (
	_USER_HEADER_KEY = "ClearBlade-UserToken"
	_USER_PREAMBLE   = "api/v/1"
)

func (u *UserClient) creds() ([][]string, error) {
	if u.UserToken != "" {
		return [][]string{
			[]string{
				_USER_HEADER_KEY,
				u.UserToken,
			},
		}, nil
	} else if u.SystemSecret != "" && u.SystemKey != "" {
		return [][]string{
			[]string{
				_HEADER_SECRET_KEY,
				u.SystemSecret,
			},
			[]string{
				_HEADER_KEY_KEY,
				u.SystemKey,
			},
		}, nil
	} else {
		return [][]string{}, errors.New("No SystemSecret/SystemKey combo, or UserToken found")
	}
}

func (u *UserClient) preamble() string {
	return _USER_PREAMBLE
}

func (u *UserClient) setToken(t string) {
	u.UserToken = t
}
func (u *UserToken) getToken() string {
	return u.UserToken
}
