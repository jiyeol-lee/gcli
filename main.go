package main

import (
	"google_app_cli/pkg/gcal"
	"log"
)

func main() {
	c := gcal.Calendar{
		Id: "primary",
	}
	c.Initialize()

	evts, err := c.GetTodayEvents(false)
	if err != nil {
		log.Fatalf("Unable to retrieve today's events: %v", err)
	}

	for _, item := range evts.Items {
		log.Printf("%v (%v - %v)\n", item.Summary, item.Start.DateTime, item.End.DateTime)
		log.Printf("Working Hours : %s\n", c.GetWorkingHoursProperty(item))
	}

	// evt, err := c.AddTotalWorkingEvent()
	// if err != nil {
	// 	log.Fatalf("Unable to add total working event: %v", err)
	// }
	//
	// log.Printf("Total Working Event: %v\n", evt.HtmlLink)

	totalWorkingEvent, err := c.GetTodayTotalWorkingEvent(evts)
	if err != nil {
		log.Fatalf("Unable to get today's total working event: %v", err)
	}
	evt, err := c.UpdateTotalWorkingEvent(totalWorkingEvent, evts)
	if err != nil {
		log.Fatalf("Unable to update total working event: %v", err)
	}
	log.Printf("Total Working Event: %v\n", evt.HtmlLink)

	// _, err = c.UpsertPendingEvent()
	// if err != nil {
	// 	log.Fatalf("Unable to upsert pending event: %v", err)
	// }
}
