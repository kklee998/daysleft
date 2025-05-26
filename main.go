package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

// Daysleft calculates the number of days left until the end the current year.
func Daysleft() int {
	t := time.Now().UTC()
	y, _, _ := t.Date()
	finalDayOfYear := time.Date(y, time.December, 31, 23, 59, 59, 0, time.UTC)

	return int(math.Floor(time.Until(finalDayOfYear).Abs().Hours() / 24))

}

func main() {
	today := time.Now().UTC()
	year, _, _ := today.Date()
	d := Daysleft()

	result := fmt.Sprintf("Days left in %d: %d", year, d)

	bskyUsername := os.Getenv("BSKY_USERNAME")
	bskyPassword := os.Getenv("BSKY_PASSWORD")

	if bskyUsername != "" && bskyPassword != "" {
		PostToBluesky(result, bskyUsername, bskyPassword)
	} else {
		fmt.Println(result)
	}

}
