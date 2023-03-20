
package main

//  ____                    _   ____    _                      _
// / ___|    ___    _   _  | | / ___|  | |__     __ _    ___  | | __
// \___ \   / _ \  | | | | | | \___ \  | '_ \   / _` |  / __| | |/ /
//  ___) | | (_) | | |_| | | |  ___) | | | | | | (_| | | (__  |   <
// |____/   \___/   \__,_| |_| |____/  |_| |_|  \__,_|  \___| |_|\_\
//  .  .  .  because  real  people  are  overrated

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	vip "github.com/spf13/viper"

	"github.com/lrstanley/girc"
	ai "github.com/sashabaranov/go-openai"
)

func getBanner() string {
	return fmt.Sprintf("%s\n%s",
		figure.NewColorFigure("SoulShack", "", "green", true).ColorString(),
		figure.NewColorFigure(" . . . because real people are overrated", "term", "green", true).ColorString())
}

func main() {
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

var root = &cobra.Command{
	Use:     "soulshack",
	Example: "soulshack --nick chatbot --server irc.freenode.net --port 6697 --channel '#soulshack' --ssl --openaikey ****************",
	Short:   getBanner(),
	Run:     run,
	Version: "0.42 - http://github.com/pkdindustries/soulshack",
}

func run(r *cobra.Command, _ []string) {

	aiClient := ai.NewClient(vip.GetString("openaikey"))

	if err := verifyConfig(vip.GetViper()); err != nil {
		log.Fatal(err)
	}

	irc := girc.New(girc.Config{
		Server:    vip.GetString("server"),
		Port:      vip.GetInt("port"),
		Nick:      vip.GetString("nick"),
		User:      "soulshack",
		Name:      "soulshack",
		SSL:       vip.GetBool("ssl"),
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})

	saslUser := vip.GetString("sasluser")
	saslPass := vip.GetString("saslpass")
	if saslUser != "" && saslPass != "" {
		irc.Config.SASL = &girc.SASLPlain{
			User: vip.GetString("saslnick"),
			Pass: vip.GetString("saslpass"),
		}
	}

	irc.Handlers.AddBg(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		ctx, cancel := CreateChatContext(context.Background(), aiClient, vip.GetViper(), c, &e)
		defer cancel()

		log.Println("joining channel:", ctx.Config.Channel)
		c.Cmd.Join(ctx.Config.Channel)

		time.Sleep(1 * time.Second)
		sendGreeting(ctx)
	})