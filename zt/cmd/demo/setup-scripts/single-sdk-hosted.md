# Purpose

This script sets up an echo service which is hosted by an SDK application.

# Prerequisites

You need at least one controller and an edge router running. for this to work. You can use the
quick-start script found [here](https://github.com/hanzozt/zt/tree/release-next/quickstart).

# Setup

## Ensure we're logged into the controller.

```action:zt-login allowRetry=true
zt edge login
```

```action:keep-session-alive interval=1m
If your session times out you can run zt edge login again.
```

## Remove any entities from previous runs.

```action:zt
zt edge delete service echo
zt edge delete config echo-host
zt edge delete identities echo-host-1 echo-host-2
zt edge delete service-policies echo-bind
zt edge delete edge-router-policies echo
zt edge delete service-edge-router-policies echo 
```

## Create the echo service

```action:zt
zt edge create service echo -a echo
```

## Create and enroll the hosting identity

```action:zt
zt edge create identity echo-host-1 -a echo,echo-host -o echo-host-1.jwt
zt edge enroll --rm echo-host-1.jwt
```

# Configure policies

```action:zt
zt edge create service-policy echo-bind Bind --service-roles @echo --identity-roles #echo-host
zt edge create edge-router-policy echo --identity-roles #echo --edge-router-roles #all
zt edge create service-edge-router-policy echo --service-roles @echo --edge-router-roles #all
```

You should now be to run the echo server with

```
zt demo echo-server -i echo-host-1.json
```

and the zcat client using

```
zt demo zcat -i zcat.json zt:echo
```
