
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
	root.PersistentFlags().DurationP("timeout", "t", time.Second*60, "timeout for each completion request to openai")
	root.PersistentFlags().IntP("history", "H", 15, "maximum number of lines of context to keep per session")
	root.PersistentFlags().DurationP("chunkdelay", "C", time.Second*15, "after this delay, bot will look to split the incoming buffer on sentence boundaries")
	root.PersistentFlags().IntP("chunkmax", "m", 350, "maximum number of characters to send as a single message")

	// personality / prompting
	root.PersistentFlags().String("goodbye", "goodbye.", "prompt to be used when the bot leaves the channel")
	root.PersistentFlags().String("greeting", "hello.", "prompt to be used when the bot joins the channel")
	root.PersistentFlags().String("prompt", "respond in a short text:", "initial system prompt for the ai")

	vip.BindPFlags(root.PersistentFlags())

	vip.SetEnvPrefix("SOULSHACK")
	vip.AutomaticEnv()
}

func initConfig() {

	fmt.Println(getBanner())

	if _, err := os.Stat(vip.GetString("directory")); errors.Is(err, fs.ErrNotExist) {
		log.Printf("? configuration directory %s does not exist", vip.GetString("directory"))
	}

	if vip.GetBool("list") {
		personalities := listPersonalities()
		log.Printf("Available personalities: %s", strings.Join(personalities, ", "))
		os.Exit(0)
	}

	vip.AddConfigPath(vip.GetString("directory"))