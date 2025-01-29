package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coding-for-fun-org/gcli/pkg/gcal"
	"github.com/coding-for-fun-org/gcli/pkg/util"
)

var maxOutputLength = 20

func main() {
	c := gcal.Calendar{
		Id: "primary",
	}
	c.Initialize()

	evts, err := c.GetTodayEvents(true)
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
			if item.Start == nil || item.End == nil || item.Start.DateTime == "" ||
				item.End.DateTime == "" {
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
			if item.Start == nil || item.End == nil || item.Start.DateTime == "" ||
				item.End.DateTime == "" ||
				c.GetWorkingHoursProperty(item) != "" {
				continue
			}

			t := time.Now().Local()

			if gap, err := util.CalculateTimeGap(t.Format(time.RFC3339), item.Start.DateTime); err == nil &&
				gap > 0 {
				output = fmt.Sprintf(
					"[%v] in %.0fmin\n",
					util.TruncateWithSuffix(item.Summary, maxOutputLength),
					gap.Minutes(),
				)
				break
			}
		}
		if output != "" {
			fmt.Print(output)
		} else {
			fmt.Print("N/A")
		}

		break

	case "in-progress":
		var output string
		for _, item := range evts.Items {
			if item.Start == nil || item.End == nil || item.Start.DateTime == "" ||
				item.End.DateTime == "" ||
				c.GetWorkingHoursProperty(item) != "" {
				continue
			}

			t := time.Now().Local()
			st, err := time.Parse(time.RFC3339, item.Start.DateTime)
			if err != nil {
				log.Fatalf("Unable to parse start time: %v", err)
			}
			et, err := time.Parse(time.RFC3339, item.End.DateTime)
			if err != nil {
				log.Fatalf("Unable to parse end time: %v", err)
			}
			startGap, errStartGap := util.CalculateTimeGap(t.Format(time.RFC3339), item.Start.DateTime)
			endGap, errEndGap := util.CalculateTimeGap(t.Format(time.RFC3339), item.End.DateTime)

			if errStartGap == nil && startGap < 0 && errEndGap == nil && endGap > 0 {
				output = fmt.Sprintf(
					"[%v] (%v-%v)\n",
					util.TruncateWithSuffix(item.Summary, maxOutputLength),
					st.Local().Format("15:04"),
					et.Local().Format("15:04"),
				)
				break
			}
		}
		if output != "" {
			fmt.Print(output)
		} else {
			fmt.Print("N/A")
		}

		break
	}
}
