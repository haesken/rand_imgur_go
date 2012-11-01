package main

import (
    "fmt"
    "math/rand"
    "time"
)

func randInt(min int, max int) int {
    return min + rand.Intn(max - min)
}

func randString (length int) string {
    var alpha = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
    buffer := make([]byte, length)
    for count := 0 ; count < length ; count++ {
        buffer[count] = alpha[rand.Intn(len(alpha)-1)]
    }
    return string(buffer)
}

func genURL() (string, string) {
    var base = "http://www.imgur.com/"
    imgurName := randString(5)
    imgurURL := base + imgurName + ".jpg"
    return imgurName, imgurURL
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    fmt.Println(randString(5))
    fmt.Println(genURL())
}
