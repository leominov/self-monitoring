# Self-monitoring tool (init script)

## Configuration

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
