# Self-monitoring tool (init script)

## Configuration Sys-v-init

Open init script. Set your workspace directory and config filename.

## Usage

Copy init script:

```shell
sudo cp gomon /etc/init.d/
```

Enable start service on a running system:

```shell
sudo rcconf
```

Run:

```shell
sudo init/gomon {start|stop|restart|status}
```


## Configuration Systemd unit

Open systemd unit. Set your workspace directory and config filename.

## Usage

Copy file script:

```shell
sudo cp gomon.service /etc/systemd/system/
```

Reload systemd units

```shell
sudo systemctl daemon-reload
```

Enable start service on a running system:

```shell
sudo systemctl enable gomon.service
```

Run:

```shell
sudo systemctl {start|stop|restart|status} gomon.service
```
