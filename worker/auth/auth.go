package auth

var authFilterMap = make(map[string]struct{})

func SetMap(a map[string]struct{}) {
	authFilterMap = a
}

func Check(token string) bool {
	if _, ok := authFilterMap[token]; ok {
		return true
	} else {
		return false
	}
}
