# carla-sandboxie-patchset

- The experimantal patchset for Carla VST Host on Windows with Sandboxied VST Plugins

## So, What?

This patchset is enable to experimental build on Ubuntu 20.04 about Carla,
and add features to load VST plugins from Sandboxie's sandboxed containers.

## Building Carla VST Host with this patchset

### Requirements

- bash
- docker
- golang with build binary supports for windows
- make
- linux environment

NOTE: I tested build about this repository is inside NixOS on WSL2.

### How to Build

Clone to this repository, and just run:

```bash
$ make
```

And custom built binary is getting from `dist` directory

## How to use

Just add plugin directory and rescan on Carla VST Host,
it supports both Sandboxied directory and Non-Sandboxied directory.

### Restrictions of use to this custom built

#### You must set `CARLA_SANDBOXIE_PREFIX` and `CARLA_SANDBOXIE_START`

These environment variable is required by working for this custom built:

- `CARLA_SANDBOXIE_PREFIX` is your _Sandbox file system root_ on Sandboxie without `%SANDBOX%` param.
- `CARLA_SANDBOXIE_START` is path to Sandboxie's `Start.exe`

#### You must always enable to plugin bridge

This patchset modify to proces of plugin bridge on Carla,
and other cases about load VST plugins is just not working, or broken (?)

How to configure:

- In `Main` section, enable to _Enable experimantal features_
- In `Exprimental` section, enable to and _Enable plugin bridges_ and _Run plugins in bridge when possible_

#### Some configuration required on Sandboxie containers

- In `General Option`, `Sandbox Indicator in title` is _Don't alter the window title_
- In `Resource Access`, these configurattions required:
  - `OpenFilePath=%temp%\*`
  - `OpenIpcPath=\Sessions\*\BaseNamedObjects\carla-bridge_sem_*`

## Known Problems

When `_carla-discovery-win{32,64}.exe` cannot load VST plugin by differenct architecture inside Sandboxie's container,

Sandboxie throws notification of `SBIE2224` error too many.

## Licenses

These patches are under the [GPL v2, or later](https://www.gnu.org/licenses/old-licenses/gpl-2.0.txt):

- fix-buildscript.patch
- fix-glib-patch.patch
- fix-pack-win.patch
- fix-pyqt-deps.patch
- sandboxie-discovery.patch
- sandboxie.patch

This patch is under the [Python Software Foundation License](https://docs.python.org/3/license.html):

- cx_Freeze.patch

Other script or files are under the MIT-licensed.

# Author

OKAMURA Naoki a.k.a. nyarla <nyarla@kalaclista.com>

# Tested Environment

- build on NixOS on WSL2
- Based version of Carla VST Host is [v2.2.0](https://github.com/falkTX/Carla/releases/tag/v2.2.0)
- Sandboxie version is [Release v0.4.3 / 5.43.7](https://github.com/sandboxie-plus/Sandboxie/releases/tag/v0.4.3)

# NOTES

- GPLv2-ed patchs includes original Carla source code
- cx_Freeze patch is based on [this pull request](https://github.com/marcelotduarte/cx_Freeze/pull/545)
