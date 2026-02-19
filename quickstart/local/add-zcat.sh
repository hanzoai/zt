zt edge login "${ZITI_CTRL_EDGE_ADVERTISED_ADDRESS}" -u "${ZITI_USER}" -p "${ZITI_PWD}" -y

zt edge delete service zcatsvc
zt edge delete config zcatconfig

zt edge create config zcatconfig zt-tunneler-client.v1 '{ "hostname" : "zcat.zt", "port" : 7256 }'
zt edge create service zcatsvc --configs zcatconfig
zt edge create terminator zcatsvc "${ZITI_ROUTER_ADVERTISED_ADDRESS}" tcp:localhost:7256


