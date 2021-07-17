package libravatar

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/c032/fork-go-resize"
)

// AvatarFunc returns an `io.ReadCloser` for the image file that corresponds to
// `hash`.
//
// If no image is found, `err` should be `ErrNotFound`.
type AvatarFunc func(hash string) (r io.ReadCloser, err error)

type Server struct {
	Avatar AvatarFunc
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if s.Avatar == nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if !strings.HasPrefix(req.URL.Path, "/avatar/") {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	imageHash := strings.TrimPrefix(req.URL.Path, "/avatar/")
	if !isSHA256(imageHash) {
		w.WriteHeader(http.StatusBadRequest)

		// TODO: Write error message to response.

		return
	}

	var (
		err error

		imageSizeStr string
		imageSize64  uint64

		imageReader io.ReadCloser
	)

	imageSize64, err = strconv.ParseUint(imageSizeStr, 10, 16)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		// TODO: Write error message to response.

		return
	}
	if imageSize64 < uint64(SizeMinimum) || imageSize64 > uint64(SizeMaximum) {
		w.WriteHeader(http.StatusBadRequest)

		// TODO: Write error message to response.

		return
	}

	imageSize := ImageSize(imageSize64)

	imageReader, err = s.Avatar(imageHash)
	if err != nil {
		if err == ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}
	defer imageReader.Close()

	var img image.Image
	img, _, err = image.Decode(imageReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	resizedImage := resize.Thumbnail(
		uint(imageSize),
		uint(imageSize),
		img,
		resize.Lanczos3,
	)

	w.Header().Set("Content-Type", "image/png")

	png.Encode(w, resizedImage)
}

func isSHA256(hash string) bool {
	if len(hash) != 64 {
		return false
	}

	for _, c := range hash {
		if c >= '0' && c <= '9' {
			continue
		}
		if c >= 'a' && c <= 'f' {
			continue
		}

		return false
	}

	return true
}
