package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/johnmikee/manifester/pkg/helpers"
)

func (c *Client) createDeptManifest(dept string) error {
	deptFile := fmt.Sprintf("%s/includes/%s", c.directory, dept)
	if _, err := os.Stat(deptFile); err == nil {
		c.log.Debug().Str("department", dept).Msg("department manifest exists")
	} else if errors.Is(err, os.ErrNotExist) {
		// does not exist - create
		c.log.Debug().Str("department", dept).Msg("creating department manifest")

		err := c.copyGroupManifest(c.directory + "/includes/" + dept)
		if err != nil {
			c.log.Info().AnErr("error", err).Str("department", dept).Msg("failed to copy manifest template")
			return err
		}
	} else {
		c.log.Debug().Str("file", deptFile).Msg("schrodinger says file may or may not exist.")
	}

	return nil
}

func (c *Client) manifests() error {
	manifestMachines, err := c.getDevices()
	if err != nil {
		c.log.Info().AnErr("error", err).Msg("failed to get devices")
		return err
	}

	// create a manifest for each machine and a map for quick lookup later
	machineMap := c.machineManifests(manifestMachines)

	// create the dept manifests
	groups := c.departmentManifest()

	// add the departments to the user manifests
	for group, members := range groups {
		for _, member := range members {
			serial, ok := machineMap[strings.Split(member, "@")[0]]
			if ok {
				err = addDeptToManifest(
					&UpdateInfo{
						directory:  c.directory,
						department: group,
						serial:     serial,
						user:       member,
					},
				)
				if err != nil {
					c.log.Info().AnErr("error", err).Str("serial", serial).Str("group", group).Msg("failed to add dept to manifest")
					continue
				}
			}
		}
	}

	return nil
}

func (c *Client) machineManifests(manifestMachines []MachineInfo) map[string]string {
	machineMap := make(map[string]string)
	for _, v := range manifestMachines {
		if !helpers.Contains(c.currentManifests(), v.Serial) {
			if !helpers.Contains(c.exclusions, v.Serial) {
				if _, err := os.Stat(v.Serial); os.IsNotExist(err) {
					c.copyTemplate(v.Serial, v.Username)
				}
			}
		}
		machineMap[v.Username] = v.Serial
	}
	return machineMap
}

func (c *Client) departmentManifest() map[string][]string {
	groups := c.oktaGroupMembers(c.filter)
	for group := range groups {
		err := c.createDeptManifest(group)
		if err != nil {
			c.log.Info().AnErr("error", err).Str("group", group).Msg("failed to create dept manifest")
			continue
		}
	}

	return groups
}
