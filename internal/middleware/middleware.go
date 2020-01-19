package middleware

type contextKey string

func (k contextKey) String() string {
	return "shorty-api/middleware context value " + string(k)
}
