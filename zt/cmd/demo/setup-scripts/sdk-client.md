# Purpose

This script sets up the SDK client side for an echo service

# Prerequisites

You need at least one controller and an edge router running. for this to work. You can use the
quick-start script found [here](https://github.com/hanzozt/zt/tree/release-next/quickstart).

# Setup

## Ensure we're logged into the controller

```action:zt-login allowRetry=true
zt edge login
```

<!--action:keep-session-alive interval=1m quiet=false-->

## Remove any entities from previous runs

```action:zt
zt edge delete identities zcat 
zt edge delete service-policies echo-dial
```

## Create and enroll the client app identity

```action:zt
zt edge create identity zcat -a echo,echo-client -o zcat.jwt
zt edge enroll --rm zcat.jwt
```

## Configure a dial policy

```action:zt
zt edge create service-policy echo-dial Dial --service-roles #echo --identity-roles #echo-client
```

## Summary

After you've configured the service side, you should now be to run the zcat client using

```
zt demo zcat -i zcat.json zt:echo
```
