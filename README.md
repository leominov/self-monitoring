Self-monitoring tool (sketch)
=============================

Usage
-----
Run with 'config.json' config file:
```shell
go run *.go
```

Run with 'local.config.json' config file:
```shell
go run *.go --config=local.config.json
```

Configuration
-------------
Example:
```javascript
{
    "interval": 15, // Update interval in seconds
    "notifyAtStart": false, // Send notification with offline services on start
	"processList": // Process list for monitoring
    [
        "acrypt",
		"capella",
		"docker"
    ],
    "logger": true, // Print status info in log
    "telegram": {
        "enable": true, // Enable Telegram notification
        "token": "", // Telegram Bot API Token
        "contactID": 0, // Contact ID (user, chat, etc.)
        "debug": false // Print debug info
    }
}
```
