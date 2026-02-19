
# Run Ziti Controller in Docker

You can use this container image to run a Ziti Controller in a Docker container.

## Container Image

The `hanzozt/zt-controller` image is thin and is based on the `hanzozt/zt-cli` image, which only provides the
`zt` CLI. The `zt-controller` image adds an entrypoint that provides controller bootstrapping when
`ZITI_BOOTSTRAP=true` and uses the same defaults and options as the Linux package.

## Docker Compose

The included `compose.yml` demonstrates how to bootstrap a controller container.

### Example

At a minimum, you must set the address and password options in the parent env or set every recurrence in the compose file.

```text
# fetch the compose file for the zt-router image
wget https://get.hanzozt.dev/dist/docker-images/zt-controller/compose.yml

ZITI_PWD="mypass" \
ZITI_CTRL_ADVERTISED_ADDRESS=ctrl.127.21.71.0.sslip.io \
    docker compose up
```

After a few seconds, `docker compose ps` will show a "healthy" status for the controller.

Then, you may log in to the controller using the `zt` CLI.

```text
zt edge login ctrl.127.21.71.0.sslip.io:1280 -u admin -p mypass
```
