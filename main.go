package main

import (
	"fmt"
	"google_calendar_cli/pkg/google_oauth"
	"google_calendar_cli/pkg/utils"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
)

func main() {
	o := google_oauth.OAuth{}

	srv, err := o.GetCalendarService()
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	calendarId := "primary"

	t := utils.StartOfDayTime()
	events, err := srv.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	var pendingEvent *calendar.Event
	for _, item := range events.Items {
		if item.Start.DateTime == item.End.DateTime && (item.Start.DateTime != "") {
			pendingEvent = item
		}

		if pendingEvent != nil {
			break
		}
	}

	currentTime := time.Now().Format(time.RFC3339)
	if pendingEvent != nil {
		duration, err := utils.CalculateDuration(currentTime, pendingEvent.End.DateTime)
		if err != nil {
			log.Fatalf("Unable to calculate duration. %v\n", err)
		}

		pendingEvent.Summary = fmt.Sprintf("__Work (%.3f hrs)__", duration.Hours())
		pendingEvent.End.DateTime = currentTime

		event, err := srv.Events.Update(calendarId, pendingEvent.Id, pendingEvent).Do()
		if err != nil {
			log.Fatalf("Unable to update event. %v\n", err)
		}
		fmt.Printf("Event updated: %s\n", event.HtmlLink)
	} else {
		event := &calendar.Event{
			Summary: "__Working__",
			Start: &calendar.EventDateTime{
				DateTime: currentTime,
			},
			End: &calendar.EventDateTime{
				DateTime: currentTime,
			},
			Visibility:      "public",
			Transparency:    "transparent",
			GuestsCanModify: false,
			ColorId:         "8",
		}
		boolFalse := false
		event.GuestsCanSeeOtherGuests = &boolFalse
		event.GuestsCanInviteOthers = &boolFalse

		event, err = srv.Events.Insert(calendarId, event).Do()
		if err != nil {
			log.Fatalf("Unable to create event. %v\n", err)
		}
		fmt.Printf("Event created: %s\n", event.HtmlLink)
	}

	// fmt.Println("Upcoming events:")
	// if len(events.Items) == 0 {
	// 	fmt.Println("No upcoming events found.")
	// } else {
	// 	for _, item := range events.Items {
	// 		date := item.Start.DateTime
	// 		if date == "" {
	// 			date = item.Start.Date
	// 		}
	// 		fmt.Printf("%v (%v)\n", item.Summary, date)
	// 	}
	// }
}
