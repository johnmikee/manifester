package cmd

import (
	"os"

	"github.com/johnmikee/manifester/mdm"
	"github.com/johnmikee/manifester/okta"
	"github.com/johnmikee/manifester/pkg/logger"
)

type Client struct {
	mdm        mdm.Provider
	okta       *okta.Client
	log        *logger.Logger
	directory  string   // munki manifest directory
	exclusions []string // serial numbers to exclude
	filter     string   // okta filter
}

func Execute() {
	client := setup()
	err := client.run()
	if err != nil {
		client.log.Info().AnErr("error", err).Msg("failed to successfully generate manifests")
		os.Exit(1)
	}
}

func (c *Client) run() error {
	/*
		first, remove the current entries. we exlude the includes/ directory
		as these are manually created and managed. this ensures that if a user
		switched departments their manifests are updated accordingly and is less expensive than
		opening each file to check if it is current and then closing it.
	*/

	err := c.removeEntries()
	if err != nil {
		c.log.Info().AnErr("error", err).Msg("failed to remove entries")
		return err
	}
	/*
		next we build our manifests. we do this by getting all the machines from kandji
		and then iterating through them taking the serial number to make the manifest.
		we take the user assigned to the device, unless we cannot, and get their department from okta.

		this allows us to target specific groups of users with specific manifests. or not.
	*/
	return c.manifests()
}
