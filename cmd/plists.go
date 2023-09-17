package cmd

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/johnmikee/manifester/pkg/helpers"
	"howett.net/plist"
)

// UpdateInfo contains the information needed to update a manifest
type UpdateInfo struct {
	directory  string
	department string
	file       string
	serial     string
	user       string
}

func (c *Client) copyGroupManifest(dest string) error {
	source := fmt.Sprintf("%s/includes/department_template", c.directory)

	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func (c *Client) copyTemplate(serial string, user string) error {
	if _, err := os.Stat(c.directory); os.IsNotExist(err) {
		c.log.Error().AnErr("error", err).Str("directory", c.directory).Msg("failed to stat")
	}

	_, err := os.Create(c.directory + "/" + serial)
	if err != nil {
		c.log.Info().AnErr("error", err).Str("directory", c.directory).Msg("failed to create")
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/%s", c.directory, serial), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		c.log.Info().AnErr("error", err).Str("directory", c.directory).Msg("failed to open")
	}

	if user == "" {
		return noUser(f)
	}

	return withUser(user, f)
}

func (c *Client) currentManifests() []string {
	files, err := os.ReadDir(c.directory)
	if err != nil {
		c.log.Error().AnErr("error", err).Str("directory", c.directory).Msg("failed to read")
	}

	contents := []string{}
	for _, file := range files {
		contents = append(contents, file.Name())
	}
	return contents
}

func addDeptToManifest(u *UpdateInfo) error {
	u.file = fmt.Sprintf("%s/%s", u.directory, u.serial)

	return updatePlist(u)
}

func noUser(f *os.File) error {
	t := template.Must(template.New("manifest").Parse(unknownUserManifestTemplate()))
	return t.Execute(f, nil)
}

func withUser(user string, f *os.File) error {
	data := struct {
		Name string
	}{
		Name: user,
	}

	t := template.Must(template.New("manifest").Parse(userManifestTemplate()))
	return t.Execute(f, data)
}

func userManifestTemplate() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>display_name</key>
	<array>
		<string>{{.Name}}</string>
	</array>
	<key>catalogs</key>
	<array>
		<string>production</string>
	</array>
	<key>included_manifests</key>
	<array>
		<string>includes/apple_apps</string>
		<string>includes/common_base</string>
		<string>includes/optional_apps</string>
		<string>includes/security</string>
	</array>
</dict>
</plist>
	`
}

func unknownUserManifestTemplate() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>catalogs</key>
	<array>
		<string>production</string>
	</array>
	<key>included_manifests</key>
	<array>
		<string>includes/apple_apps</string>
		<string>includes/common_base</string>
		<string>includes/optional_apps</string>
	</array>
</dict>
</plist>
	`
}

func updatePlist(u *UpdateInfo) error {
	f, err := os.Open(u.file)
	if err != nil {
		return fmt.Errorf("failed to open plist: %w", err)
	}
	defer f.Close()

	p, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read plist: %w", err)
	}

	var pl map[string]interface{}
	_, err = plist.Unmarshal(p, &pl)
	if err != nil {
		return fmt.Errorf("failed to unmarshal plist: %w", err)
	}

	department := fmt.Sprintf("includes/%s", u.department)
	for key, val := range pl {
		if key == "included_manifests" {
			included := make([]string, len(val.([]interface{})))
			for i, v := range val.([]interface{}) {
				included[i] = v.(string)
			}
			if !helpers.Contains(included, department) {
				included = append(included, department)
				pl[key] = included
			}
		}
	}

	updatedPlist, err := plist.MarshalIndent(pl, plist.XMLFormat, "\t")
	if err != nil {
		return err
	}

	manifest := fmt.Sprintf("%s/%s", u.directory, u.serial)

	return os.WriteFile(manifest, updatedPlist, 0o644)
}
