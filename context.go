
package main

import (
	"context"
	"log"
	"strings"

	"github.com/lrstanley/girc"
	ai "github.com/sashabaranov/go-openai"
	vip "github.com/spf13/viper"
)

type Personality struct {
	Prompt   string
	Greeting string
	Nick     string
	Model    string
	Goodbye  string
}

type Config struct {
	Channel   string
	Nick      string
	Admins    []string
	Directory string
	Verbose   bool
	Server    string
	Port      int
	SSL       bool
	Addressed bool
}

type ChatContext struct {
	context.Context
	AI          *ai.Client
	Personality *Personality
	Config      *Config
	Client      *girc.Client
	Event       *girc.Event
	Session     *ChatSession
	Args        []string
}

func PersonalityFromViper(v *vip.Viper) *Personality {
	return &Personality{
		Prompt:   v.GetString("prompt"),
		Greeting: v.GetString("greeting"),
		Nick:     v.GetString("nick"),
		Model:    v.GetString("model"),
		Goodbye:  v.GetString("goodbye"),
	}
}

func IrcFromViper(v *vip.Viper) *Config {
	return &Config{
		Channel:   v.GetString("channel"),
		Nick:      v.GetString("nick"),
		Admins:    v.GetStringSlice("admins"),
		Directory: v.GetString("directory"),
		Verbose:   v.GetBool("verbose"),
		Server:    v.GetString("server"),
		Port:      v.GetInt("port"),
		SSL:       v.GetBool("ssl"),