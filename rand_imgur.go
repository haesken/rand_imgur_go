package main

import (
    "fmt"
    "math/rand"
    "time"
)


var alphanum = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"


func randInt(min int, max int) int {
    return min + rand.Intn(max - min)
}

func randString (length int, charset string) string {
    buffer := make([]byte, length)

    for count := 0 ; count < length ; count++ {
        buffer[count] = charset[randInt(0, len(charset))]
    }

    return string(buffer)
}

func genURL() (string, string) {
    base := "http://www.imgur.com/"
    imgurName := randString(5, alphanum)
    imgurURL := base + imgurName + ".jpg"

    return imgurName, imgurURL
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    fmt.Println(genURL())
}
