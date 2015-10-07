# Self-monitoring tool (sketch)

## Usage
Install package and dependencies:
```shell
go get github.com/leominov/self-monitoring
cd $GOPATH/src/github.com/leominov/self-monitoring
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

## Configuration
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
