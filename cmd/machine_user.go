package cmd

import (
	"strings"
)

type MachineInfo struct {
	Serial   string
	Username string
}

func (c *Client) getDevices() ([]MachineInfo, error) {
	machines, err := c.mdm.ListAllDevices()
	if err != nil {
		c.log.Fatal().AnErr("error", err).Msg("failed to get devices")
		return nil, err
	}
	var manifestMachines []MachineInfo
	for _, machine := range machines {
		m := MachineInfo{
			Serial: machine.Device.SerialNumber,
		}
		if machine.Users != nil {
			m.Username = strings.Split(machine.Users.Email, "@")[0]
		} else {
			m.Username = ""
		}

		manifestMachines = append(manifestMachines, m)
	}

	return manifestMachines, nil
}

func (c *Client) oktaGroupMembers(filter string) map[string][]string {
	oktaGroups, err := c.okta.ListGroups()
	if err != nil {
		c.log.Info().AnErr("error", err).Msg("failed to get okta groups")
		return nil
	}

	gm := oktaGroups.GetMembers(c.okta, &filter)

	return gm
}
