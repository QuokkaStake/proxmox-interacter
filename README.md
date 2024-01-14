# proxmox-interacter

![Latest release](https://img.shields.io/github/v/release/QuokkaS/proxmox-interacter)
[![Actions Status](https://github.com/QuokkaStake/proxmox-interacter/workflows/test/badge.svg)](https://github.com/QuokkaS/proxmox-interacter/actions)

proxmox-interacter is a tool to interact with your Proxmox instances via a Telegram bot.
Here's what it can do:
- List Proxmox instances
- List containers/VMs and show their details
- Start/stop/restart containers
- List Proxmox storages
- Scale Proxmox containers/VMs on a fly
- (planned) See SMART stats of your storage devices
- (planned) Backup VMs/containers

Why is it cool?
- Only a single binary required for it to work
- Fully open-source
- No need to go log in to Proxmox instances anymore to do simple tasks
- Can manage one or multiple Proxmox nodes at once


## How can I set it up?

Before starting, you need to create a Telegram bot.
Go to @Botfather at Telegram and create a new bot there.
For bot commands, put the following:

```
status - See status of your Proxmox instances
containers - List your Proxmox containers and VMs
container - Show the details of a Proxmox container/VM
node - Show the detailed info about a Proxmox node
start - Start a Proxmox container/VM
stop - Stop a Proxmox container/VM
restart - Restart a Proxmox container/VM
disks - List your Proxmox disks and their info
scale - Scale a Proxmox container/VM
```

Save the bot token somewhere, you'll need it later to for proxmox-interacter to function.

Next, you need to create the API token which will be used for interacting with your Proxmox instance.
Log in into your Proxmox, select "Datacenter" -> "Users" -> "API tokens", and create a new token.
Do not forget to uncheck the "Privilege Separation" checkbox, as otherwise it'll display data partially.
Then copy the token ID (should look like "root@pam!proxmox-interacter") and secret, and put it into your config file.

Then, you need to download the latest release from [the releases page](https://github.com/QuokkaS/proxmox-interacter/releases/). After that, you should unzip it and you are ready to go:

```sh
wget <the link from the releases page>
tar xvfz proxmox-interacter-*
./proxmox-interacter --config <path to config>
```

What you probably want to do is to have it running in the background in a detached mode. For that, first of all, we have to copy the file to the system apps folder:

```sh
sudo cp ./proxmox-interacter /usr/bin
```

Then we need to create a systemd service for our app:

```sh
sudo nano /etc/systemd/system/proxmox-interacter.service
```

You can use this template (change the user to whatever user you want this to be executed from. It's advised to create a separate user for that instead of running it from root):

```
[Unit]
Description=proxmox-interacter
After=network-online.target

[Service]
User=<username>
TimeoutStartSec=0
CPUWeight=95
IOWeight=95
ExecStart=proxmox-interacter --config <path to config>
Restart=always
RestartSec=2
LimitNOFILE=800000
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
```

Then we'll add this service to the autostart and run it:

```sh
sudo systemctl enable proxmox-interacter # set it to start on system load
sudo systemctl start proxmox-interacter  # start it
sudo systemctl status proxmox-interacter # validate it's running
```

If you need to, you can also see the logs of the process:

```sh
sudo journalctl -u proxmox-interacter -f --output cat
```

## How does it work?

It queries Proxmox instance/instances via API tokens and gets the data from a Proxmox instance/instances via
Proxmox API.

## How can I configure it?

All configuration is executed via a `.toml` config, which is passed as a `--config` variable. Check out `config.example.toml` for reference.

## How can I contribute?

Bug reports and feature requests are always welcome! If you want to contribute, feel free to open issues or PRs.
