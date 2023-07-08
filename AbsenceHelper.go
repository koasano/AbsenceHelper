package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/jwt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

const (
	credentialsFile = "credentials.json"
	configFile = "config.json"
)

var (
	summary		string
	day			string
	endDay		string
	startTime	string
	endTime		string
	timeZone	string
	language	string
)

type Credentials struct {
	Installed InstalledCredentials `json:"installed"`
}

type InstalledCredentials struct {
	ClientID		string		`json:"client_id"`
	ClientSecret	string		`json:"client_secret"`
	RedirectURIs	[]string	`json:"redirect_uris"`
}

type Config struct {
	CalendarId string `json:"calendar_id"`
	TimeZone   string `json:"time_zone"`
	Language   string `json:"language"`
}

/*
 The function parses command-line arguments and adds an event to Google Calendar.
 コマンドライン引数を解析しGoogleカレンダーにイベントを追加する.
*/
func main() {
	log.Println("Starting the application...")

	creds, err := getCredentials(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to get credentials: %v", err)
	}

	config, err := getConfig(configFile)
	if err != nil {
		log.Fatalf("Unable to get config: %v", err)
	}

	err = parseAndValidateArgs()
	if err != nil {
		fmt.Println(err.Error())
		displayHelp()
		flag.Usage()
		return
	}

	event, err := createEvent()
	if err != nil {
		log.Fatalf("Unable to create event: %v", err)
	}

	err = addEventToCalendar(creds, config.CalendarId, event)
	if err != nil {
		log.Fatalf("Unable to add event to calendar: %v", err)
	}

	log.Println("Finished successfully.")
}

func getCredentials(file string) (*jwt.Config, error) {
	log.Println("Getting credentials...")

	//  Load the authentication information file for the service account.
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	// Create *jwt.Config from the authentication information.
	config, err := google.JWTConfigFromJSON(fileBytes, calendar.CalendarEventsScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	return config, nil
}

func getConfig(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		return nil, err
	}

	timeZone = config.TimeZone
	language = config.Language

	return &config, nil
}

func parseAndValidateArgs() error {
	log.Println("Parse and validate args...")

	if language == "ja" {
		flag.StringVar(&summary,   "s",  "", "不在情報の要約を格納する.例：'Asano 終日不在'")
		flag.StringVar(&day,       "d",  "", "不在日の日付を格納する.例：'2023-07-15'")
		flag.StringVar(&endDay,    "de", "", "複数日指定の場合、終了日を格納する.例：'2023-07-20'")
		flag.StringVar(&startTime, "tb", "", "時間指定の場合、開始時間を格納する.例：'13:00'")
		flag.StringVar(&endTime,   "te", "", "時間指定の場合、終了時間を格納する.例：'18:00'")
	} else {
		flag.StringVar(&summary,   "s",  "", "Stores the summary of the absence information. Example: 'Asano is absent all day'")
		flag.StringVar(&day,       "d",  "", "Stores the date of absence. Example: '2023-07-15'")
		flag.StringVar(&endDay,    "de", "", "In case of multiple days specified, stores the end date. Example: '2023-07-20'")
		flag.StringVar(&startTime, "tb", "", "In case of time specified, stores the start time. Example: '13:00'")
		flag.StringVar(&endTime,   "te", "", "In case of time specified, stores the end time. Example: '18:00'")
	}

	flag.Parse()

	// Checking option combinations.
	if summary == "" || day == "" {
		return errors.New("Required options -s and -d are not provided.")
	}

	if (endDay != "" && (startTime != "" || endTime != "")) || (startTime != "" && endTime == "" || startTime == "" && endTime != "") {
		return errors.New("Invalid combination of options.")
	}

	if len(os.Args) == 1 {
		return errors.New("No options provided.")
	}
	return nil
}

func displayHelp() {
	if language == "ja" {
		fmt.Println("\n使い方: ./AbsenceHelper -s <summary> -d <date> [-de <end date>] [-tb <start time>] [-te <end time>]")
		fmt.Println("\n例:")
		fmt.Println("   終日指定  : ./AbsenceHelper -s \"Asano 終日不在\" -d 2023-07-15")
		fmt.Println("   複数日指定: ./AbsenceHelper -s \"Asano 終日不在\" -d 2023-07-15 -de 2023-07-20")
		fmt.Println("   時間指定  : ./AbsenceHelper -s \"Asano AM休\" -d 2023-07-16 -tb 09:30 -te 13:00")
	} else {
		fmt.Println("\nUsage: ./AbsenceHelper -s <summary> -d <date> [-de <end date>] [-tb <start time>] [-te <end time>]")
		fmt.Println("\nExample:")
		fmt.Println("   All-day        : ./AbsenceHelper -s \"Asano is absent all day\" -d 2023-07-15")
		fmt.Println("   Multiple days  : ./AbsenceHelper -s \"Asano is absent all day\" -d 2023-07-15 -de 2023-07-20")
		fmt.Println("   Time specified : ./AbsenceHelper -s \"Asano is absent in the morning\" -d 2023-07-16 -tb 09:30 -te 13:00")
	}
}

func createEvent() (*calendar.Event, error) {
	var start *calendar.EventDateTime
	var end *calendar.EventDateTime

	if endDay == "" && startTime == "" && endTime == "" {
		log.Printf("Creating event: %s %s\n", summary, day)
		start = &calendar.EventDateTime{
			Date: day,
			TimeZone: timeZone,
		}
		end = &calendar.EventDateTime{
			Date: day,
			TimeZone: timeZone,
		}
	} else if endDay != "" {
		log.Printf("Creating event: %s %s - %s\n", summary, day, endDay)
		start = &calendar.EventDateTime{
			Date: day,
			TimeZone: timeZone,
		}
		end = &calendar.EventDateTime{
			Date: endDay,
			TimeZone: timeZone,
		}
	} else if startTime != "" && endTime != "" {
		log.Printf("Creating event: %s %s %s-%s\n", summary, day, startTime, endTime)
		start = &calendar.EventDateTime{
			DateTime: day + "T" + startTime + ":00",
			TimeZone: timeZone,
		}
		end = &calendar.EventDateTime{
			DateTime: day + "T" + endTime + ":00",
			TimeZone: timeZone,
		}
	} else {
		return nil, errors.New("Creating event: Invalid command line parameters for event date/time")
	}

	event := &calendar.Event{
		Summary: summary,
		Start: start,
		End: end,
	}

	return event, nil
}

func addEventToCalendar(config *jwt.Config, calendarId string, event *calendar.Event) error {
	log.Println("Adding event to calendar...")

	// Create a token source for the service account.
	client := config.Client(context.Background())

	// Create a service for Google Calendar.
	srv, err := calendar.New(client)
	if err != nil {
		return fmt.Errorf("Unable to retrieve calendar Client %v", err)
	}

	// Add the event to the calendar.
	_, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		return fmt.Errorf("Unable to create event. %v\n", err)
	}
	return nil
}

