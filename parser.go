package libravatar

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"strings"
)

type Parsed struct {
	Hash     string
	Hostname string
}

func Parse(input string) *Parsed {
	p := &Parsed{}

	var (
		err error
		u   *url.URL
	)

	u, err = url.Parse(input)
	if err != nil {
		return nil
	}
	if u.Scheme == "" {
		u, err = url.Parse("mailto://" + input)
		if err != nil {
			return nil
		}
	}

	u.Host = strings.ToLower(u.Host)

	p.Hostname = u.Hostname()

	var rawHash [32]byte
	if u.Scheme == "mailto" {
		// Email.

		cleanEmail := fmt.Sprintf("%s@%s", u.User.String(), p.Hostname)
		rawHash = sha256.Sum256([]byte(cleanEmail))
	} else {
		// URL.

		u.Scheme = strings.ToLower(u.Scheme)
		rawHash = sha256.Sum256([]byte(u.String()))
	}

	p.Hash = fmt.Sprintf("%x", rawHash)

	return p
}
