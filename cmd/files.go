package cmd

import (
	"os"

	"github.com/johnmikee/manifester/pkg/helpers"
)

func (c *Client) removeEntries() error {
	files, err := os.ReadDir(c.directory)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !helpers.Contains(c.exclusions, file.Name()) {
			filePath := c.directory + "/" + file.Name()
			// make sure its not a directory
			if file.IsDir() {
				c.log.Debug().Str("file", filePath).Msg("skipping directory")
				continue
			}
			c.log.Debug().Str("file", filePath).Msg("removing")
			err := os.Remove(filePath)
			if err != nil {
				c.log.Info().AnErr("error", err).Msg("failed to remove file")
			}
		}
	}
	return nil
}

func verifyManifestDir(dir string) error {
	_, err := os.Stat(dir)

	return err
}
