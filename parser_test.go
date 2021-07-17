package libravatar_test

import (
	"testing"

	"github.com/c032/go-libravatar"
)

func TestParse(t *testing.T) {
	tt := map[string]*libravatar.Parsed{
		"HTTP://UserName:Password@EXAMPLE.COM/ID/User": &libravatar.Parsed{
			Hash:     "8cf939a9ba0c432e1f5e7b73e649e656b7238f614cb1967165b7c56ecd28dae4",
			Hostname: "example.com",
		},
		"http://UserName:Password@example.com/ID/User": &libravatar.Parsed{
			Hash:     "8cf939a9ba0c432e1f5e7b73e649e656b7238f614cb1967165b7c56ecd28dae4",
			Hostname: "example.com",
		},
		"user@Example.com": &libravatar.Parsed{
			Hash:     "b4c9a289323b21a01c3e940f150eb9b8c542587f1abfd8f0e1cc1ffc5e475514",
			Hostname: "example.com",
		},
		"user@example.com": &libravatar.Parsed{
			Hash:     "b4c9a289323b21a01c3e940f150eb9b8c542587f1abfd8f0e1cc1ffc5e475514",
			Hostname: "example.com",
		},
	}

	for input, expected := range tt {
		p := libravatar.Parse(input)
		if p == nil {
			t.Errorf("libravatar.Parse(%q) = nil; want non-nil", input)

			continue
		}

		if got, want := p.Hash, expected.Hash; got != want {
			t.Errorf("libravatar.Parse(%q).Hash = %q; want %q", input, got, want)

			continue
		}

		if got, want := p.Hostname, expected.Hostname; got != want {
			t.Errorf("libravatar.Parse(%q).Hostname = %q; want %q", input, got, want)

			continue
		}
	}
}
