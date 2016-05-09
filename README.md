# Self-monitoring tool (sketch)

## How to Start

Install package and dependencies:

```shell
$ go get github.com/leominov/self-monitoring
```

Create own config file:

```shell
$ cp example.config.json config.json
```

Edit configuration:

```shell
$ vim config.json
```

Build binary:

```shell
$ make build
```

Run:

```shell
$ ./gomon
```

Or install as a service:

```shell
$ ./install.sh
```

## Special service commands

Reload configuration:

```shell
$ service gomon reload
```

## Telegram commands

```
sh - Exec shell command (shell)
service - Alias for /sh service (srv)
calc - Calculator (bc)
uptime - Server uptime (up)
status - Service list (st)
who - Who is logged in (w)
vote - Random vote (v)
version - Monitoring version (ver)
reload - Reload configuration (rld)
service-add - Add service to monitoring list (srvadd)
service-dev - Delete service from monitoring list (srvdel)
```

## Configuration

Example with description:

```javascript
{
    "nodeName": "local", // Alias for logs and messages
    "interval": "15s", // Update interval
    "notifyAtStart": false, // Send notification with offline services on start
	"processList": [ // Process list for monitoring
        "acrypt",
		"capella",
		"docker"
    ],
    "logLevel": "info", // Logger level (debug, info, warning, error, fatal, panic)
    "telegram": {
        "enable": false, // Enable Telegram notification
        "token": "", // Telegram Bot API Token
        "contactID": 0, // Contact ID (user, chat, etc.)
        "admins": [] // Admin list to control monitoring
    }
}
```

## Building release

### Environment variables

* `GITHUB_TOKEN` - GitHub token for uploading releases;
* `TELEGRAM_API_TOKEN` - Telegram token for release notice;
* `TELEGRAM_CHAT_ID` - Telegram Chat ID for release notice.

### Installing dependencies

```shell
$ go get github.com/aktau/github-release
```

### Building

Check release version:

```shell
$ vim VERSION
```

Build binaries:

```shell
$ make build-all
```

Prepare and upload GitHub release:

```shell
$ make upload-release
```
