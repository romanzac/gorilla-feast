package cmd

import (
	"fmt"
	"github.com/romanzac/gorilla-feast/controller/dbhandler"
	"github.com/romanzac/gorilla-feast/controller/httphandler"
	"github.com/romanzac/gorilla-feast/infra/config"
	"github.com/romanzac/gorilla-feast/infra/database"
	"github.com/romanzac/gorilla-feast/infra/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	cfgFile string

	// GorillaFeastCmd to start server application
	GorillaFeastCmd = &cobra.Command{
		Use:              "gorilla-feast",
		Short:            "My API tribute to Gorilla",
		Long:             `My API tribute to Gorilla`,
		Version:          "1.0.0",
		PersistentPreRun: initConfig,
		Run:              startGorillaFeast,
	}
)

func init() {
	GorillaFeastCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is gorilla-feast.yaml)")
}

// initConfig loads values from configuration file Gorilla Feast API
func initConfig(cmd *cobra.Command, args []string) {

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
	err := viper.ReadInConfig()

	if err == nil {
		// Read values from environment or configuration file
		config.Cfg.Web.Listen = viper.Get("Listen").(string)
		config.Cfg.Web.Port = viper.Get("Port").(string)
		config.Cfg.Web.DisableTLS = strings.ToLower(viper.Get("DisableTLS").(string))
		config.Cfg.Web.Key = viper.Get("Key").(string)
		config.Cfg.Web.Cert = viper.Get("Cert").(string)
		config.Cfg.Web.JWTPrivKey = viper.Get("JWTPrivKey").(string)
		config.Cfg.Web.JWTPubKey = viper.Get("JWTPubKey").(string)
		config.Cfg.Database.PostgresURI = viper.Get("PostgresURI").(string)
	} else {
		fmt.Fprintf(os.Stdout, "err loading config: %s", err)
		os.Exit(1)
	}
}

// startGorillaFeast starts Gorilla Feast API controller
func startGorillaFeast(cmd *cobra.Command, args []string) {

	// Initialize DB
	database.InitDB(config.Cfg.Database.PostgresURI)

	// Initialize router and failure channel
	r := router.NewRouter()

	// Initialize repositories
	userDBRepo := dbhandler.NewDbUserRepo()

	// Initialize APIs
	apiv1 := httphandler.NewAPIv1(userDBRepo)

	// Add routes
	httphandler.InitRoutes(r, apiv1)

	// Prepare config for HTTP server
	s := &http.Server{
		Addr:              config.Cfg.Web.Listen + ":" + config.Cfg.Web.Port,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		Handler:           r,
	}

	// Start HTTP server
	if config.Cfg.Web.DisableTLS == "yes" || config.Cfg.Web.DisableTLS == "true" {
		log.Fatal(s.ListenAndServe())
	} else {
		log.Fatal(s.ListenAndServeTLS(config.Cfg.Web.Cert, config.Cfg.Web.Key))
	}
}
