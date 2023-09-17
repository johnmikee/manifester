package mdm

import (
	"net/http"

	"github.com/johnmikee/manifester/pkg/logger"
)

// MDM represents the type of MDM provider.
type MDM string

const (
	Jamf   MDM = "jamf"
	Kandji MDM = "kandji"
)

// Provider represents the interface for an MDM provider.
type Provider interface {
	Setup(config Config)
	ListAllDevices() ([]MachineInfo, error)
}

// Config is the struct that is used to configure the MDM client.
type Config struct {
	MDM                    MDM           `json:"mdm,omitempty"`
	User                   string        `json:"user,omitempty"`
	Password               string        `json:"password,omitempty"`
	URL                    string        `json:"url,omitempty"`
	Token                  string        `json:"token,omitempty"`
	Client                 *http.Client  `json:"client,omitempty"`
	Log                    logger.Logger `json:"log,omitempty"`
	ProviderSpecificConfig interface{}   `json:"provider_specific_config,omitempty"`
}

// MachineInfo holds the device and user information
type MachineInfo struct {
	Device Device `json:"general"`
	Users  *User  `json:"users"`
}

// Device holds the general purpose information of the device
type Device struct {
	DeviceID     string `json:"device_id"`
	Hostname     string `json:"host_name"`
	SerialNumber string `json:"serial_number"`
}

// User holds the general purpose information of the user
type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
}
