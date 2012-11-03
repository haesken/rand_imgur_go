package main

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "time"
)


var alphanum = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"


// randInt generates a random int between a set minimum and maximum
func randInt(min int, max int) int {
    return min + rand.Intn(max - min)
}

// randString generates a random string of a set length.
// The length argument determines the length of the string to generate,
// and charset argument serves as the list of possible characters
// to choose from.
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
// Note: if the image hosted at imgur is not actually a jpg,
// imgur will respond with the image anyway.
func genImgurURL() (string, string) {
    base := "http://www.imgur.com/"
    imgurName := randString(5, alphanum)
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

func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    imgurName, imgurURL := genImgurURL()
    fmt.Println(imgurURL)
    image := getUrl(imgurURL)
    image_hash := hashImage(image)

    if image_hash != "d835884373f4d6c8f24742ceabe74946" {
        filename := imgurName + ".jpg"
        ioutil.WriteFile(filename, image, 0666)
    } else {
        fmt.Println("Found 404 gif!")
    }
}
