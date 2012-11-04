package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/droundy/goopt"
)

const alphanum string = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// randInt generates a random int between a set minimum and maximum
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// randString generates a random string of a set length.
func randString(length int, charset string) string {
	buffer := make([]byte, length)

	for count := 0; count < length; count++ {
		buffer[count] = charset[randInt(0, len(charset))]
	}

	return string(buffer)
}

// genImgurURL generates a possible url for an image hosted at imgur.com
// It returns the "name" of the image (five randomized alphanumeric characters)
// and the respective url for that image.
func genImgurURL() (string, string) {
	base := "http://i.imgur.com/"
	imgurName := randString(5, alphanum)
	// Note: if the image hosted at imgur is not actually a jpg,
	// imgur will respond with the image anyway.
	imgurURL := base + imgurName + ".jpg"

	return imgurName, imgurURL
}

// hashImage takes and image takes a md5 hash of it.
// It returns a string containing the hex representation of the hash.
func hashImage(image []byte) string {
	hasher := md5.New()
	hasher.Write(image)

	return hex.EncodeToString(hasher.Sum(nil))
}

// getUrl fetches a url with a http GET.
// It returns if the contents and filetype of the the response.
func getUrl(url string) ([]byte, string, error) {
	resp, httpErr := http.Get(url)
	if httpErr != nil {
		log.Printf("http.Get -> %v", httpErr)
		log.Printf("http.Get -> Continuing...")
		return nil, "", httpErr
	}

	contents, fileErr := ioutil.ReadAll(resp.Body)
	if fileErr != nil {
		log.Printf("ioutil.ReadAll -> %v", fileErr)
		log.Printf("ioutil.ReadAll -> Continuing...")
		return nil, "", fileErr
	}

	filetype := strings.Split(resp.Header.Get("content-type"), "/")[1]
	resp.Body.Close()

	return contents, filetype, nil
}

// pathExists determines if a file/directory already exists.
// It returns a bool, and an error if it gets one.
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// writeFile writes the contents of a byte slice to the specified
// directory.
func writeFile(contents []byte, pathDirectory string, filename string) {
	pathDirectoryExists, err := pathExists(pathDirectory)
	if err != nil {
		log.Printf("writeFile -> %v", err)
		log.Printf("writeFile  -> Continuing...")
	}
	if pathDirectoryExists != true {
		os.MkdirAll(pathDirectory, 0777)
	}

	imagePath := path.Join(pathDirectory, filename)
	ioutil.WriteFile(imagePath, contents, 0666)
}

// findImages searches imgur for images. It requests a random image,
// but only writes it to disk if it is not the 404 gif.
func findImages(interval int, directory string, threadNum int) {
	for {
		imgurName, imgurURL := genImgurURL()
		image, filetype, err := getUrl(imgurURL)
		if err == nil {
			image_hash := hashImage(image)
			// Hash here is the 404 gif's hash.
			if image_hash != "d835884373f4d6c8f24742ceabe74946" {
				timestamp := time.Now().Format("2006-01-02 15-04-05")
				filename := imgurName + "." + filetype
				log.Printf("| Thread: %d | Found: %s", threadNum, filename)
				writeFile(image, directory, timestamp+" "+filename)
			} else {
				log.Printf("| Thread: %d | Found: 404 gif", threadNum)
			}
		}

		// Throttle connections to one per second per thread.
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

// main handles parsing command line arguments and spawning instances of
// findImages()
func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	var interval = goopt.Int([]string{"-i", "--interval"}, 2000,
		"Interval between requests. (Milliseconds)")
	var connections = goopt.Int([]string{"-c", "--connections"}, 4,
		"Number of simultanious connections.")
	var directory = goopt.String([]string{"-d", "--directory"}, "images",
		"Directory to save images to.")

	goopt.Description = func() string {
		return "Download random images from imgur"
	}
	goopt.Version = "0.0.1"
	goopt.Summary = "Random imgur downloader"
	goopt.Parse(nil)

	// Create requested number of connections.
	for threadNum := 1; threadNum < *connections; threadNum++ {
		go findImages(*interval, *directory, threadNum)
	}
	findImages(*interval, *directory, 0)
}
