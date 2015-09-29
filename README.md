Self-monitoring tool (sketch)
=============================

Usage
-----
Create own config file:
```shell
cp example.config.json config.json
```
Edit configuration:
```shell
vim config.json
```
Run:
```shell
go run gomon.go
```
Or build:
```shell
go build gomon.go
```

Configuration
-------------
Example with description:
```javascript
{
    "interval": 15000, // Update interval in milliseconds
    "notifyAtStart": false, // Send notification with offline services on start
	"processList": [ // Process list for monitoring
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
