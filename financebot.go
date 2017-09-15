package main

import (
  "fmt"
  "flag"
  "os"
  "os/signal"
  "syscall"

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
