package kandji

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/johnmikee/manifester/mdm"
)

// DeviceResults is a list of DeviceResult
type DeviceResults []DeviceResult

func (c *Client) GetDeviceDetails(d string) (*DeviceDetails, error) {
	return c.deviceDetails(d)
}

func (c *Client) deviceDetails(d string) (*DeviceDetails, error) {
	u := fmt.Sprintf("devices/%s/details", d)

	req, err := c.newRequest(http.MethodGet, u, false, nil)
	if err != nil {
		c.log.Debug().AnErr("error", err).Msg("building request")
		return nil, err
	}

	var deviceDetails DeviceDetails
	_, err = c.do(req, &deviceDetails)
	if err != nil {
		c.log.Debug().AnErr("error", err).Msg("making request")
		return nil, err
	}

	return &deviceDetails, nil
}

func (c *Client) list(limit, offset int) (DeviceResults, error) {
	url := fmt.Sprintf("devices?limit=%d&offset=%d", limit, offset)

	req, err := c.newRequest(http.MethodGet, url, false, nil)
	if err != nil {
		c.log.Debug().AnErr("error", err).Msg("building request")
		return nil, err
	}

	var res DeviceResults
	_, err = c.do(req, &res)
	if err != nil {
		c.log.Debug().
			AnErr("err", err).
			Str("url", req.URL.String()).
			Msg("error making request")
		return nil, err
	}

	return res, nil
}

func (c *Client) listDevices(limit, offset int) (DeviceResults, error) {
	return c.list(limit, offset)
}

// ListAllDevices will paginate through devices until there is no response and return the results
func (c *Client) listAllDevices() (DeviceResults, error) {
	opts := &offsetRange{
		Limit:  300,
		Offset: 0,
	}

	res := DeviceResults{}
	for {
		results, err := c.listDevices(opts.Limit, opts.Offset)
		if err != nil {
			c.log.Info().AnErr("error", err).Msg("listing devices")
			return results, err
		}
		opts.Offset += opts.Limit
		if len(results) == 0 {
			break
		}
		res = append(res, results...)
	}

	return res, nil
}

func (c *Client) ListAllDevices() ([]mdm.MachineInfo, error) {
	devices, err := c.listAllDevices()
	if err != nil {
		c.log.Info().AnErr("error", err).Msg("listing devices")
		return nil, err
	}

	var res []mdm.MachineInfo
	for _, device := range devices {
		m := mdm.MachineInfo{
			Device: mdm.Device{
				DeviceID:     device.DeviceID,
				Hostname:     device.DeviceName,
				SerialNumber: device.SerialNumber,
			},
		}
		if device.User != nil {
			m.Users = &mdm.User{
				Email: device.User.UserClass.Email,
				Name:  device.User.UserClass.Name,
				ID:    int(device.User.UserClass.ID),
			}
		} else {
			m.Users = nil
		}

		res = append(res, m)
	}

	return res, nil
}

func unmarshalDeviceResults(data []byte) (DeviceResults, error) {
	var r DeviceResults
	err := json.Unmarshal(data, &r)
	return r, err
}

func unmarshalDeviceDetails(data []byte) (DeviceDetails, error) {
	var r DeviceDetails
	err := json.Unmarshal(data, &r)
	return r, err
}
