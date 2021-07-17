package libravatar

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
)

func AvatarURL(email string) (imageURL string, err error) {
	return AvatarURLSize(email, SizeDefault)
}

func AvatarURLSize(email string, size ImageSize) (imageURL string, err error) {
	if size < SizeMinimum {
		size = 1
	} else if size > SizeMaximum {
		size = SizeMaximum
	}

	p := Parse(email)
	if p == nil {
		return "", fmt.Errorf("could not parse: %#v", email)
	}

	var host string

	host, err = lookup(p.Hostname)
	if err != nil {
		return "", fmt.Errorf("could not lookup host: %w", err)
	}

	imageURL = fmt.Sprintf("https://%s/avatar/%s", host, p.Hash)

	return imageURL, nil
}

func lookup(hostname string) (string, error) {
	_, addrs, err := net.LookupSRV("avatars-sec", "tcp", hostname)
	if err != nil {
		return "", fmt.Errorf("could not lookup hostname %#v: %w", hostname, err)
	}

	choice := choose(addrs)
	if choice == nil {
		return "", fmt.Errorf("could not lookup hostname %#v: %w", hostname, err)
	}

	// TODO: Sanitize.
	addr := fmt.Sprintf("%s:%d", choice.Target, choice.Port)

	return addr, nil
}

func choose(addrs []*net.SRV) *net.SRV {
	if len(addrs) == 0 {
		return nil
	}

	choices := []*net.SRV{addrs[0]}
	for i := 1; i < len(addrs); i++ {
		srv := addrs[i]
		diff := srv.Priority - choices[0].Priority
		if diff > 0 {
			choices = []*net.SRV{srv}

			continue
		}
		if diff < 0 {
			continue
		}

		choices = append(choices, srv)
	}

	totalWeight := int64(0)
	for _, srv := range choices {
		// Add `1` to weight to avoid zero weight.
		totalWeight += int64(srv.Weight) + 1
	}

	var (
		err  error
		bigN *big.Int
	)

	bigN, err = rand.Int(rand.Reader, big.NewInt(totalWeight))
	if err != nil {
		// If there's an error we ignore weights and just return the first
		// element.
		return choices[0]
	}

	var n int64 = bigN.Int64()

	for n > 0 && len(choices) > 1 {
		n -= int64(choices[0].Weight)
		choices = choices[1:]
	}

	return choices[0]
}
