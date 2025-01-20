package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coding-for-fun-org/gcli/pkg/gcal"
	"github.com/coding-for-fun-org/gcli/pkg/util"
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

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		log.Fatalf("No command provided")
	}

	switch argsWithoutProg[0] {
	case "list":
		for _, item := range evts.Items {
			if item.Start.DateTime == "" || item.End.DateTime == "" {
				continue
			}

			tStart, err := time.Parse(time.RFC3339, item.Start.DateTime)
			if err != nil {
				log.Fatalf("Unable to parse start time: %v", err)
			}

			tEnd, err := time.Parse(time.RFC3339, item.End.DateTime)
			if err != nil {
				log.Fatalf("Unable to parse end time: %v", err)
			}

			fmt.Printf(
				"%v (%v - %v)\n",
				item.Summary,
				fmt.Sprintf("%02d:%02d", tStart.Local().Hour(), tStart.Local().Minute()),
				fmt.Sprintf("%02d:%02d", tEnd.Local().Hour(), tEnd.Local().Minute()),
			)
		}

		break

	case "soon":
		var output string
		for _, item := range evts.Items {
			if item.Start.DateTime == "" || item.End.DateTime == "" ||
				c.GetWorkingHoursProperty(item) != "" {
				continue
			}

			t := time.Now().Local()

			if gap, err := util.CalculateTimeGap(t.Format(time.RFC3339), item.Start.DateTime); err == nil &&
				gap > 0 {
				output = fmt.Sprintf(
					"[%v] %.0f min left\n",
					util.TruncateWithSuffix(item.Summary, 7),
					gap.Minutes(),
				)
				break
			}
		}
		if output != "" {
			fmt.Print(output)
		} else {
			fmt.Print("No event soon\n")
		}

		break
	}

	// evt, err := c.AddTotalWorkingEvent()
	// if err != nil {
	// 	log.Fatalf("Unable to add total working event: %v", err)
	// }
	//
	// log.Printf("Total Working Event: %v\n", evt.HtmlLink)

	// totalWorkingEvent, err := c.GetTodayTotalWorkingEvent(evts)
	// if err != nil {
	// 	log.Fatalf("Unable to get today's total working event: %v", err)
	// }
	// evt, err := c.UpdateTotalWorkingEvent(totalWorkingEvent, evts)
	// if err != nil {
	// 	log.Fatalf("Unable to update total working event: %v", err)
	// }
	// log.Printf("Total Working Event: %v\n", evt.HtmlLink)

	// _, err = c.UpsertPendingEvent()
	// if err != nil {
	// 	log.Fatalf("Unable to upsert pending event: %v", err)
	// }
}
