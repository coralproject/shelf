package galleryfix

import "os"

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/ask/form/gallery/galleryfix/"
}
