package jamf

import (
	"strconv"
	"sync"

	"github.com/DataDog/jamf-api-client-go/classic"
	"github.com/johnmikee/manifester/mdm"
	"github.com/johnmikee/manifester/pkg/logger"
)

var wg sync.WaitGroup

type Client struct {
	log    logger.Logger
	client *classic.Client
	info   []mdm.MachineInfo
}

// Setup implements mdm.Provider.
func (c *Client) Setup(config mdm.Config) {
	c.log = logger.ChildLogger("jamf", &config.Log)

	jc, err := classic.NewClient(config.URL,
		config.User,
		config.Password,
		nil,
	)
	if err != nil {
		config.Log.Fatal().AnErr("error", err).Msg("building jamf client")
	}

	c.client = jc
}

func (c *Client) ListAllDevices() ([]mdm.MachineInfo, error) {
	computers, err := c.client.Computers()
	if err != nil {
		c.log.Info().AnErr("error", err).Msg("getting computers")
		return nil, err
	}

	ch := make(chan classic.BasicComputerInfo, 20)

	for t := 0; t < 5; t++ {
		wg.Add(1)
		go c.infoGrabber(ch, &wg)
	}

	for _, v := range computers {
		ch <- v
	}

	close(ch)
	wg.Wait()

	return c.info, nil
}

func (c *Client) infoGrabber(ch chan classic.BasicComputerInfo, wg *sync.WaitGroup) {
	for v := range ch {
		res, err := c.client.ComputerDetails(v.ID)
		if err != nil {
			c.log.Info().AnErr("error", err).Msg("getting computer details")
		}

		c.info = append(c.info, mdm.MachineInfo{
			Device: mdm.Device{
				DeviceID:     strconv.Itoa(res.Info.ID),
				Hostname:     res.Info.General.Name,
				SerialNumber: res.Info.General.SerialNumber,
			},
			Users: &mdm.User{
				Email: res.Info.UserLocation.EmailAddress,
				Name:  res.Info.UserLocation.RealName,
				ID:    res.Info.ID,
			},
		})
		wg.Done()
	}
}
