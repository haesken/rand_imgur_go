package main

import (
    "fmt"
    "math/rand"
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

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    fmt.Println(genImgurURL())
}
