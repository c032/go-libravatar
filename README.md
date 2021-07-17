# go-libravatar

## Explicitly not supported

* **MD5 hashes.** This library always uses and expects SHA256.
* **Default URLs.** This library always forces a `d=404` parameter on
  all requests sent by the client, and received by the server.

## Server example

* `"/usr/local/var/libravatar/images/${HASH}"` is an existing image file.
* `curl "http://localhost:3000/avatar/${HASH}"` responds with the image
  file.

```go
package main

import (
	"net/http"
	"log"
	"path"
	"fmt"

	"github.com/c032/go-libravatar"
)

// ImagesDirectory is the path to the directory containing all avatar
// images.
const ImagesDirectory = "/usr/local/var/libravatar/images"

// Avatar returns an `io.ReadCloser` for the image file that corresponds
// to `hash`.
func Avatar(hash string) (io.ReadCloser, error) {
	imagePath := path.Join(ImagesDirectory, hash)

	f, err := os.Open(imagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, libravatar.ErrNotFound
		}

		return nil, fmt.Errorf("could not open image: %w", err)
	}

	return f, nil
}

func main() {
	const addr = ":3000"

	s := &libravatar.Server{
		Avatar: Avatar,
	}

	err := http.ListenAndServe(addr, s)
	if err != nil {
		log.Fatal(err)
	}
}
```

## License

Apache 2.0
