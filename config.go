
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	ai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	vip "github.com/spf13/viper"
)

func init() {

	cobra.OnInitialize(initConfig)

	// irc client configuration
	root.PersistentFlags().StringP("nick", "n", "soulshack", "bot's nickname on the irc server")
	root.PersistentFlags().StringP("server", "s", "localhost", "irc server address")
	root.PersistentFlags().BoolP("ssl", "e", false, "enable SSL for the IRC connection")
	root.PersistentFlags().IntP("port", "p", 6667, "irc server port")
	root.PersistentFlags().StringP("channel", "c", "", "irc channel to join")
	root.PersistentFlags().StringP("saslnick", "", "", "nick used for SASL")
	root.PersistentFlags().StringP("saslpass", "", "", "password for SASL plain")

	// bot configuration
	root.PersistentFlags().StringP("become", "b", "chatbot", "become the named personality")
	root.PersistentFlags().StringP("directory", "d", "./personalities", "personalities configuration directory")
	root.PersistentFlags().StringSliceP("admins", "A", []string{}, "comma-separated list of allowed users to administrate the bot (e.g., user1,user2,user3)")

	// informational
	root.PersistentFlags().BoolP("list", "l", false, "list configured personalities")
	root.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging of sessions and configuration")

	// openai configuration
	root.PersistentFlags().String("openaikey", "", "openai api key")
	root.PersistentFlags().Int("maxtokens", 512, "maximum number of tokens to generate")
	root.PersistentFlags().String("model", ai.GPT4, "model to be used for responses (e.g., gpt-4)")

	// timeouts and behavior
	root.PersistentFlags().BoolP("addressed", "a", true, "require bot be addressed by nick for response")
	root.PersistentFlags().DurationP("session", "S", time.Minute*3, "duration for the chat session; message context will be cleared after this time")