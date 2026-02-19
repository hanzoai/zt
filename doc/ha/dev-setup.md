# HA Setup for Development

**NOTE: HA is in beta. Bug reports are appreciated**

To set up a local three node HA cluster, do the following.

## Create The Necessary PKI

Run the `create-pki.sh` script found in the folder.

## Running the Controllers

1. The controller configuration files have relative paths, so make sure you're running things from
   this directory.
2. Start all three controllers
    1. `zt controller run ctrl1.yml`
    1. `zt controller run ctrl2.yml`
    1. `zt controller run ctrl3.yml`
1. Initialize the first controller using the agent
    1. `zt agent cluster init -i ctrl1 admin admin 'Default Admin'`
1. Add the other two nodes to the cluster
    1. `zt agent cluster add -i ctrl1 tls:localhost:6363`
    1. `zt agent cluster add -i ctrl1 tls:localhost:6464`

You should now have a three node cluster running. You can log into each controller individually.

1. `zt edge login localhost:1280`
2. `zt edge -i ctrl2 login localhost:1380`
3. `zt edge -i ctrl3 login localhost:1480`

You could then create some model data on any controller:

```
# This will create the client side identity and policies
zt demo setup echo client 

# This will create the server side identity and policies
zt demo setup echo single-sdk-hosted
```

Any view the results on any controller

```
zt edge login localhost:1280
zt edge ls services

zt edge login -i ctrl2 localhost:1380
zt edge -i ctrl2 ls services

zt edge login -i ctrl3 localhost:1480
zt edge -i ctrl3 ls services
```

## Running HA Go SDKs 

As of Golang SDK v1.2.3, no special changes should be required to work with HA systems.