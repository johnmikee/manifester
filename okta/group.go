package okta

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Groups []Group

type Group []struct {
	ID                    string       `json:"id"`
	Created               time.Time    `json:"created"`
	LastUpdated           time.Time    `json:"lastUpdated"`
	LastMembershipUpdated time.Time    `json:"lastMembershipUpdated"`
	ObjectClass           []string     `json:"objectClass"`
	Type                  string       `json:"type"`
	Profile               GroupProfile `json:"profile,omitempty"`
	Links                 Links        `json:"_links,omitempty"`
	Source                Source       `json:"source,omitempty"`
}

type GroupProfile struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	GroupType      string `json:"groupType,omitempty"`
	SamAccountName string `json:"samAccountName,omitempty"`
	ObjectSid      string `json:"objectSid,omitempty"`
	GroupScope     string `json:"groupScope,omitempty"`
	Dn             string `json:"dn,omitempty"`
	ExternalID     string `json:"externalId,omitempty"`
}

type GroupMembers []struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	Created       time.Time `json:"created"`
	Activated     time.Time `json:"activated"`
	StatusChanged time.Time `json:"statusChanged"`
	LastLogin     time.Time `json:"lastLogin"`
	LastUpdated   time.Time `json:"lastUpdated"`
	Profile       Profile   `json:"profile"`
	Links         Links     `json:"_links"`
}

type Provider struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Self struct {
	Href string `json:"href"`
}

type Source struct {
	ID string `json:"id"`
}

type Links struct {
	Self Self `json:"self"`
}

var groupBase = "groups"

// GetMembers returns a map of group names to a slice of email addresses
func (g Groups) GetMembers(o *Client, filter *string) map[string][]string {
	m := make(map[string][]string)
	idNameMap := g.idNameMap(filter)

	for group, name := range idNameMap {
		gr, err := o.getGroupsMembers(group)
		if err != nil {
			o.log.Error().Err(err).Msg("error getting group members")
			continue
		}
		for _, member := range gr {
			m[name] = append(m[name], member.Profile.Email)
		}
	}

	return m
}

func (o *Client) getGroupsMembers(groupID string) (GroupMembers, error) {
	req, err := o.newRequest(http.MethodGet, fmt.Sprintf("%s/%s/users", groupBase, groupID), false, nil)
	if err != nil {
		return nil, err
	}

	var groupMembers GroupMembers
	_, err = o.do(req, &groupMembers)
	if err != nil {
		return nil, err
	}

	return groupMembers, nil
}

// ListGroups queries the groups endpoint and paginates until all groups have been returned.
func (o *Client) ListGroups() (Groups, error) {
	var override bool
	var group Group

	groups := Groups{}
	url := groupBase
	for {
		resp, err := o.listGroups(url, override, &group)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
		value := resp.Header["Link"]

		link := linkSorter(value)
		if link == "" {
			o.log.Trace().Msg("no more responses from okta")
			break
		}

		url = link
		override = true
		o.log.Trace().Msg("checking next link..")
	}

	return groups, nil
}

func (o *Client) listGroups(url string, override bool, group *Group) (*http.Response, error) {
	req, err := o.newRequest(http.MethodGet, url, override, nil)
	if err != nil {
		o.log.Error().Err(err).Msg("error creating request")
		return nil, err
	}

	resp, err := o.do(req, group)
	if err != nil {
		o.log.Error().Err(err).Msg("error making request")
		return nil, err
	}

	return resp, nil
}

// MakeIDNameMap returns a map of group IDs to group names with
// an optional filter. right now this is a simple strings.hasPrefix
// check on the group name.
func (g Groups) MakeIDNameMap(filter *string) map[string]string {
	return g.idNameMap(filter)
}

func (g Groups) idNameMap(f *string) map[string]string {
	m := make(map[string]string)
	for _, groups := range g {
		for group := range groups {
			if f != nil {
				if !strings.HasPrefix(groups[group].Profile.Name, *f) {
					continue
				}
			}
			m[groups[group].ID] = groups[group].Profile.Name
		}
	}

	return m
}
