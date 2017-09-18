package main

import (
  "fmt"
  "flag"
  "os"
  "os/signal"
  "syscall"
  "strings"

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
  fmt.Println("Finance bot is now running.  Meep!  Press CTRL-C to exit.")
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
