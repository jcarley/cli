# Overview 

```

Usage:  [OPTIONS] COMMAND [arg...]


Options:
  -U, --username    Catalyze Username ($CATALYZE_USERNAME)
  -P, --password    Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env         The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version     Show the version and exit

Commands:
  associate	Associates an environment
  associated	Lists all associated environments
  certs	Manage your SSL certificates and domains
  clear	Clear out information in the global settings file to fix a misconfigured CLI.
  console	Open a secure console to a service
  dashboard	Open the Catalyze Dashboard in your default browser
  db	Tasks for databases
  default	[DEPRECATED] Set the default associated environment
  deploy-keys	Tasks for SSH deploy keys
  disassociate	Remove the association with an environment
  domain	Print out the temporary domain name of the environment
  environments	Manage environments for which you have access
  files	Tasks for managing service files
  git-remote	Manage git remotes to Catalyze code services
  invites	Manage invitations for your organizations
  keys	Tasks for SSH keys
  logout	Clear the stored user information from your local machine
  logs	Show the logs in your terminal streamed from your logging dashboard
  metrics	Print service and environment metrics in your local time zone
  rake	Execute a rake task
  redeploy	Redeploy a service without having to do a git push
  releases	Manage releases for code services
  rollback	Rollback a code service to a specific release
  services	Perform operations on an environment's services
  sites	Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl	Perform operations on local certificates to verify their validity
  status	Get quick readout of the current status of your associated environment and all of its services
  support-ids	Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update	Checks for available updates and updates the CLI if a new update is available
  users	Manage users who have access to the given organization
  vars	Interaction with environment variables for the associated environment
  version	Output the version and quit
  whoami	Retrieve your user ID
  worker	Manage a service's workers

Run ' COMMAND --help' for more information on a command.

```



#  Associate

```

Usage:  associate ENV_NAME SERVICE_NAME [-a] [-r] [-d]

Associates an environment

Arguments:
  ENV_NAME=""       The name of your environment
  SERVICE_NAME=""   The name of the primary code service to associate with this environment (i.e. 'app01')

Options:
  -a, --alias=""            A shorter name to reference your environment by for local commands
  -r, --remote="catalyze"   The name of the remote
  -d, --default=false       [DEPRECATED] Specifies whether or not the associated environment will be the default

```

