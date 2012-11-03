package main

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "os"
    "path"
    "time"
    "strconv"
)


var alphanum = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// randInt generates a random int between a set minimum and maximum
func randInt(min int, max int) int {
    return min + rand.Intn(max - min)
}

// randString generates a random string of a set length.
func randString (length int, charset string) string {
    buffer := make([]byte, length)

    for count := 0 ; count < length ; count++ {
        buffer[count] = charset[randInt(0, len(charset))]
    }

    return string(buffer)
}

// genImgurURL generates a possible url for an image hosted at imgur.com
// It returns the "name" of the image (five randomized alphanumeric characters)
// and the respective url for that image.
func genImgurURL() (string, string) {
    base := "http://www.imgur.com/"
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
// It returns if the contents of the the response.
func getUrl(url string) []byte {
    resp, err := http.Get(url)
    if err != nil {
        log.Fatalf("http.Get -> %v", err)
    }

    contents, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("ioutil.ReadAll -> %v", err)
    }

    resp.Body.Close()

    return contents
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
func writeFile(contents []byte, pathDirectory string, filename string) int {
    pathDirectoryExists, err := pathExists(pathDirectory)
    if err != nil {
        return 1
    }
    if pathDirectoryExists != true {
        os.MkdirAll(pathDirectory, 0777)
    }

    imagePath := path.Join(pathDirectory, filename)
    ioutil.WriteFile(imagePath, contents, 0666)
    return 0
}


// findImages searches imgur for images. It requests a random image,
// but only writes it to disk if it is not the 404 gif.
func findImages() {
    for {
        imgurName, imgurURL := genImgurURL()
        image := getUrl(imgurURL)
        image_hash := hashImage(image)
        timestamp := strconv.FormatInt(time.Now().Unix(), 10)

        // Hash here is the 404 gif's hash.
        if image_hash != "d835884373f4d6c8f24742ceabe74946" {
            fmt.Println("Found image!")
            fmt.Println(imgurURL)
            filename := timestamp + " " + imgurName + ".jpg"
            pathDirectory := "images"
            writeFile(image, pathDirectory, filename)
        }

        // Throttle connects to one per second per thread.
        time.Sleep(1 * time.Second)
    }
}


func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    go findImages()
    go findImages()
    go findImages()
    findImages()
}
