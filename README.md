# carla-sandboxie-patchset

The experimantal patchset for [Carla](https://github.com/falkTX/Carla) VST Host on Windows with supports loading VST plugins from Sandboxie containers.

## DESCRIPTION

This patchset adds feature about load VST plugins from under the Sandbox-Plus sandbox.

Base version of Carla is [v2.3-RC1](https://github.com/falkTX/Carla/releases/tag/v2.3.0-RC1),
and Tested Sandboxie-Plus version is [v0.7.1](https://github.com/sandboxie-plus/Sandboxie/releases/tag/0.7.1)

## BUILD

### Requirements

- docker
- git
- golang (with cross-compile supprts for windows)
- make

NOTE: I tested compiles of this project under NixOS on WSL2.

### Builds

```bash
$ git clone https://github.com/nyarla/carla-sandboxie-patchset.git
$ cd carla-sandboxie-patchset
$ make
```

## USE

### How

First, copy patched carla application to Windows from `dist` directory in `carla-sandboxie-patchset`,
and launch `Carla.exe` same as non-patched Carla.

Second, setup some configuration to Sandboxie-Plus, patched Carla and Environemnt variables,
and information about this configuration is on the next _Restrictions_ section.

Last, add full path about scan plugin directory inside sandbox container to Carla's configuration,
and rescan plugin by Carla like as non-patched.

Finally, you're suceed all steps, you're enable to load VST pluguins from Sandboxie-Plus's sandbox container.

### Restrictions

this patchset features has some restrictions, see below:

#### `CARLA_SANDBOXIE_PREFIX` and `CARLA_SANDBOXIE_START` must sets

`CARLA_SANDBOXIE_PREFIX` is root directory about Sandboxie-plus sandboxes,
and `CARLA_SANDBOXIE_START` is path to `Start.exe` via Sandboxie-plus.

For example, your sandboxes root directory is `C:\Sandboxie`,
`CARLA_SANDBOXIE_PREFIX` is `C:\Sandboxie`.

And if your installed directory of Sandboxie-Plus is `C:\Sandboxie-Plus`,
`CARLA_SANDBOXIE_START` is `C:\Sandboxie-Plus\Start.exe`.

#### You must add some configuration to Sandboxie-Plus's `Sandboxie.ini`

In sandbox container, you must set this configuration to `Sandboxie.ini`
for getting works fine about load VST plugin from sandbox container:

```ini
OpenPipePath=%temp%\*
OpenIpcPath=\Sessions\*\BaseNamedObjects\carla-bridge_sem_*
```

#### Always use plugin bridge on Carla

Last, patched carla requires to always use carla's plugin bridge feature,
And you're enable to this feature by these steps:

1. Launch carla application, and click `Configure Carla` button.
2. In `Main`, checked `Enable experimental features` in `Experimantal` section
3. In `Experimental`, checked both `Enable plugin bridges` and `Run plugins in bridge mode when possible`

#### Sandboxie-Plus deletes sandbox container after the some days by default

You should block this feature by Sandboxie-Plus,
and you're enable to block this behavior by adds this configuration to `Sandboxie.ini`:

```ini
NeverDelete=y
AutoDelete=n
```

## KNOWN PROBLREMS

### `_carla-discoery-win{32,64}.exe` sometime reports crashed

This problems maybe trigger by carla default behavior,
but I'm not understood why triggered this behavior.

This problem has spamming `SBIE2224` error to notification,
but workaround is disable `SBIE2224` error notification on Sandboxie-Plus notification panel.

### `Carla.vst` cannot load as synth on FL Studio

This is a upstream bug of Carla, See Carla's issue report [#1054](https://github.com/falkTX/Carla/issues/1054).

## LICENSE

License of patch files are same as [Carla](https://github.com/falkTX/Carla/tree/main/doc),
And other files are under the MIT-licensed.

## AUTHOR

OKAMURA Naoki a.k.a. nyarla <nyarla@kalaclista.com>
