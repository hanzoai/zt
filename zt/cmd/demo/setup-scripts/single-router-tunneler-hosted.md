# Purpose

This script sets up hosting for an echo service which is hosted by a router-embedded tunneler.

# Prerequisites

You need at least one controller and an edge router running. for this to work.
You can use the quick-start script found [here](https://github.com/hanzozt/zt/tree/release-next/quickstart).

# Setup

## Ensure we're logged into the controller.

```action:zt-login allowRetry=true
zt edge login
```

<!--action:keep-session-alive interval=1m quiet=false-->

## Remove any entities from previous runs.

```action:zt
zt edge delete service echo
zt edge delete config echo-host
zt edge delete identities echo-host-1 echo-host-2
zt edge delete service-policies echo-bind
zt edge delete edge-router-policies echo
zt edge delete service-edge-router-policies echo 
```

## Create the echo-host config

```action:zt-create-config name=echo-host type=host.v2
{
    "terminators" : [
        {
            "address" : "localhost",
            "port" : 1234,
            "protocol" : "tcp"   
        }
    ]
}
```

## Create the echo service

```action:zt
zt edge create service echo -c echo-host -a echo
```

## Update edge-routers

Make sure demo edge routers are tunneler enabled and the associated identity has the `echo-host` attribute.
Only routers with the `demo` role attribute will be updated.

```action:zt-for-each type=edge-routers minCount=1 maxCount=2 filter='anyOf(roleAttributes)="demo"'
zt edge update edge-router ${entityName} --tunneler-enabled
zt edge update identity ${entityName} --role-attributes echo-host 
```

## Configure policies

```action:zt
zt edge create service-policy echo-bind Bind --service-roles @echo --identity-roles #echo-host
zt edge create edge-router-policy echo --identity-roles #echo --edge-router-roles #all
zt edge create service-edge-router-policy echo --service-roles @echo --edge-router-roles #all
```

## Summary

You should now be to run the echo server with

```
zt demo echo-server -p 1234
```

and the zcat client using

```
zt demo zcat -i zcat.json zt:echo
```
