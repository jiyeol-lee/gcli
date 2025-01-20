package gcal

import (
	"context"
	"fmt"
	"google_app_cli/pkg/goauth"
	"google_app_cli/pkg/utils"
	"log"
	"regexp"
	"strconv"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	workingHoursKey      string = "WORKING_HOURS"
	totalWorkingHoursKey string = "TOTAL_WORKING_HOURS"
)

type Calendar struct {
	Id      string
	Service *calendar.Service
}

func (c *Calendar) Initialize() {
	if c.Id == "" {
		log.Fatalf("Calendar ID is required")
	}

	o := goauth.OAuth{}

	err := o.SetClient(calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("Unable to set client: %v", err)
	}

	ctx := context.Background()

	svc, err := calendar.NewService(ctx, option.WithHTTPClient(o.Client))
	if err != nil {
		log.Fatalf("Unable to create service: %v", err)
	}

	c.Service = svc
}

func (c *Calendar) GetTodayEvents(onlySingleEvent bool) (*calendar.Events, error) {
	tmin := util.StartOfDayTime()
	tmax := util.EndOfDayTime()
	evtListCall := c.Service.Events.List(c.Id).ShowDeleted(false).
		SingleEvents(onlySingleEvent).TimeMin(tmin).TimeMax(tmax)
	var evts *calendar.Events
	var err error
	if onlySingleEvent {
		evts, err = evtListCall.OrderBy("startTime").Do()
	} else {
		evts, err = evtListCall.Do()
	}
	if err != nil {
		return nil, err
	}

	return evts, nil
}

func (_ *Calendar) GetWorkingHoursProperty(event *calendar.Event) string {
	if event.ExtendedProperties == nil {
		return ""
	}

	return event.ExtendedProperties.Private[workingHoursKey]
}

func (_ *Calendar) setWorkingHoursProperty(event *calendar.Event, workingHours float64) {
	if event.ExtendedProperties == nil {
		event.ExtendedProperties = &calendar.EventExtendedProperties{
			Private: map[string]string{},
		}
	}
	event.ExtendedProperties.Private[workingHoursKey] = fmt.Sprintf(
		"%.3f",
		workingHours,
	)
}

func (_ *Calendar) GetTotalWorkingHoursProperty(event *calendar.Event) string {
	if event.ExtendedProperties == nil {
		return ""
	}

	return event.ExtendedProperties.Private[totalWorkingHoursKey]
}

func (_ *Calendar) setTotalWorkingHoursProperty(event *calendar.Event, totalWorkingHours float64) {
	if event.ExtendedProperties == nil {
		event.ExtendedProperties = &calendar.EventExtendedProperties{
			Private: map[string]string{},
		}
	}
	event.ExtendedProperties.Private[totalWorkingHoursKey] = fmt.Sprintf(
		"%.3f",
		totalWorkingHours,
	)
}

func (c *Calendar) GetTodayPendingEvent(events *calendar.Events) (*calendar.Event, error) {
	var pendingEvent *calendar.Event
	for _, item := range events.Items {
		if c.GetWorkingHoursProperty(item) == "0.000" {
			pendingEvent = item
			break
		}
	}

	return pendingEvent, nil
}

func (c *Calendar) GetTodayTotalWorkingEvent(events *calendar.Events) (*calendar.Event, error) {
	var totalWorkingHoursEvent *calendar.Event
	for _, item := range events.Items {
		if matched, err := regexp.MatchString("^\\d.\\d{3}$", c.GetTotalWorkingHoursProperty(item)); matched == true &&
			err == nil {
			totalWorkingHoursEvent = item
			break
		}
	}

	return totalWorkingHoursEvent, nil
}

func (c *Calendar) AddPendingEvent() (*calendar.Event, error) {
	currentTime := time.Now().Format(time.RFC3339)
	event := &calendar.Event{
		Summary: "Working",
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
	c.setWorkingHoursProperty(event, 0)

	evt, err := c.Service.Events.Insert(c.Id, event).Do()
	if err != nil {
		return nil, err
	}

	return evt, nil
}

func (c *Calendar) UpdatePendingEvent(event *calendar.Event) (*calendar.Event, error) {
	if event == nil {
		return nil, fmt.Errorf("event is nil")
	}

	currentTime := time.Now().Format(time.RFC3339)
	duration, err := util.CalculateTimeGap(event.Start.DateTime, currentTime)
	if err != nil {
		return nil, err
	}

	event.Summary = fmt.Sprintf("Work (%.3f hrs)", duration.Hours())
	event.End.DateTime = currentTime
	c.setWorkingHoursProperty(event, duration.Hours())

	evt, err := c.Service.Events.Update(c.Id, event.Id, event).Do()
	if err != nil {
		return nil, err
	}

	return evt, nil
}

func (c *Calendar) AddTotalWorkingEvent() (*calendar.Event, error) {
	currentTime := time.Now()
	event := &calendar.Event{
		Summary: "Total Work",
		Start: &calendar.EventDateTime{
			Date: fmt.Sprintf(
				"%02d-%02d-%02d",
				currentTime.Year(),
				currentTime.Month(),
				currentTime.Day(),
			),
		},
		End: &calendar.EventDateTime{
			Date: fmt.Sprintf(
				"%02d-%02d-%02d",
				currentTime.Year(),
				currentTime.Month(),
				currentTime.Day(),
			),
		},
		Visibility:      "public",
		Transparency:    "transparent",
		GuestsCanModify: false,
		ColorId:         "8",
		Reminders:       nil,
	}
	boolFalse := false
	event.GuestsCanSeeOtherGuests = &boolFalse
	event.GuestsCanInviteOthers = &boolFalse
	c.setTotalWorkingHoursProperty(event, 0)

	evt, err := c.Service.Events.Insert(c.Id, event).Do()
	if err != nil {
		return nil, err
	}

	return evt, nil
}

func (c *Calendar) UpdateTotalWorkingEvent(
	totalWorkingEvent *calendar.Event,
	workingEvents *calendar.Events,
) (*calendar.Event, error) {
	if totalWorkingEvent == nil {
		return nil, fmt.Errorf("event is nil")
	}

	var hasPendingEvent bool = false
	var totalWorkingHours float64 = 0
	for _, item := range workingEvents.Items {
		workingHoursPropertyValue := c.GetWorkingHoursProperty(item)
		if workingHoursPropertyValue == "0.000" {
			hasPendingEvent = true
			break
		}
		if matched, err := regexp.MatchString("^\\d.\\d{3}$", workingHoursPropertyValue); matched == true &&
			err == nil {
			hrs, err := strconv.ParseFloat(workingHoursPropertyValue, 3)
			if err != nil {
				continue
			}
			totalWorkingHours += hrs
		}
	}

	if hasPendingEvent {
		return nil, fmt.Errorf("pending event exists")
	}

	totalWorkingEvent.Summary = fmt.Sprintf("Total Work (%.3f hrs)", totalWorkingHours)
	c.setTotalWorkingHoursProperty(totalWorkingEvent, totalWorkingHours)
	totalWorkingEvent.Reminders = nil

	evt, err := c.Service.Events.Update(c.Id, totalWorkingEvent.Id, totalWorkingEvent).Do()
	if err != nil {
		log.Printf("Total Working Event: %v\n", totalWorkingEvent.HtmlLink)
		return nil, err
	}

	return evt, nil
}
