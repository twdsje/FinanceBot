package main

import (
  "fmt"
  "flag"
  "os"
  "os/signal"
  "syscall"
  "strings"
  "golang.org/x/net/html"
  "net/http"
  "time"

  "github.com/bwmarrin/discordgo"
)

func init() {
  flag.StringVar(&token, "t", "", "Bot Token")
  flag.Parse()
}
var token string
var buffer = make([][]byte, 0)

func main() {
  fmt.Println("Starting discord finance bot")

  if token == "" {
    fmt.Println("No token provided. Please run: financebot -t <bot token>")
    return
  }

  newsQueue := initCalendar()
  for _, item := range newsQueue {
    fmt.Printf("Id: %v Date: %v Name: %v Forecast: %v Previous: %v \n", item.eventId, item.date, item.event, item.forecast, item.previous)
  }

  // Create a new Discord session using the provided bot token.
  fmt.Println("Starting discord session")
  dg, err := discordgo.New("Bot " + token)
  if err != nil {
    fmt.Println("Error creating Discord session: ", err)
    return
  }

  // Register messageCreate as a callback for the messageCreate events.
  dg.AddHandler(messageCreate)

  // Open the websocket and begin listening.
  fmt.Println("Opening discord session")
  err = dg.Open()
  if err != nil {
    fmt.Println("Error opening Discord session: ", err)
  }

  // Wait here until CTRL-C or other term signal is received.
  fmt.Println("Finance bot is now running.  Press CTRL-C to exit.")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc

  // Cleanly close down the Discord session.
  dg.Close()


}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

  // Ignore all messages created by the bot itself
  // This isn't required in this specific example but it's a good practice.
  if m.Author.ID == s.State.User.ID {
    return
  }


  // check if the message is "!airhorn"
  if strings.HasPrefix(m.Content, "!financebot ping") {
    fmt.Println("PING")
    s.ChannelMessageSend(m.ChannelID, "PONG!")

  }
}

//Read the calendar.
func initCalendar() []*NewsEvent {
  curDate, curTime := "", ""

  newsQueue := make([]*NewsEvent, 0)

  fmt.Println("Grabbing data from forexfactory")
  resp, _ := http.Get("https://www.forexfactory.com/calendar.php")

  z := html.NewTokenizer(resp.Body)

  for {
    tt := z.Next()

    switch {
      case tt == html.ErrorToken:
        // End of the document, we're done
	return newsQueue
      case tt == html.StartTagToken:
        t := z.Token()

        isRow := t.Data == "tr"
        if isRow {

          ok, classname := getClass(t)
          if !ok {
            continue
          }
          if classname != "calendar__row calendar__expand calendar__row--alt" && classname != "calendar__row calendar__expand" {
            ok, eventId := getEventId(t)
            if ok {
              mydate, mytime, event, forecast, previous := parseRow(z)
              if mydate != "" {
                curDate = mydate
              }
              if mytime != "" {
                curTime = mytime
              }

              //Create struct
              const longForm = "Mon Jan 2 3:04pm MST 2006"
              dt, _ := time.Parse(longForm, curDate + " " + curTime + " EST 2017")
              item := NewsEvent{eventId: eventId, date: dt, event: event, forecast: forecast, previous: previous }
              newsQueue = append(newsQueue, &item)
            }
          }

          continue
        }
    }
  }
}

func parseRow(z *html.Tokenizer)(date string, time string, event string, forecast string, previous string){
  date, time, event, forecast, previous = "", "", "", "", ""
  for {
    tt := z.Next()

    switch {
      case tt == html.ErrorToken:
        // End of the document, we're done
        return
      case tt == html.StartTagToken:
        t := z.Token()

        isCell := t.Data == "td"
        if isCell {
          ok, classname := getClass(t)
          if !ok {
            continue
          }

          if classname == "calendar__cell calendar__date date" {
            //fmt.Printf("Found date: %v \n", 0)
            ok, tmp := getText(z)
            if ok {
              date = tmp
            }
            continue
          }

          if classname == "calendar__cell calendar__time time" {
            ok, tmp := getText(z)
            if ok {
              time = tmp
            }
            continue
          }

          if classname == "calendar__cell calendar__event event" {
            ok, tmp := getText(z)
            if ok {
              event = tmp
            }
            continue
          }

          if classname == "calendar__cell calendar__forecast forecast" {
            ok, tmp := getText(z)
            if ok {
              forecast = tmp
            }
            continue
          }

          if classname == "calendar__cell calendar__previous previous" {
            ok, tmp := getText(z)
            if ok {
              previous = tmp
            }
            continue
          }
        }

      case tt == html.EndTagToken:
        t := z.Token()

        isRow := t.Data == "tr"
        if isRow {
          return
        }
    }
  }
}

func getText(z *html.Tokenizer)(ok bool, val string) {
  val = ""
  for {
    tt := z.Next()

    switch {
      case tt == html.TextToken:
        t := z.Token()
        val += t.Data + " "

     case tt == html.EndTagToken:
        t := z.Token()

        isEndOfCell := t.Data == "td"
        if isEndOfCell {
          if val != "" {
            ok = true
          }
          val = strings.TrimSpace(val)
          return
        }

    }
  }
}


func getClass(t html.Token) (ok bool, href string) {
  // Iterate over all of the Token's attributes until we find an "href"
  for _, a := range t.Attr {
    if a.Key == "class" {
      href = a.Val
      ok = true
    }
  }

  // "bare" return will return the variables (ok, href) as defined in
  // the function definition
  return
}

func getEventId(t html.Token) (ok bool, val string) {
  for _, a := range t.Attr {
    if a.Key == "data-eventid" {
      val = a.Val
      ok = true
    }
  }

  return
}

type NewsEvent struct {
  eventId string
  date time.Time
  event string
  forecast string
  previous string
}

