package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	today := time.Now().UTC()
	year, _, _ := today.Date()
	finalDayOfYear := time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)
	daysLeft := int(math.Floor(time.Until(finalDayOfYear).Abs().Hours() / 24))

	result := fmt.Sprintf("Days left in %d: %d", year, daysLeft)

	bskyUsername := os.Getenv("BSKY_USERNAME")
	bskyPassword := os.Getenv("BSKY_APP_PASSWORD")
	sentryDSN := os.Getenv("SENTRY_DSN")

	if sentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			EnableTracing:    true,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		// Flush buffered events before the program terminates.
		defer sentry.Flush(2 * time.Second)
	}

	log.Println(result)
	if bskyUsername != "" && bskyPassword != "" {
		checkinId := sentry.CaptureCheckIn(
			&sentry.CheckIn{
				MonitorSlug: "daysleft-monitor",
				Status:      sentry.CheckInStatusInProgress,
			},
			nil,
		)
		resp, err := PostToBluesky(result, bskyUsername, bskyPassword)
		if err != nil {
			sentry.CaptureCheckIn(
				&sentry.CheckIn{
					ID:          *checkinId,
					MonitorSlug: "daysleft-monitor",
					Status:      sentry.CheckInStatusError,
				},
				nil,
			)
			log.Printf("Failed to post to Bluesky: %v", err)
		} else {
			sentry.CaptureCheckIn(
				&sentry.CheckIn{
					ID:          *checkinId,
					MonitorSlug: "daysleft-monitor",
					Status:      sentry.CheckInStatusOK,
				},
				nil,
			)
			log.Printf("Posted to Bluesky successfully: %v", resp.Uri)
		}
	}

}
