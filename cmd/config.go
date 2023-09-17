package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/johnmikee/manifester/mdm"
	"github.com/johnmikee/manifester/mdm/client"
	"github.com/johnmikee/manifester/okta"
	"github.com/johnmikee/manifester/pkg/logger"
	"github.com/johnmikee/yae"
)

type Config struct {
	MDMToken   string `json:"mdm_token"`
	MDMURL     string `json:"mdm_url"`
	MDMUser    string `json:"mdm_user"`
	MDMPass    string `json:"mdm_pass"`
	OktaToken  string `json:"okta_token"`
	OktaURL    string `json:"okta_url"`
	OktaDomain string `json:"okta_domain"`
}

type Flags struct {
	configFile  string
	dryRun      bool
	service     string
	env         string
	logLevel    string
	logToFile   bool
	manifestDir string
	mdm         string
}

type Opts struct {
	Filter      string   `json:"department-filter"`
	Exclusions  []string `json:"exclusions"`
	DisplayName string   `json:"display-name"`
}

func readConf(cf string) *Opts {
	data, err := os.ReadFile(cf)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return nil
	}

	var opts Opts
	err = json.Unmarshal(data, &opts)
	if err != nil {
		fmt.Println("Error unmarshaling config:", err)
		return nil
	}

	return &opts
}

func setup() *Client {
	f := &Flags{
		configFile:  "config.json",
		dryRun:      false,
		env:         "dev",
		logLevel:    "debug",
		mdm:         "kandji",
		logToFile:   false,
		manifestDir: "munki_repo/manifests",
		service:     "manifester",
	}
	flag.StringVar(
		&f.configFile,
		"config-file",
		f.configFile,
		"Change config file location. [default: config.json]",
	)
	flag.BoolVar(
		&f.dryRun,
		"dry-run",
		f.dryRun,
		"Run without making changes.",
	)
	flag.StringVar(
		&f.env,
		"env",
		f.env,
		"Set the environment. [prod | dev]",
	)
	flag.BoolVar(
		&f.logToFile,
		"log-to-file",
		f.logToFile,
		"Log results to file.",
	)
	flag.StringVar(
		&f.logLevel,
		"log-level",
		f.logLevel,
		"Set the log level.",
	)
	flag.StringVar(
		&f.mdm,
		"mdm",
		f.mdm,
		"Select which mdm [jamf | kandji].",
	)
	flag.StringVar(
		&f.manifestDir,
		"manifest-dir",
		f.manifestDir,
		"Set the manifest directory.",
	)
	flag.StringVar(
		&f.service,
		"service",
		f.service,
		"Set the service name.",
	)
	flag.Parse()
	log := logger.NewLogger(
		&logger.Config{
			ToFile:  f.logToFile,
			Level:   f.logLevel,
			Service: f.service,
			Env:     f.env,
		},
	)

	var cfg Config
	err := yae.Get(yae.PROD,
		&yae.Env{
			Name:         f.service,
			Type:         yae.JSON,
			ConfigStruct: &cfg,
		},
	)
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to get config")
	}

	err = verifyManifestDir(f.manifestDir)
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to verify manifest directory")
	}

	opts := readConf(f.configFile)
	if opts == nil {
		log.Fatal().Msg("failed to read config")
		return nil
	}

	client := &Client{
		directory:  f.manifestDir,
		exclusions: opts.Exclusions,
		filter:     opts.Filter,
		log:        &log,
		mdm: client.New(
			&client.MDM{
				MDM: mdm.MDM(f.mdm),
				Config: mdm.Config{
					MDM:                    mdm.MDM(f.mdm),
					URL:                    cfg.MDMURL,
					User:                   cfg.MDMUser,
					Password:               cfg.MDMPass,
					Token:                  cfg.MDMToken,
					Client:                 nil,
					Log:                    log,
					ProviderSpecificConfig: nil,
				},
			},
		),
		okta: okta.New(
			&okta.Config{
				Domain: cfg.OktaDomain,
				Token:  cfg.OktaToken,
				URL:    cfg.OktaURL,
				Log:    &log,
				Client: nil,
			},
		),
	}

	return client
}
