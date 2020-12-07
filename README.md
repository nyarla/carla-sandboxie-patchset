# carla-sandboxie-patchset

- The experimantal patchset for Carla VST Host on Windows with supports Sandboxied VST Plugins

## DESCRIPTION

This patchset is enable to build for experimantal windows on Ubuntu 20.04 about Carla,
and add feature of sandboxie-plus supports about load VST plugins can loading from sandboxed containers.

## BUILD

### Requirements

- docker
- bash
- golang (with supports for build windows binary)
- make

NOTE: Build of this patchset is tested on linux environment with mingw cross compile.

### How to build

`git clone` from this repository, enter to repository dir, and just run:

```bash
$ make
```

built artifacts can getting from `dist/` directory on repository dir.

## USE

### How to use

Copy carla application from `dist/` dir, and put on Windows machine,
so you're able to use it as nomarlly carla.

However, sandboxie supports on patched carla has some restrictions,
about this restrictions see below:

- setenv `CARLA_SANDBOXIE_PREFIX` and `CARLA_SANDBOXIE_START` is required.
  - `CARLA_SANDBOXIE_PREFIX` is same as root directory of Sandboxie's sandboxes
    - this value is able to get from Sandboxie-plus's `Sandbox file system root`
    - and notes, this value must not include `%SANDBOX%` placeholder
  - `CARLA_SANDBOXIE_START` is path to `Start.exe` about Sandboxie's command.
- set configuration about always use plugin bridge on carla
  - this features can enable on experimantal features in carla configuration
  - checked `Enable experimantal features` in `Main` section on carla config
  - and checked`Enable plugin bridges` in `Exprimental` section too.

And, you must include this configuration in Sandboxie's container config:

```ini
OpenPipePath=%temp%\*
OpenIpcPath=\Sessions\*\BaseNamedObjects\carla-bridge_sem_*
```

in additional, you want VST container is persistents, you should add this configuration:

```ini
NeverDelete=y
AutoDelete=n
```

### How to add VST plugins from Sandboxie's containers

You're able to add VST directory as normally VST plugins dir.

For example, Your sandbox directory is `P:\sandbox\VST`,
and your VST directory is `P:\sandbox\VST\drive\C\VST2`,
you just add `P:\sandbox\VST\drive\C\VST2` directory to VST2 directory configuration,
and rescan plugins on carla.

So you added this configuration and rescan plugins on carla,
patched carla is auto-detected Sandboxie's sandbox containers,
and patched carla uses Sandboxie's `Start.exe` command when launch plugin discovery and load plugins.

More information, when you added non-sandboxed VST directory,
patched carla does not use `Start.exe` command.

## KNOWN PROBLREMS

- real `_carla-discovery-win{32,64}.exe` is crash some(?) plugins
  - this problem is spamming `SBIE2224` error to windows nortification
  - workaround is disable nortification by global settings on sandboxie-plus
- when try to display VST GUI from carla rack UI, VST GUI not shown at first time
  - workaround is retry to display VST GUI twice

## LICENSE

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

## AUTHOR

OKAMURA Naoki a.k.a. nyarla <nyarla@kalaclista.com>

## TESTED ENVIRONEMT

- cross-compile by mingw is test on docker in NixOS on WSL2
- based carla version is [v2.2.0](https://github.com/falkTX/Carla/releases/tag/v2.2.0)
- Sandboxie version is [Release v0.4.5 / 5.44.1](https://github.com/sandboxie-plus/Sandboxie/releases/tag/v0.4.5)

# NOTES

- cx_Freeze patch is based on [this pull request](https://github.com/marcelotduarte/cx_Freeze/pull/545)
