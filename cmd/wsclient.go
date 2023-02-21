package cmd

import (
	"github.com/gorilla/websocket"
	"github.com/romanzac/gorilla-feast/infra/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

var (
	// Command to start ws client
	startWSClientCmd = &cobra.Command{
		Use:              "wsclient",
		Short:            "Start websocket client process",
		Long:             `Start websocket client process which receives messages about failed logins`,
		Version:          "1.0.0",
		PersistentPreRun: initWSConfig,
		Run:              startWSClient,
	}
)

func init() {
	GorillaFeastCmd.AddCommand(startWSClientCmd)
}

// initWSConfig loads config values for WSClient
func initWSConfig(cmd *cobra.Command, args []string) {

	// Read config from cfgFile
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in working directory with name "GorillaFeast" (without YAML extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.SetConfigName("gorilla-feast")
	}

	// Bind environment to read configuration values
	viper.SetEnvPrefix("GORILLA_FEAST")
	viper.AutomaticEnv()

	// Read the environment and configuration file
	if err := viper.ReadInConfig(); err == nil {
		config.Cfg.Web.Listen = viper.Get("Listen").(string)
		if config.Cfg.Web.Listen == "" {
			config.Cfg.Web.Listen = "localhost"
		}
		config.Cfg.Web.Port = viper.Get("Port").(string)
		config.Cfg.Web.DisableTLS = strings.ToLower(viper.Get("DisableTLS").(string))
	} else {
		os.Exit(1)
	}
}

// startWSClient starts websocket client to listen to failed sign ins
func startWSClient(cmd *cobra.Command, args []string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Prepare URL for the client
	var u url.URL
	if config.Cfg.Web.DisableTLS == "yes" || config.Cfg.Web.DisableTLS == "true" {
		u = url.URL{Scheme: "ws", Host: config.Cfg.Web.Listen + ":" + config.Cfg.Web.Port, Path: "/login-failures"}
	} else {
		u = url.URL{Scheme: "wss", Host: config.Cfg.Web.Listen + ":" + config.Cfg.Web.Port, Path: "/login-failures"}
	}

	log.Printf("Connected to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to the server: ", err)
	}
	defer c.Close()

	// Wait for new messages
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Error reading from the server: ", err)
				return
			}
			log.Println(string(message))
		}
	}()

	// Wait for the interrupt from keyboard and send close connection message to the server
	for range interrupt {
		log.Println("Closing connection on interrupt from keyboard")
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("Error sending close message to the server: ", err)
			return
		} else {
			c.Close()
			return
		}
	}
}