`associate` is the entry point of the cli. You need to associate an environment before you can run most other commands. Check out [scope](#global-scope) and [aliases](#environment-aliases) for more info on the value of the alias and default options. Here is a sample command

```
catalyze associate My-Production-Environment app01 -a prod
```

#  Associated

```

Usage:  associated

Lists all associated environments

```

`associated` outputs information about all previously associated environments on your local machine. The information that is printed out includes the alias, environment ID, actual environment name, service ID, and the git repo directory. Here is a sample command

```
catalyze associated
```

#  Certs

The `certs` command gives access to certificate and private key management for public facing services. The certs command cannot be run directly but has sub commands.

##  Certs Create

```

Usage:  certs create HOSTNAME PUBLIC_KEY_PATH PRIVATE_KEY_PATH [-s] [-r]

Create a new domain with an SSL certificate and private key

Arguments:
  HOSTNAME=""           The hostname of this domain and SSL certificate plus private key pair
  PUBLIC_KEY_PATH=""    The path to a public key file in PEM format
  PRIVATE_KEY_PATH=""   The path to an unencrypted private key file in PEM format

Options:
  -s, --self-signed=false   Whether or not the given SSL certificate and private key are self signed
  -r, --resolve=true        Whether or not to attempt to automatically resolve incomplete SSL certificate issues

```



##  Certs List

```

Usage:  certs list

List all existing domains that have SSL certificate and private key pairs

```



##  Certs Rm

```

Usage:  certs rm HOSTNAME

Remove an existing domain and its associated SSL certificate and private key pair

Arguments:
  HOSTNAME=""   The hostname of the domain and SSL certificate and private key pair

```



##  Certs Update

```

Usage:  certs update HOSTNAME PUBLIC_KEY_PATH PRIVATE_KEY_PATH [-s] [-r]

Update the SSL certificate and private key pair for an existing domain

Arguments:
  HOSTNAME=""           The hostname of this domain and SSL certificate and private key pair
  PUBLIC_KEY_PATH=""    The path to a public key file in PEM format
  PRIVATE_KEY_PATH=""   The path to an unencrypted private key file in PEM format

Options:
  -s, --self-signed=false   Whether or not the given SSL certificate and private key are self signed
  -r, --resolve=true        Whether or not to attempt to automatically resolve incomplete SSL certificate issues

```



#  Clear

```

Usage:  clear [--private-key] [--session] [--environments] [--default] [--pods] [--all]

Clear out information in the global settings file to fix a misconfigured CLI.

Options:
  --private-key=false    Clear out the saved private key information
  --session=false        Clear out all session information
  --environments=false   Clear out all associated environments
  --default=false        [DEPRECATED] Clear out the saved default environment
  --pods=false           Clear out all saved pods
  --all=false            Clear out all settings

```

`clear` allows you to manage your global settings file in case your CLI becomes misconfigured. The global settings file is stored in your home directory at `~/.catalyze`. You can clear out all settings or pick and choose which ones need to be removed. After running the `clear` command, any other CLI command will reset the removed settings to their appropriate values. Here are some sample commands

```
catalyze clear --all
catalyze clear --environments # removes your associated environments
catalyze clear --session --private-key # removes all session and private key authentication information
```

#  Console

```

Usage:  console SERVICE_NAME [COMMAND]

Open a secure console to a service

Arguments:
  SERVICE_NAME=""   The name of the service to open up a console for
  COMMAND=""        An optional command to run when the console becomes available

```

`console` gives you direct access to your database service or application shell. For example, if you open up a console to a postgres database, you will be given access to a psql prompt. You can also open up a mysql prompt, mongo cli prompt, rails console, django shell, and much more. When accessing a database service, the `COMMAND` argument is not needed because the appropriate prompt will be given to you. If you are connecting to an application service the `COMMAND` argument is required. Here are some sample commands

```
catalyze console db01
catalyze console app01 "bundle exec rails console"
```

#  Dashboard

```

Usage:  dashboard

Open the Catalyze Dashboard in your default browser

```

`dashboard` opens up the Catalyze Dashboard homepage in your default web browser. Here is a sample command

```
catalyze dashboard
```

#  Db

The `db` command gives access to backup, import, and export services for databases. The db command can not be run directly but has sub commands.

##  Db Backup

```

Usage:  db backup DATABASE_NAME [-s]

Create a new backup

Arguments:
  DATABASE_NAME=""   The name of the database service to create a backup for (i.e. 'db01')

Options:
  -s, --skip-poll=false   Whether or not to wait for the backup to finish

```



##  Db Download

```

Usage:  db download DATABASE_NAME BACKUP_ID FILEPATH [-f]

Download a previously created backup

Arguments:
  DATABASE_NAME=""   The name of the database service which was backed up (i.e. 'db01')
  BACKUP_ID=""       The ID of the backup to download (found from "catalyze backup list")
  FILEPATH=""        The location to save the downloaded backup to. This location must NOT already exist unless -f is specified

Options:
  -f, --force=false   If a file previously exists at "filepath", overwrite it and download the backup

```



##  Db Export

```

Usage:  db export DATABASE_NAME FILEPATH [-f]

Export data from a database

Arguments:
  DATABASE_NAME=""   The name of the database to export data from (i.e. 'db01')
  FILEPATH=""        The location to save the exported data. This location must NOT already exist unless -f is specified

Options:
  -f, --force=false   If a file previously exists at `filepath`, overwrite it and export data

```



##  Db Import

```

Usage:  db import DATABASE_NAME FILEPATH [-d [-c]]

Import data into a database

Arguments:
  DATABASE_NAME=""   The name of the database to import data to (i.e. 'db01')
  FILEPATH=""        The location of the file to import to the database

Options:
  -c, --mongo-collection=""   If importing into a mongo service, the name of the collection to import into
  -d, --mongo-database=""     If importing into a mongo service, the name of the database to import into

```



##  Db List

```

Usage:  db list DATABASE_NAME [-p] [-n]

List created backups

Arguments:
  DATABASE_NAME=""   The name of the database service to list backups for (i.e. 'db01')

Options:
  -p, --page=1         The page to view
  -n, --page-size=10   The number of items to show per page

```



##  Db Logs

```

Usage:  db logs DATABASE_NAME BACKUP_ID

Print out the logs from a previous database backup job

Arguments:
  DATABASE_NAME=""   The name of the database service (i.e. 'db01')
  BACKUP_ID=""       The ID of the backup to download logs from (found from "catalyze backup list")

```



#  Default

```

Usage:  default ENV_ALIAS

[DEPRECATED] Set the default associated environment

Arguments:
  ENV_ALIAS=""   The alias of an already associated environment to set as the default

```

The `default` command has been deprecated! It will be removed in a future version. Please specify `-E` on all commands instead of using the default.

`default` sets the default environment for all commands that don't specify an environment with the `-E` flag. See [scope](#global-scope) for more information on scope and default environments. When setting a default environment, you must give the alias of the environment if one was set when it was associated and not the real environment name. Here is a sample command

```
catalyze default prod
```

#  Deploy-keys

The `deploy-keys` command gives access to SSH deploy keys for environment services. The deploy-keys command can not be run directly but has sub commands.

##  Deploy-keys Add

```

Usage:  deploy-keys add NAME KEY_PATH SERVICE_NAME

Add a new deploy key

Arguments:
  NAME=""           The name for the new key, for your own purposes
  KEY_PATH=""       Relative path to the SSH key file
  SERVICE_NAME=""   The name of the code service to add this deploy key to

```



##  Deploy-keys List

```

Usage:  deploy-keys list SERVICE_NAME

List all deploy keys

Arguments:
  SERVICE_NAME=""   The name of the code service to list deploy keys

```



##  Deploy-keys Rm

```

Usage:  deploy-keys rm NAME SERVICE_NAME

Remove a deploy key

Arguments:
  NAME=""           The name of the key to remove
  SERVICE_NAME=""   The name of the code service to remove this deploy key from

```



#  Disassociate

```

Usage:  disassociate ENV_ALIAS

Remove the association with an environment

Arguments:
  ENV_ALIAS=""   The alias of an already associated environment to disassociate

```

`disassociate` removes the environment from your list of associated environments but **does not** remove the catalyze git remote on the git repo. Disassociate does not have to be run from within a git repo. Here is a sample command

```
catalyze disassociate myprod
```

#  Domain

```

Usage:  domain

Print out the temporary domain name of the environment

```

`domain` prints out the temporary domain name setup by Catalyze for an environment. This domain name typically takes the form podXXXXX.catalyzeapps.com but may vary based on the environment. Here is a sample command

```
catalyze domain
```

#  Environments

This command has been moved! Please use [environments list](#environments-list) instead. This alias will be removed in the next CLI update.

The `environments` command allows you to manage your environments. The environments command can not be run directly but has sub commands.

##  Environments List

```

Usage:  environments list

List all environments you have access to

```



##  Environments Rename

```

Usage:  environments rename NAME

Rename an environment

Arguments:
  NAME=""      The new name of the environment

```



#  Files

The `files` command gives access to service files on your environment's services. Service files can include Nginx configs, SSL certificates, and any other file that might be injected into your running service. The files command can not be run directly but has sub commands.

##  Files Download

```

Usage:  files download [SERVICE_NAME] FILE_NAME [-o] [-f]

Download a file to your localhost with the same file permissions as on the remote host or print it to stdout

Arguments:
  SERVICE_NAME="service_proxy"   The name of the service to download a file from
  FILE_NAME=""                   The name of the service file from running "catalyze files list"

Options:
  -o, --output=""     The downloaded file will be saved to the given location with the same file permissions as it has on the remote host. If those file permissions cannot be applied, a warning will be printed and default 0644 permissions applied. If no output is specified, stdout is used.
  -f, --force=false   If the specified output file already exists, automatically overwrite it

```



##  Files List

```

Usage:  files list [SERVICE_NAME]

List all files available for a given service

Arguments:
  SERVICE_NAME="service_proxy"   The name of the service to list files for

```



#  Git-remote

The `git-remote` command allows you to interact with code service remote git URLs. The git-remote command can not be run directly but has sub commands.

##  Git-remote Add

```

Usage:  git-remote add SERVICE_NAME [-r]

Add the git remote for the given code service to the local git repo

Arguments:
  SERVICE_NAME=""   The name of the service to add a git remote for

Options:
  -r, --remote="catalyze"   The name of the git remote to be added

```



##  Git-remote Show

```

Usage:  git-remote show SERVICE_NAME

Print out the git remote for a given code service

Arguments:
  SERVICE_NAME=""   The name of the service to add a git remote for

```



#  Invites

The `invites` command gives access to organization invitations. Every environment is owned by an organization and users join organizations in order to access individual environments. You can invite new users by email and manage pending invites through the CLI. You cannot call the `invites` command directly, but must call one of its subcommands.

##  Invites Accept

```

Usage:  invites accept INVITE_CODE

Accept an organization invite

Arguments:
  INVITE_CODE=""   The invite code that was sent in the invite email

```



##  Invites List

```

Usage:  invites list

List all pending organization invitations

```



##  Invites Rm

```

Usage:  invites rm INVITE_ID

Remove a pending organization invitation

Arguments:
  INVITE_ID=""   The ID of an invitation to remove

```



##  Invites Send

```

Usage:  invites send EMAIL [-m | -a]

Send an invite to a user by email for a given organization

Arguments:
  EMAIL=""     The email of a user to invite to the associated environment. This user does not need to have a Catalyze account prior to sending the invitation

Options:
  -m, --member=true   Whether or not the user will be invited as a basic member
  -a, --admin=false   Whether or not the user will be invited as an admin

```



#  Keys

The `keys` command gives access to SSH key management for your user account. SSH keys can be used for authentication and pushing code to the Catalyze platform. Any SSH keys added to your user account should not be shared but be treated as private SSH keys. Any SSH key uploaded to your user account will be able to be used with all code services and environments that you have access to. The keys command can not be run directly but has sub commands.

##  Keys Add

```

Usage:  keys add NAME PUBLIC_KEY_PATH

Add a public key

Arguments:
  NAME=""              The name for the new key, for your own purposes
  PUBLIC_KEY_PATH=""   Relative path to the public key file

```



##  Keys List

```

Usage:  keys list

List your public keys

```



##  Keys Rm

```

Usage:  keys rm NAME

Remove a public key

Arguments:
  NAME=""      The name of the key to remove.

```



##  Keys Set

```

Usage:  keys set PRIVATE_KEY_PATH

Set your auth key

Arguments:
  PRIVATE_KEY_PATH=""   Relative path to the private key file

```



#  Logout

```

Usage:  logout

Clear the stored user information from your local machine

```

When using the CLI, your username and password are **never** stored in any file on your filesystem. However, in order to not type in your username and password each and every command, a session token is stored in the CLI's configuration file and used until it expires. `logout` removes this session token from the configuration file. Here is a sample command

```
catalyze logout
```

#  Logs

```

Usage:  logs [QUERY] [(-f | -t)] [--hours] [--minutes] [--seconds]

Show the logs in your terminal streamed from your logging dashboard

Arguments:
  QUERY="*"    The query to send to your logging dashboard's elastic search (regex is supported)

Options:
  -f, --follow=false   Tail/follow the logs (Equivalent to -t)
  -t, --tail=false     Tail/follow the logs (Equivalent to -f)
  --hours=0            The number of hours before now (in combination with minutes and seconds) to retrieve logs
  --minutes=0          The number of minutes before now (in combination with hours and seconds) to retrieve logs
  --seconds=0          The number of seconds before now (in combination with hours and minutes) to retrieve logs

```

`logs` prints out your application logs directly from your logging Dashboard. If you do not see your logs, try adjusting the number of hours, minutes, or seconds of logs that are retrieved with the `--hours`, `--minutes`, and `--seconds` options respectively. You can also follow the logs with the `-f` option. When using `-f` all logs will be printed to the console within the given time frame as well as any new logs that are sent to the logging Dashboard for the duration of the command. When using the `-f` option, hit ctrl-c to stop. Here are some sample commands

```
catalyze logs --hours=6 --minutes=30
catalyze logs -f
```

#  Metrics

The `metrics` command gives access to environment metrics or individual service metrics through a variety of formats. This is useful for checking on the status and performance of your application or environment as a whole. The metrics command cannot be run directly but has sub commands.

##  Metrics Cpu

```

Usage:  metrics cpu [SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]

Print service and environment CPU metrics in your local time zone

Arguments:
  SERVICE_NAME=""   The name of the service to print metrics for

Options:
  --json=false     Output the data as json
  --csv=false      Output the data as csv
  --spark=false    Output the data using spark lines
  --stream=false   Repeat calls once per minute until this process is interrupted.
  -m, --mins=1     How many minutes worth of metrics to retrieve.

```



##  Metrics Memory

```

Usage:  metrics memory [SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]

Print service and environment memory metrics in your local time zone

Arguments:
  SERVICE_NAME=""   The name of the service to print metrics for

Options:
  --json=false     Output the data as json
  --csv=false      Output the data as csv
  --spark=false    Output the data using spark lines
  --stream=false   Repeat calls once per minute until this process is interrupted.
  -m, --mins=1     How many minutes worth of metrics to retrieve.

```



##  Metrics Network-in

```

Usage:  metrics network-in [SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]

Print service and environment received network data metrics in your local time zone

Arguments:
  SERVICE_NAME=""   The name of the service to print metrics for

Options:
  --json=false     Output the data as json
  --csv=false      Output the data as csv
  --spark=false    Output the data using spark lines
  --stream=false   Repeat calls once per minute until this process is interrupted.
  -m, --mins=1     How many minutes worth of metrics to retrieve.

```



##  Metrics Network-out

```

Usage:  metrics network-out [SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]

Print service and environment transmitted network data metrics in your local time zone

Arguments:
  SERVICE_NAME=""   The name of the service to print metrics for

Options:
  --json=false     Output the data as json
  --csv=false      Output the data as csv
  --spark=false    Output the data using spark lines
  --stream=false   Repeat calls once per minute until this process is interrupted.
  -m, --mins=1     How many minutes worth of metrics to retrieve.

```



#  Rake

```

Usage:  rake [SERVICE_NAME] TASK_NAME

Execute a rake task

Arguments:
  SERVICE_NAME=""   The service that will run the rake task. Defaults to the associated service.
  TASK_NAME=""      The name of the rake task to run

```

`rake` executes a rake task by its name asynchronously. Once executed, the output of the task can be seen through your logging Dashboard. Here is a sample command

```
catalyze rake code-1 db:migrate
```

#  Redeploy

```

Usage:  redeploy SERVICE_NAME

Redeploy a service without having to do a git push

Arguments:
  SERVICE_NAME=""   The name of the service to redeploy (i.e. 'app01')

```

`redeploy` deploys an identical copy of the given service. For code services, this avoids having to perform a code push. You skip the git push and the build. For service proxies, new instances simply replace the old ones. All other service types cannot be redeployed with this command. Here is a sample command

```
catalyze redeploy app01
```

#  Releases

The `releases` command allows you to manage your code service releases. A release is automatically created each time you perform a git push. The release is tagged with the git SHA of the commit. Releases are a way of tagging specific points in time of your git history. You can rollback to a specific release by using the [rollback](#rollback) command. The releases command cannot be run directly but has sub commands.

##  Releases List

```

Usage:  releases list SERVICE_NAME

List all releases for a given code service

Arguments:
  SERVICE_NAME=""   The name of the service to list releases for

```



##  Releases Rm

```

Usage:  releases rm SERVICE_NAME RELEASE_NAME

Remove a release from a code service

Arguments:
  SERVICE_NAME=""   The name of the service to remove a release from
  RELEASE_NAME=""   The name of the release to remove

```



##  Releases Update

```

Usage:  releases update SERVICE_NAME RELEASE_NAME [--notes] [--release]

Update a release from a code service

Arguments:
  SERVICE_NAME=""   The name of the service to update a release for
  RELEASE_NAME=""   The name of the release to update

Options:
  -n, --notes=""     The new notes to save on the release. If omitted, notes will be unchanged.
  -r, --release=""   The new name of the release. If omitted, the release name will be unchanged.

```



#  Rollback

```

Usage:  rollback SERVICE_NAME RELEASE_NAME

Rollback a code service to a specific release

Arguments:
  SERVICE_NAME=""   The name of the service to rollback
  RELEASE_NAME=""   The name of the release to rollback to

```

`rollback` is a way to redeploy older versions of your code service. You must specify the name of the service to rollback and the name of an existing release to rollback to. Releases can be found with the [releases list](#releases-list) command. Here are some sample commands

```
catalyze rollback code-1 f93ced037f828dcaabccfc825e6d8d32cc5a1883
```

#  Services

The `services` command allows you to manage your services. The services command cannot be run directly but has sub commands.

##  Services List

```

Usage:  services list

List all services for your environment

```



##  Services Stop

```

Usage:  services stop SERVICE_NAME

Stop all instances of a given service (including all workers, rake tasks, and open consoles)

Arguments:
  SERVICE_NAME=""   The name of the service to stop

```



##  Services Rename

```

Usage:  services rename SERVICE_NAME NEW_NAME

Rename a service

Arguments:
  SERVICE_NAME=""   The service to rename
  NEW_NAME=""       The new name for the service

```



#  Sites

The `sites` command gives access to hostname and SSL certificate usage for public facing services. `sites` are different from `certs` in that `sites` use an instance of a `cert` and are associated with a single service. `certs` can be used by multiple sites. The sites command can not be run directly but has sub commands.

##  Sites Create

```

Usage:  sites create SITE_NAME SERVICE_NAME HOSTNAME [--client-max-body-size] [--proxy-connect-timeout] [--proxy-read-timeout] [--proxy-send-timeout] [--proxy-upstream-timeout] [--enable-cors] [--enable-websockets]

Create a new site linking it to an existing cert instance

Arguments:
  SITE_NAME=""      The name of the site to be created. This will be used in this site's nginx configuration file (i.e. ".example.com")
  SERVICE_NAME=""   The name of the service to add this site configuration to (i.e. 'app01')
  HOSTNAME=""       The hostname used in the creation of a certs instance with the 'certs' command (i.e. "star_example_com")

Options:
  --client-max-body-size=-1     The 'client_max_body_size' nginx config specified in megabytes
  --proxy-connect-timeout=-1    The 'proxy_connect_timeout' nginx config specified in seconds
  --proxy-read-timeout=-1       The 'proxy_read_timeout' nginx config specified in seconds
  --proxy-send-timeout=-1       The 'proxy_send_timeout' nginx config specified in seconds
  --proxy-upstream-timeout=-1   The 'proxy_next_upstream_timeout' nginx config specified in seconds
  --enable-cors=false           Enable or disable all features related to full CORS support
  --enable-websockets=false     Enable or disable all features related to full websockets support

```



##  Sites List

```

Usage:  sites list

List details for all site configurations

```



##  Sites Rm

```

Usage:  sites rm NAME

Remove a site configuration

Arguments:
  NAME=""      The name of the site configuration to delete

```



##  Sites Show

```

Usage:  sites show NAME

Shows the details for a given site

Arguments:
  NAME=""      The name of the site configuration to show

```



#  Ssl

The `ssl` command offers access to subcommands that deal with SSL certificates. You cannot run the SSL command directly but must call a subcommand.

##  Ssl Resolve

```

Usage:  ssl resolve CHAIN PRIVATE_KEY HOSTNAME [OUTPUT] [-f]

Verify that an SSL certificate is signed by a valid CA and attempt to resolve any incomplete certificate chains that are found

Arguments:
  CHAIN=""         The path to your full certificate chain in PEM format
  PRIVATE_KEY=""   The path to your private key in PEM format
  HOSTNAME=""      The hostname that should match your certificate (i.e. "*.catalyze.io")
  OUTPUT=""        The path of a file to save your properly resolved certificate chain (defaults to STDOUT)

Options:
  -f, --force=false   If an output file is specified and already exists, setting force to true will overwrite the existing output file

```



##  Ssl Verify

```

Usage:  ssl verify CHAIN PRIVATE_KEY HOSTNAME [-s]

Verify whether a certificate chain is complete and if it matches the given private key

Arguments:
  CHAIN=""         The path to your full certificate chain in PEM format
  PRIVATE_KEY=""   The path to your private key in PEM format
  HOSTNAME=""      The hostname that should match your certificate (i.e. "*.catalyze.io")

Options:
  -s, --self-signed=false   Whether or not the certificate is self signed. If set, chain verification is skipped

```



#  Status

```

Usage:  status

Get quick readout of the current status of your associated environment and all of its services

```

`status` will give a quick readout of your environment's health. This includes your environment name, environment ID, and for each service the name, size, build status, deploy status, and service ID. Here is a sample command

```
catalyze status
```

#  Support-ids

```

Usage:  support-ids

Print out various IDs related to your associated environment to be used when contacting Catalyze support

```

`support-ids` is helpful when contacting Catalyze support by sending an email to support@catalyze.io. If you are having an issue with a CLI command or anything with your environment, it is helpful to run this command and copy the output into the initial correspondence with a Catalyze engineer. This will help Catalyze identify the environment faster and help come to resolution faster. Here is a sample command

```
catalyze support-ids
```

#  Update

```

Usage:  update

Checks for available updates and updates the CLI if a new update is available

```

`update` is a shortcut to update your CLI instantly. If a newer version of the CLI is available, it will be downloaded and installed automatically. This is used when you want to apply an update before the CLI automatically applies it on its own. Here is a sample command

```
catalyze update
```

#  Users

The `users` command allows you to manage who has access to your environment through the organization that owns the environment. The users command can not be run directly but has three sub commands.

##  Users List

```

Usage:  users list

List all users who have access to the given organization

```



##  Users Rm

```

Usage:  users rm EMAIL

Revoke access to the given organization for the given user

Arguments:
  EMAIL=""     The email address of the user to revoke access from for the given organization

```



#  Vars

The `vars` command allows you to manage environment variables for your code services. The vars command can not be run directly but has sub commands.

##  Vars List

```

Usage:  vars list [SERVICE_NAME] [--json | --yaml]

List all environment variables

Arguments:
  SERVICE_NAME=""   The name of the service containing the environment variables. Defaults to the associated service.

Options:
  --json=false   Output environment variables in JSON format
  --yaml=false   Output environment variables in YAML format

```



##  Vars Set

```

Usage:  vars set [SERVICE_NAME] -v...

Set one or more new environment variables or update the values of existing ones

Arguments:
  SERVICE_NAME=""   The name of the service on which the environment variables will be set. Defaults to the associated service.

Options:
  -v, --variable    The env variable to set or update in the form "<key>=<value>"

```



##  Vars Unset

```

Usage:  vars unset [SERVICE_NAME] VARIABLE

Unset (delete) an existing environment variable

Arguments:
  SERVICE_NAME=""   The name of the service on which the environment variables will be unset. Defaults to the associated service.
  VARIABLE=""       The name of the environment variable to unset

```



#  Version

```

Usage:  version

Output the version and quit

```

`version` prints out the current CLI version as well as the architecture it was built for (64-bit or 32-bit). This is useful to see if you have the latest version of the CLI and when working with Catalyze support engineers to ensure you have the correct CLI installed. Here is a sample command

```
catalyze version
```

#  Whoami

```

Usage:  whoami

Retrieve your user ID

```

`whoami` prints out the currently logged in user's users ID. This is used with Catalyze support engineers. Here is a sample command

```
catalyze whoami
```

#  Worker

This command has been moved! Please use [worker deploy](#worker-deploy) instead. This alias will be removed in the next CLI update.

The `worker` commands allow you to manage your environment variables per service. The `worker` command cannot be run directly, but has subcommands.

##  Worker Deploy

```

Usage:  worker deploy SERVICE_NAME TARGET

Deploy new workers for a given service

Arguments:
  SERVICE_NAME=""   The name of the service to use to deploy a worker
  TARGET=""         The name of the Procfile target to invoke as a worker

```



##  Worker List

```

Usage:  worker list SERVICE_NAME

Lists all workers for a given service

Arguments:
  SERVICE_NAME=""   The name of the service to list workers for

```



##  Worker Rm

```

Usage:  worker rm SERVICE_NAME TARGET

Remove all workers for a given service and target

Arguments:
  SERVICE_NAME=""   The name of the service running the workers
  TARGET=""         The worker target to remove

```



##  Worker Scale

```

Usage:  worker scale SERVICE_NAME TARGET SCALE

Scale existing workers up or down for a given service and target

Arguments:
  SERVICE_NAME=""   The name of the service running the workers
  TARGET=""         The worker target to scale up or down
  SCALE=""          The new scale (or change in scale) for the given worker target. This can be a single value (i.e. 2) representing the final number of workers that should be running. Or this can be a change represented by a plus or minus sign followed by the value (i.e. +2 or -1). When using a change in value, be sure to insert the "--" operator to signal the end of options. For example, "catalyze worker scale code-1 worker -- -1"

```



 Here is a sample command

```
catalyze support-ids
```

#  Update

```

Usage:  update

Checks for available updates and updates the CLI if a new update is available

```

`update` is a shortcut to update your CLI instantly. If a newer version of the CLI is available, it will be downloaded and installed automatically. This is used when you want to apply an update before the CLI automatically applies it on its own. Here is a sample command

```
catalyze update
```

#  Users

```

Usage:  users COMMAND [arg...]

Manage users who have access to the given organization

Commands:
  list	List all users who have access to the given organization
  rm	Revoke access to the given organization for the given user

Run ' users COMMAND --help' for more information on a command.

```

The `users` command allows you to manage who has access to your environment through the organization that owns the environment. The users command can not be run directly but has three sub commands.

##  Users List

```

Usage:  users list

List all users who have access to the given organization

```



##  Users Rm

```

Usage:  users rm EMAIL

Revoke access to the given organization for the given user

Arguments:
  EMAIL=""     The email address of the user to revoke access from for the given organization

```



#  Vars

```

Usage:  vars COMMAND [arg...]

Interaction with environment variables for the associated environment

Commands:
  list	List all environment variables
  set	Set one or more new environment variables or update the values of existing ones
  unset	Unset (delete) an existing environment variable

Run ' vars COMMAND --help' for more information on a command.

```

The `vars` command allows you to manage environment variables for your code services. The vars command can not be run directly but has sub commands.

##  Vars List

```

Usage:  vars list [SERVICE_NAME] [--json | --yaml]

List all environment variables

Arguments:
  SERVICE_NAME=""   The name of the service containing the environment variables. Defaults to the associated service.

Options:
  --json=false   Output environment variables in JSON format
  --yaml=false   Output environment variables in YAML format

```



##  Vars Set

```

Usage:  vars set [SERVICE_NAME] -v...

Set one or more new environment variables or update the values of existing ones

Arguments:
  SERVICE_NAME=""   The name of the service on which the environment variables will be set. Defaults to the associated service.

Options:
  -v, --variable    The env variable to set or update in the form "<key>=<value>"

```



##  Vars Unset

```

Usage:  vars unset [SERVICE_NAME] VARIABLE

Unset (delete) an existing environment variable

Arguments:
  SERVICE_NAME=""   The name of the service on which the environment variables will be unset. Defaults to the associated service.
  VARIABLE=""       The name of the environment variable to unset

```



#  Version

```

Usage:  version

Output the version and quit

```

`version` prints out the current CLI version as well as the architecture it was built for (64-bit or 32-bit). This is useful to see if you have the latest version of the CLI and when working with Catalyze support engineers to ensure you have the correct CLI installed. Here is a sample command

```
catalyze version
```

#  Whoami

```

Usage:  whoami

Retrieve your user ID

```

`whoami` prints out the currently logged in user's users ID. This is used with Catalyze support engineers. Here is a sample command

```
catalyze whoami
```

#  Worker

```

Usage:  worker [SERVICE_NAME] [TARGET] COMMAND [arg...]

Manage a service's workers

Arguments:
  SERVICE_NAME=""   The name of the service to use to start a worker. Defaults to the associated service.
  TARGET=""         The name of the Procfile target to invoke as a worker

Commands:
  deploy	Deploy new workers for a given service
  list	Lists all workers for a given service
  rm	Remove all workers for a given service and target
  scale	Scale existing workers up or down for a given service and target

Run ' worker COMMAND --help' for more information on a command.

```

This command has been moved! Please use [worker deploy](#worker-deploy) instead. This alias will be removed in the next CLI update.

The `worker` commands allow you to manage your environment variables per service. The `worker` command cannot be run directly, but has subcommands.

##  Worker Deploy

```

Usage:  worker deploy SERVICE_NAME TARGET

Deploy new workers for a given service

Arguments:
  SERVICE_NAME=""   The name of the service to use to deploy a worker
  TARGET=""         The name of the Procfile target to invoke as a worker

```



##  Worker List

```

Usage:  worker list SERVICE_NAME

Lists all workers for a given service

Arguments:
  SERVICE_NAME=""   The name of the service to list workers for

```



##  Worker Rm

```

Usage:  worker rm SERVICE_NAME TARGET

Remove all workers for a given service and target

Arguments:
  SERVICE_NAME=""   The name of the service running the workers
  TARGET=""         The worker target to remove

```



##  Worker Scale

```

Usage:  worker scale SERVICE_NAME TARGET SCALE

Scale existing workers up or down for a given service and target

Arguments:
  SERVICE_NAME=""   The name of the service running the workers
  TARGET=""         The worker target to scale up or down
  SCALE=""          The new scale (or change in scale) for the given worker target. This can be a single value (i.e. 2) representing the final number of workers that should be running. Or this can be a change represented by a plus or minus sign followed by the value (i.e. +2 or -1). When using a change in value, be sure to insert the "--" operator to signal the end of options. For example, "catalyze worker scale code-1 worker -- -1"

```



CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Db Logs

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Default

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `default` command has been deprecated! It will be removed in a future version. Please specify `-E` on all commands instead of using the default.

`default` sets the default environment for all commands that don't specify an environment with the `-E` flag. See [scope](#global-scope) for more information on scope and default environments. When setting a default environment, you must give the alias of the environment if one was set when it was associated and not the real environment name. Here is a sample command

```
catalyze default prod
```

#  Deploy-keys

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `deploy-keys` command gives access to SSH deploy keys for environment services. The deploy-keys command can not be run directly but has sub commands.

##  Deploy-keys Add

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Deploy-keys List

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Deploy-keys Rm

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Disassociate

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`disassociate` removes the environment from your list of associated environments but **does not** remove the catalyze git remote on the git repo. Disassociate does not have to be run from within a git repo. Here is a sample command

```
catalyze disassociate myprod
```

#  Domain

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`domain` prints out the temporary domain name setup by Catalyze for an environment. This domain name typically takes the form podXXXXX.catalyzeapps.com but may vary based on the environment. Here is a sample command

```
catalyze domain
```

#  Environments

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

This command has been moved! Please use [environments list](#environments-list) instead. This alias will be removed in the next CLI update.

The `environments` command allows you to manage your environments. The environments command can not be run directly but has sub commands.

##  Environments List

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Environments Rename

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Files

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `files` command gives access to service files on your environment's services. Service files can include Nginx configs, SSL certificates, and any other file that might be injected into your running service. The files command can not be run directly but has sub commands.

##  Files Download

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Files List

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Git-remote

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `git-remote` command allows you to interact with code service remote git URLs. The git-remote command can not be run directly but has sub commands.

##  Git-remote Add

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Git-remote Show

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Invites

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `invites` command gives access to organization invitations. Every environment is owned by an organization and users join organizations in order to access individual environments. You can invite new users by email and manage pending invites through the CLI. You cannot call the `invites` command directly, but must call one of its subcommands.

##  Invites Accept

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Invites List

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Invites Rm

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Invites Send

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Keys

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `keys` command gives access to SSH key management for your user account. SSH keys can be used for authentication and pushing code to the Catalyze platform. Any SSH keys added to your user account should not be shared but be treated as private SSH keys. Any SSH key uploaded to your user account will be able to be used with all code services and environments that you have access to. The keys command can not be run directly but has sub commands.

##  Keys Add

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Keys List

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Keys Rm

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Keys Set

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Logout

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

When using the CLI, your username and password are **never** stored in any file on your filesystem. However, in order to not type in your username and password each and every command, a session token is stored in the CLI's configuration file and used until it expires. `logout` removes this session token from the configuration file. Here is a sample command

```
catalyze logout
```

#  Logs

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`logs` prints out your application logs directly from your logging Dashboard. If you do not see your logs, try adjusting the number of hours, minutes, or seconds of logs that are retrieved with the `--hours`, `--minutes`, and `--seconds` options respectively. You can also follow the logs with the `-f` option. When using `-f` all logs will be printed to the console within the given time frame as well as any new logs that are sent to the logging Dashboard for the duration of the command. When using the `-f` option, hit ctrl-c to stop. Here are some sample commands

```
catalyze logs --hours=6 --minutes=30
catalyze logs -f
```

#  Metrics

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `metrics` command gives access to environment metrics or individual service metrics through a variety of formats. This is useful for checking on the status and performance of your application or environment as a whole. The metrics command cannot be run directly but has sub commands.

##  Metrics Cpu

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Metrics Memory

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Metrics Network-in

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Metrics Network-out

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Rake

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`rake` executes a rake task by its name asynchronously. Once executed, the output of the task can be seen through your logging Dashboard. Here is a sample command

```
catalyze rake code-1 db:migrate
```

#  Redeploy

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`redeploy` deploys an identical copy of the given service. For code services, this avoids having to perform a code push. You skip the git push and the build. For service proxies, new instances simply replace the old ones. All other service types cannot be redeployed with this command. Here is a sample command

```
catalyze redeploy app01
```

#  Releases

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `releases` command allows you to manage your code service releases. A release is automatically created each time you perform a git push. The release is tagged with the git SHA of the commit. Releases are a way of tagging specific points in time of your git history. You can rollback to a specific release by using the [rollback](#rollback) command. The releases command cannot be run directly but has sub commands.

##  Releases List

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Releases Rm

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Releases Update

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Rollback

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`rollback` is a way to redeploy older versions of your code service. You must specify the name of the service to rollback and the name of an existing release to rollback to. Releases can be found with the [releases list](#releases-list) command. Here are some sample commands

```
catalyze rollback code-1 f93ced037f828dcaabccfc825e6d8d32cc5a1883
```

#  Services

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `services` command allows you to manage your services. The services command cannot be run directly but has sub commands.

##  Services List

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Services Stop

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Services Rename

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Sites

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `sites` command gives access to hostname and SSL certificate usage for public facing services. `sites` are different from `certs` in that `sites` use an instance of a `cert` and are associated with a single service. `certs` can be used by multiple sites. The sites command can not be run directly but has sub commands.

##  Sites Create

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Sites List

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Sites Rm

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Sites Show

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Ssl

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `ssl` command offers access to subcommands that deal with SSL certificates. You cannot run the SSL command directly but must call a subcommand.

##  Ssl Resolve

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Ssl Verify

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Status

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`status` will give a quick readout of your environment's health. This includes your environment name, environment ID, and for each service the name, size, build status, deploy status, and service ID. Here is a sample command

```
catalyze status
```

#  Support-ids

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`support-ids` is helpful when contacting Catalyze support by sending an email to support@catalyze.io. If you are having an issue with a CLI command or anything with your environment, it is helpful to run this command and copy the output into the initial correspondence with a Catalyze engineer. This will help Catalyze identify the environment faster and help come to resolution faster. Here is a sample command

```
catalyze support-ids
```

#  Update

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`update` is a shortcut to update your CLI instantly. If a newer version of the CLI is available, it will be downloaded and installed automatically. This is used when you want to apply an update before the CLI automatically applies it on its own. Here is a sample command

```
catalyze update
```

#  Users

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `users` command allows you to manage who has access to your environment through the organization that owns the environment. The users command can not be run directly but has three sub commands.

##  Users List

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Users Rm

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Vars

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `vars` command allows you to manage environment variables for your code services. The vars command can not be run directly but has sub commands.

##  Vars List

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Vars Set

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Vars Unset

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Version

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`version` prints out the current CLI version as well as the architecture it was built for (64-bit or 32-bit). This is useful to see if you have the latest version of the CLI and when working with Catalyze support engineers to ensure you have the correct CLI installed. Here is a sample command

```
catalyze version
```

#  Whoami

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`whoami` prints out the currently logged in user's users ID. This is used with Catalyze support engineers. Here is a sample command

```
catalyze whoami
```

#  Worker

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

This command has been moved! Please use [worker deploy](#worker-deploy) instead. This alias will be removed in the next CLI update.

The `worker` commands allow you to manage your environment variables per service. The `worker` command cannot be run directly, but has subcommands.

##  Worker Deploy

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Worker List

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Worker Rm

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Worker Scale

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.1.5

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



ificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



ns
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

The `vars` command allows you to manage environment variables for your code services. The vars command can not be run directly but has sub commands.

##  Vars List

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.3.0

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Vars Set

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.3.0

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Vars Unset

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.3.0

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



#  Version

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.3.0

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`version` prints out the current CLI version as well as the architecture it was built for (64-bit or 32-bit). This is useful to see if you have the latest version of the CLI and when working with Catalyze support engineers to ensure you have the correct CLI installed. Here is a sample command

```
catalyze version
```

#  Whoami

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.3.0

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

`whoami` prints out the currently logged in user's users ID. This is used with Catalyze support engineers. Here is a sample command

```
catalyze whoami
```

#  Worker

```
Error: incorrect usage

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.3.0

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```

This command has been moved! Please use [worker deploy](#worker-deploy) instead. This alias will be removed in the next CLI update.

The `worker` commands allow you to manage your environment variables per service. The `worker` command cannot be run directly, but has subcommands.

##  Worker Deploy

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.3.0

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Worker List

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.3.0

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Worker Rm

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.3.0

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



##  Worker Scale

```
Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 3.3.0

Options:
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit

Commands:
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



erts          Manage your SSL certificates and domains
  clear          Clear out information in the global settings file to fix a misconfigured CLI.
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        [DEPRECATED] Set the default associated environment
  deploy-keys    Tasks for SSH deploy keys
  disassociate   Remove the association with an environment
  domain         Print out the temporary domain name of the environment
  environments   Manage environments for which you have access
  files          Tasks for managing service files
  git-remote     Manage git remotes to Catalyze code services
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logout         Clear the stored user information from your local machine
  logs           Show the logs in your terminal streamed from your logging dashboard
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  releases       Manage releases for code services
  rollback       Rollback a code service to a specific release
  services       Perform operations on an environment's services
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  support-ids    Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Manage a service's workers
  version        Output the version and quit

Run 'catalyze COMMAND --help' for more information on a command.
```



