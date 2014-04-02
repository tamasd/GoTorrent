package magnet

import "net/url"

const (
	Prefix = "urn:btih:"
)

type Magnet struct {
	InfoHash string
	Name     string
	Tracker  string
}

func Parse(s string) (*Magnet, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	m := new(Magnet)
	m.InfoHash = getSuffix(getOneParameter(u, "xt"), Prefix)
	m.Name = getOneParameter(u, "dn")
	m.Tracker = getOneParameter(u, "tr")

	return m, nil
}

func getSuffix(s, prefix string) string {
	if len(s) >= len(prefix) && s[0:len(prefix)] == prefix {
		return s[len(prefix):]
	}

	return ""
}

func getOneParameter(u *url.URL, param string) string {
	arr, ok := u.Query()[param]
	if ok && len(arr) > 0 {
		return arr[0]
	}

	return ""
}
