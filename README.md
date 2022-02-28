# vinit

[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=vinyl-linux_vin&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=vinyl-linux_vin)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=vinyl-linux_vin&metric=security_rating)](https://sonarcloud.io/dashboard?id=vinyl-linux_vin)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=vinyl-linux_vin&metric=sqale_index)](https://sonarcloud.io/dashboard?id=vinyl-linux_vin)

`vinit` is the experimental Vinyl Linux init system. It borrows heavily from many init systems, is written in go, and is distributed as static binary with no external dependencies.

It:

1. Uses `toml` files to configure both the init process, and individual services
1. Supports 'groups' of services, with the ability to override groups for customisation of boot order
1. Supports simple `SysV` style ordering by prefixing services with numbers, for example, `00-foo`, `10-bar`, etc.
1. Uses directories and symlinks to provide entrypoints, working directories, and logs in a way inspired massively by `s6`
1. Exposes operations over a gRPC powered unix socket, allowing for better programmatic control of a system

It specifically doesn't:

1. Require complex boot scripts- a service expects a file called `bin`, which usually is just a symlink
1. Require thought around boot order- I don't care about complicated boot dependency trees; use a group/ alphabetical order to order things properly
1. Do anything particularly fancy- it just supervises some services, restarting anything that needs restarting

## Installation

`vinit` can be installed/ upgraded with `vin`:

```bash
vin install vinit
```

Or simply downloaded from [https://github.com/vinyl-linux/vinit/releases](https://github.com/vinyl-linux/vinit/releases). The script `install.sh` in the latest version of `vinit` [here](https://github.com/vinyl-linux/vin-packages-stable/tree/main/vinit) can be run from within the latest release to install `vinit` to a system.

If installing manually, make sure you have some `vinit` services configured or nothing much will happen on boot (you'll get nothing but a terminal). Sample scripts may be found at [https://github.com/jspc/vinit-bootscripts](https://github.com/jspc/vinit-bootscripts).

## `vinit` Services

Services are usually installed at `/etc/vinit/services` - certainly this is the default location. Other directories can be used instead, simply by setting the environment variable `SVC_DIR` before starting `vinit`.

Each service is contained in its own directory. The directory name is used as the service name, though an optional numeric prefix may be used to affect the order in which services start.

For instance: the directory `my-application` will contain the service `my-application`. Similarly, the directory `10-my-application` also contains the service `my-application`. Because services are read alphabetically, using the numeric prefix will allow for rudimentary boot ordering.

A service directory looks like:

```bash
10-my-application
├  .config.toml
├  bin
├  environment
├  logs
│   ├ stderr
│   └ stdout
└  wd
```

Where:

1. `.config.toml` is the service configuration, including arguments to pass to `bin`, setuid configuration, and grouping information. This is described in depth below.
1. `bin` is the script/ binary/ application to run (which is usually a symlink; see: [20-dropbear/bin](https://github.com/vinyl-linux/vin-packages-stable/blob/main/dropbear/2020.81/20-dropbear/bin), which points to `/usr/sbin/dropbear`)
1. `environment` is a file containing `KEY=value` pairs, and is used to set the environment in which `bin` runs
1. `logs` is a directory containing a file for both `stdout` and `stderr` (this directory/ these files will be created if they don't exist, with each file being appended to- `vinit` doesn't handle log rotation)
1. `wd` is a directory (which is also often a symlink; see [99-vind/wd](https://github.com/vinyl-linux/vin-packages-stable/blob/main/vin/0.7.0/99-vind/wd), which points to `/etc/vinyl`

In essence, then, when the service `my-application` is started `vinit` will start `10-my-application/bin` with the args from `10-my-application/.config.toml`, from within the directory `10-my-application/wd`, and with logs going to `10-application/logs/[stderr,stdout]`.

### `.config.toml` file

A fully featured example, with optional values listed, looks like:

```toml
type = "service"           # The different types are: "service", "oneoff", "cron"
reload_signal = "SIGHUP"   # The signal to send to a process during reload- such as to reload config. Defaults to SIGHUP

[user]
user = "nobody"            # Default: root
group = "nobody"           # Default: root

[grouping]
name = "none"              # Required, but setting to an unknown group will stop it autobooting

[command]
args = "-v"                # Optional; if empty then ./bin is started with no args
ignore_output = false      # Defaults to false; governs whether stdout/stderr is ignoresd
```

Additionally, configuration for types `cron` and `oneoff` must contain (respectively):

```toml
[oneoff]
valid_exit_codes = [0]    # Exit statuses for this job that governs whether a job failed successfully


[cron]
schedule = "@daily"       # Any of the standard crontab (* * * * *) style schedule, plus the less standard (but common) things like @daily, @hourly, etc.
```


## Licence

BSD 3-Clause License

Copyright (c) 2022, Vinyl Linux
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
