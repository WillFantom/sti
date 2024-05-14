# **`sti`**: Network Performance Tests    ![GitHub release (latest SemVer)](https://img.shields.io/github/v/tag/willfantom/sti?display_name=tag&label=%20&sort=semver)

A simple test suite for probing network performance on the running host. This
runs a set of configurable network tests including:
  - Speedtest (via [speedtest.net](https://speedtest.net))
  - Iperf3
  - Ping

The results from these tests are pushed to a given InfluxDB v2 server.

Since these tests can be taxing on the network, each is run sequentially with a
configurable interval between each.

## Config

Your config file should be a _yaml_ file named `config.yaml` (or specified in the
`--config` flag if named otherwise). Config files are looked for in:
  - The current working dir when the `sti` command is executed
  - `/etc/sti/config.yaml`
  - The `.config/sti` dir within the executing user's home

See the [example config](./example-config.yaml) for more on what should be in a
configuration.


## Usage

Assuming you have Docker installed and a config file for `sti` located in your
current working directory, the following command will run the program.

```
docker run --rm --name sti --network host \
  -v "$(pwd)":/etc/sti/ -v /tmp:/tmp \
  ghcr.io/willfantom/sti:latest
```
