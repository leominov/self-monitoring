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
    "interval": 15,
    "notifyAtStart": false,
	"processList":
    [
        "acrypt",
		"capella",
		"docker"
    ],
    "logger": true,
    "telegram": {
        "enable": true,
        "token": "",
        "contactID": 0,
        "debug": true
    }
}
```
