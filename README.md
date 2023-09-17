## Manifester
Manifester is a simple, and slightly opinionated, program to generate a personalized manifest for each user in your MDM.

Munki manifests can be configured in a vast number of ways, usually by department, team, or some other high-level grouping. This approach works well, until it doesnt.
By creating a personalized manifest for each user we enable the ability to target a specific user(s) machine. You may not often need to do this,
but when you do it will be a life saver. Better to have the ability and not need it than to need it and not have it.

The use of the higher-level manifests is not without merit though and to that end, we will do both.

## Before you begin
* This program assumes you have somewhat standard munki repo layout. If you do not, you will need to modify the code to match your layout.
* Under `manifests/includes` you will need to create a manifest named `department_template`. This manifest will be used to create any department manifests that do not already exist.
* The user template has some default manifests that are included. You will need to create these manifests in your repo. If you do not want to use these manifests you will need to modify the code to replace the items listed below. 
    ```
    <key>included_manifests</key>
	<array>
		<string>includes/apple_apps</string>
		<string>includes/common_base</string>
		<string>includes/optional_apps</string>
		<string>includes/security</string>
	</array>
    ``` 
    
## How it works
First, we gather all the device information from the MDM. We use that to create a manifest for each device. While doing this we also add a key to the manifest `display_name` that contains the value set in the [config](config.json) for the `display-name` key. This information is not used by munki, but it makes it easier to identify the manifest a user has without knowing their serial, mac address, or information used for the manifest name.

Once the manifests are created we then query Okta to build a map of departments and the users which belong to them. We then use this map to add the department to the manifest.

### Department Filter
When the departments are pulled from Okta and optional filter can be applied to only include departments that match the filter.
If no filter is specified all departments will be included.
To add a filter edit the [config](config.json) and modify the `department-filter` value.

## Exclusions
To add a machine to the exclusion's edit the [config](config.json) and add the serial number to the list under the `exclusions` key.
Ex:
    Say you want to exclude C02ABC123
```
{
    "exclusions": [
        "C02ABC123"
    ]
}
```

##
## Note
The code under Jamf is not currently used. I no longer have access to a Jamf instance to test with, and this was put together by pasting together old memories and referencing the API docs. If you would like to add support for Jamf, or anything else, please feel free to submit a PR.
