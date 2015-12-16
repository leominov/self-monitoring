# Self-monitoring tool (sketch)

## Usage
Install package and dependencies:
```shell
go get github.com/leominov/self-monitoring
```
Create own config file:
```shell
cp example.config.json config.json
```
Edit configuration:
```shell
vim config.json
```
Install:
```shell
sudo ./install-gomon.sh
```

## Special service commands
Reload configuration:
```shell
service gomon reload
```

## Configuration
Example with description:
```javascript
{
    "nodeName": "local", // Alias for logs and messages
    "interval": 15000, // Update interval in milliseconds
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
        "contactID": 0 // Contact ID (user, chat, etc.)
    }
}
```
