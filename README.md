# FlyStatic

Documents in other languages: [中文文档](docs/README.zh-CN.md)

FlyStatic is yet another lightweight and open source static file server.

## Features

- Basic authentication
- Multiple users
  - allows each user to have a different root directory and path prefix which are extremely well isolated
- TLS/SSL support
- HTTP Range support (for video streaming and large file downloading)

The goal of FlyStatic and FlyDav are to keep things easy and simple, with users able to deploy the service in no time. Therefore any extraneous features are not added.

## Get started in 30 seconds

1. Start by downloading FlyStatic from their website at [release page](https://github.com/pluveto/flystatic/releases).
2. Run `./flystatic -H 0.0.0.0` to start the server. Warning: this will start the server with default settings, which is not secure. Please read the [Configuring FlyStatic](#configuring-flystatic) section for more information.
3. Open `http://YOUR_IP:8080/static file` in your web browser to view files.

## Command line options

```bash
$ FlyStatic -h
--------------------------------------------------------------------------------
Usage: FlyStatic [--host HOST] [--port PORT] [--user USER] [--verbose] [--config CONFIG]

Options:
  --host HOST, -H HOST   host address
  --port PORT, -p PORT   port
  --user USER, -u USER   username
  --verbose, -v          verbose output
  --config CONFIG, -c CONFIG
                         config file
  --help, -h             display this help and exit
```

If you have a config file, you can ignore the command line options. Run `FlyStatic -c /path/to/config.toml` to start the server.

If you want to quickly start the server with host, port, username and a one-time password, you can run `FlyStatic -H IP -p PORT -u USERNAME` to start the server. Then you'll input the password for the user. And the server will serve at `http://IP:PORT/`.

## Configuring FlyStatic

1. Start by downloading FlyStatic from their website at [release page](https://github.com/pluveto/flystatic/releases).
2. Now that you have the software, you need to create a configuration file for it. Start by creating a new file called `FlyStatic.toml`.
3. Inside the configuration file, you will need to add the following information:
    - `[server]`: This section will define the host, port, and path of the static file server.
    - `host`: The IP address of the host. This should be set to “0.0.0.0” if you want to make the server accessible from any IP address.
    - `port`: The port number to use for the static file server.
    - `path`: The path of the static file server.
    - `fs_dir`: The directory on the server where the static file files will be stored.
    - `[auth]`: This section will define the authentication settings for the static file server.
    - `[[auth.user]]`: This subsection will define the username and credentials for each user that has access to the static file server.
        - `username`: The username of the user.
        - `sub_fs_dir`: The subdirectory of the fs_dir to which the user will have access.
        - `sub_path`: The path that the user will access the static file server from.
        - `password_hash`: The hashed password of the user.
        - `password_crypt`: The type of hashing algorithm used to hash the password. This should be set to “bcrypt” or “sha256”.
    - `[log]`: This section will define the logging settings for the static file server.
    - `level`: The log level of the server. This can be set to “debug”, “info”, “warn”, “error”, or “fatal”.
    - `[[log.file]]`: This subsection will define the settings for the log file. Ignore this subsection if you do not want to log to a file.
        - `format`: The format of the log file. This can be set to “json” or “text”.
        - `path`: The path of the log file.
        - `max_size`: The maximum size of the log file in megabytes.
        - `max_age`: The maximum age of the log file in days.
    - `[[log.stdout]]`: This subsection will define the settings for the log output to the console. Ignore this subsection if you do not want to log to the console.
        - `format`: The format of the log output. This can be set to “json” or “text”.
        - `output`: The output stream for the log output. This can be set to “stdout” or “stderr”.
4. Save the configuration file and run the FlyStatic server. You should now be able to access the static file server with the configured settings.

To get a example configuration file, go to [conf dir](https://github.com/pluveto/flystatic/blob/main/conf).

## Install as a service

### Install as a service on Linux

1. Create a new file called `FlyStatic.service` in `/etc/systemd/system/` and add the following information:

File `/etc/systemd/system/FlyStatic.service`

```ini
[Unit]
Description = FlyStatic Server
After = network.target syslog.target
Wants = network.target

[Service]
Type = simple
# !!! Change the binary location and config path to your own !!!
ExecStart = /usr/local/bin/FlyStatic -c /etc/FlyStatic/FlyStatic.toml

[Install]
WantedBy = multi-user.target
```

2 Run `systemctl daemon-reload` to reload the systemd daemon.
3 Run `systemctl enable FlyStatic` to enable the service.
4 Run `systemctl start FlyStatic` to start the service.

### Manage the service

- Run `systemctl status FlyStatic` to check the status of the service.
- Run `systemctl stop FlyStatic` to stop the service.

## Features

- [x] Basic authentication
- [x] Multiple users
- [x] Different root directory for each user
- [x] Different path prefix for each user
- [x] Logging
- [ ] SSL
  - Work in progress
  - You can use a reverse proxy like Nginx to enable SSL.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
