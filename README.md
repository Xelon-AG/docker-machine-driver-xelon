# Docker Machine Xelon Driver

[![Build Status](https://circleci.com/gh/Xelon-AG/docker-machine-driver-xelon.svg?style=shield)](https://circleci.com/gh/Xelon-AG/docker-machine-driver-xelon)

Create Docker machines on [Xelon](https://www.xelon.ch/).

You need to use your token and pass that to `docker-machine create` with `--xelon-token` option.


## Usage

    $ docker-machine create --driver xelon \
        --xelon-token <YOUR-TOKEN> \
        MY_INSTANCE

If you encounter any troubles, activate the debug mode with `docker-machine --debug create ...`.

### When explicitly passing environment variables

    $ export XELON_TOKEN=<YOUR-TOKEN>
    $ docker-machine create --driver xelon MY_INSTANCE


## Options

- `--xelon-api-base-url`: Xelon API base URL.
- `--xelon-cpu-cores`: Number of CPU cores for the device.
- `--xelon-device-password`: Password for the device.
- `--xelon-disk-size`: Drive size for the device in GB.
- `--xelon-kubernetes-id`: Kubernetes ID for the device.
- `--xelon-memory`: Size of memory for the device in GB.
- `--xelon-ssh-port`: SSH port to connect.
- `--xelon-ssh-user`: SSH username to connect.
- `--xelon-swap-disk-size`: Swap disk size for the device in GB.
- `--xelon-token`: **required** Xelon authentication token.

#### Environment variables and default values

 CLI option                 | Environment variable    | Default                           |
| ------------------------- | ----------------------- | --------------------------------- |
| `--xelon-api-base-url`    | `XELON_API_BASE_URL`    | `https://vdc.xelon.ch/api/user/`  |
| `--xelon-cpu-cores`       | `XELON_CPU_CORES`       | `2`                               |
| `--xelon-device-password` | `XELON_DEVICE_PASSWORD` | `Xelon22`                         |
| `--xelon-disk-size`       | `XELON_DISK_SIZE`       | `20`                              |
| `--xelon-kubernetes-id`   | `XELON_KUBERNETES_ID`   | `kub1`                            |
| `--xelon-memory`          | `XELON_MEMORY`          | `2`                               |
| `--xelon-ssh-port`        | `XELON_SSH_PORT`        | `22`                              |
| `--xelon-ssh-user`        | `XELON_SSH_USER`        | `root`                            |
| `--xelon-swap-disk-size`  | `XELON_SWAP_DISK_SIZE`  | `2`                               |
| **`--xelon-token`**       | `XELON_TOKEN`           | -                                 |


## Release process

The release process uses continuous integration from CircelCI which means the code base should
be ready to release any time.

#### Checklist for releasing Xelon Driver

The below steps are for final release:

1. Check that latest build compiles and passes tests.
2. Trigger release creation by tagging git commit with
   ```bash
   # X.Y.Z. is a release version, e.g. 1.2.0
   $ git tag -a vX.Y.Z -m "Release vX.Y.Z."
   ```
3. Wait for the build to finish and release created.
4. Update release information if needed


## Contributing

We hope you'll get involved! Read our [Contributors' Guide](.github/CONTRIBUTING.md) for details.
