zt edge controller login "${ZITI_CTRL_EDGE_ADVERTISED_ADDRESS}" -u "${ZITI_USER}" -p "${ZITI_PWD}" -y

zt edge delete service netcatsvc
zt edge delete service zcatsvc
zt edge controller delete service httpbinsvc
zt edge controller delete service iperfsvc

zt edge controller delete config netcatconfig
zt edge controller delete config zcatconfig
zt edge controller delete config httpbinsvcconfig
zt edge controller delete config iperfsvcconfig

zt edge controller create config httpbinsvcconfig zt-tunneler-client.v1 '{ "hostname" : "httpbin.zt", "port" : 8000 }'
zt edge controller create service httpbinsvc --configs httpbinsvcconfig
zt edge controller create terminator httpbinsvc "${ZITI_ROUTER_ADVERTISED_ADDRESS}" tcp:localhost:80

zt edge controller create config netcatconfig zt-tunneler-client.v1 '{ "hostname" : "localhost", "port" : 7256 }'
zt edge controller create service netcatsvc --configs netcatconfig
zt edge controller create terminator netcatsvc "${ZITI_ROUTER_ADVERTISED_ADDRESS}" tcp:localhost:7256

zt edge controller create config zcatconfig zt-tunneler-client.v1 '{ "hostname" : "zcat.zt", "port" : 7256 }'
zt edge controller create service zcatsvc --configs zcatconfig
zt edge controller create terminator zcatsvc "${ZITI_ROUTER_ADVERTISED_ADDRESS}" tcp:localhost:7256

zt edge controller create config iperfsvcconfig zt-tunneler-client.v1 '{ "hostname" : "iperf3.zt", "port" : 15000 }'
zt edge controller create service iperfsvc --configs iperfsvcconfig 
zt edge controller create terminator iperfsvc "${ZITI_ROUTER_ADVERTISED_ADDRESS}" tcp:localhost:5201

zt edge controller delete identity "test_identity"
zt edge controller create identity device "test_identity" -o "${ZITI_HOME}/test_identity".jwt

zt edge controller delete service-policy dial-all
zt edge controller create service-policy dial-all Dial --service-roles '#all' --identity-roles '#all'

#zt-enroller --jwt "${ZITI_HOME}/test_identity.jwt" -o "${ZITI_HOME}/test_identity".json

#zt-tunnel proxy netcatsvc:8145 -i "${ZITI_HOME}/test_identity".json > "${ZITI_HOME}/zt-test_identity.log" 2>&1 &
cp "${ZITI_HOME}/test_identity.jwt" /mnt/v/temp/zt-windows-tunneler/_new_id.jwt

