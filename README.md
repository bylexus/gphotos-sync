# Google Photos Sync/Backup

This Go program downloads a (filtered) collection of a Google Photos account into a local directory
hierarchy.

Unfortunately, Google does not offer any simple solution to just download photos by certain filters,
so this tool fills this gap.

Luckily, Google Photos comes with a REST API, which allows querying and downloading the Photos library.
This script is a very simple approach to fetch your Google Photos automatically.

**Note**: This is a very first draft of the script. It is NOT tested thoroughly, and it is NOT implemented in a secure and stable way! Be aware!

## Requirements

* The Go tools, if you want to compile it by yourself
* A Google Photos account
* A Google Photos account and access to the Google Developer Console (https://console.developers.google.com/)

## Setup

### Step 1: Enable the Photos API

You need to enable the Photos API in order to use it.

1. Log in to https://console.developers.google.com/
2. In the Dashboard, click on "Enable APIs and Services":<br/>
   ![Enable APIs Menu](./doc/img/gcloud_enable_apis.png)
3. Search for "Photos", and select the Photos entry:<br />
   ![Photos Menu Entry](./doc/img/gcloud_photos_entry.png)
4. Click the "Enable API" button

Return to the Developer Console Dashboard.

### Step 2: Create a Google OAuth 2.0 Client ID

The Google Photos API only works with OAuth 2.0. You need to create an OAuth 2.0 Client ID and download your secret credentials
in order to use them with the pyhton tool.

1. In the Developer Console, click "Credentials" in the left menu, then select "Create Credentials" &gt; "OAuth Client ID":<br/>
   ![OAuth Client ID Menu Entry](./doc/img/google-cloud-api-create_login_menu.png)
2. Select "Other" in the Create OAuth client ID menu, and provide a name. This name is not important for the application, just for you to
   recognize it later.<br/>
   ![OAuth Client ID creation](./doc/img/google-cloud-api-create-oauth-client-id.png)
3. Click "Create".
4. Back in the Dashboard's "Credentials" menu, click the download button on the newly created client it. The JSON file you get contains
   the secret client tokens. KEEP IT IN A SAFE PLACE, it contains the keys that allow access the Google Photos library! <br />
   ![OAuth Client ID credential download](./doc/img/google-cloud-api-download-client_creds.png)


### Step 3: Build the program

```bash
$ git clone git@github.com:bylexus/gphotos-sync.git
$ cd gphotos-sync
$ go build
```

## Usage

Once the setup is complete, you start the tool with the following command:

```bash
$ ./gphoto-sync /path/to/photo/backup/folder
```

If you start the tool for the first time, it initiates the Google Authentication workflow:

1. Your Browser will start and asks you for your Google permissions for the gphotos-sync app
2. After the authentication was successful, you can close the browser and the command line
   tool should resume its operation.

***NOTE**

The Google Authentication information is stored locally in a settings file in the User's config directory (`secrets.json`). 
This file contains sensitive information (your authentication token). Do NOT share this file!

If all goes well, the tool begins to transfer your photos. Enjoy!

## Notes


* This tool is in an early experimental stage. Do NOT expect it to work flawlessly.
* I have not tested it over a longer period of time (e.g. until the client tokens expire). I do not yet know what happens then.
* I will implement some additional config params in the future:
  * credentials store location
  * credentials store encryption
  * output folder configuration: folder level generation
  * ...?

If you have any more ideas what this tool should be able to do, please drop me an issue or a note.

