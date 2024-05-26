package imagery

import (
	"fmt"
	"image/jpeg"
	"image/png"

	// "image/webp"

	// "io"
	"net/http"
	"os"

	// "image"
	// "image/color"
	// "log"
	// "os"
	// webpbin "github.com/CapsLock-Studio/go-webpbin"
	// "path/filepath"

	httperrors "github.com/myrachanto/custom-http-error"
	"github.com/nfnt/resize"
)

// Imagery ...
var (
	Imageryrepository imageryrepository = imageryrepository{}
)

type imageryrepository struct{}

// Imagetype ...
func (i imageryrepository) Imagetype(f, filename string, height, width int) {

	// maximize CPU usage for maximum performance
	// runtime.GOMAXPROCS(runtime.NumCPU())

	// open the uploaded file
	// f := "./img.png"
	file, err := os.Open(f)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = file.Read(buff)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filetype := http.DetectContentType(buff)

	fmt.Println(filetype)
	// Define the new width and height while maintaining the aspect ratio
	newWidth := uint(width) // Desired width
	newHeight := uint(0)    // Automatically calculate to maintain aspect ratio

	switch filetype {
	case "image/jpeg", "image/jpg":
		ResizeJPG(f, filename, newHeight, newWidth)

	case "image/png":
		ResizePng(f, filename, newHeight, newWidth)

	default:
		fmt.Println("Wrong file format")
	}

}

// ResizePng ...
func ResizePng(f, filename string, height, width uint) (*os.File, httperrors.HttpErr) {

	file, err := os.Open(f)
	if err != nil {
		return nil, httperrors.NewNotFoundError("error opening the file")
	}
	// decode png into image.Image
	img, err := png.Decode(file)
	if err != nil {
		return nil, httperrors.NewNotFoundError("error decoding png")
	}

	// resize to width 60 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(width, height, img, resize.Lanczos3)
	// return out, nil
	out, err := os.Create(filename)
	if err != nil {
		return nil, httperrors.NewNotFoundError("error decoding png")
	}
	defer out.Close()

	// write new image to file
	f1 := png.Encode(out, m)
	if f1 != nil {
		return nil, httperrors.NewNotFoundError("error encoding png")
	}
	return out, nil
}

// ResizeJPG ...
func ResizeJPG(f, filename string, height, width uint) (*os.File, httperrors.HttpErr) {
	file, err := os.Open(f)
	if err != nil {
		return nil, httperrors.NewNotFoundError("error opening the file")
	}
	// decode png into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, httperrors.NewNotFoundError("error decoding jpeg")
	}

	// resize to width 60 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(width, height, img, resize.Lanczos3)
	out, err := os.Create(filename)
	if err != nil {
		return nil, httperrors.NewNotFoundError("error decoding jpeg")
	}
	defer out.Close()

	// write new image to file
	f1 := jpeg.Encode(out, m, nil)
	if f1 != nil {
		return nil, httperrors.NewNotFoundError("error encoding jpeg")
	}
	return out, nil
}

// //ResizeJPG ...
// func Resizewebp(filename string) (*os.File, *httperrors.HttpError) {
// 	err := webpbin.NewCWebP().
// 		Quality(80).
// 		InputFile(filename).
// 		OutputFile("image.webp").
// 		Run()
// 	return out, nil
// }
