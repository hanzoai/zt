suffix=$(date +"%b-%d-%H%M")
idname="User${suffix}"

zt edge login "${ZITI_CTRL_EDGE_ADVERTISED_ADDRESS}" -u "${ZITI_USER}" -p "${ZITI_PWD}" -y

zt edge delete identity "${idname}"
zt edge create identity device "${idname}" -o "${ZITI_HOME}/test_identity".jwt

cp "${ZITI_HOME}/test_identity".jwt /mnt/v/temp/zt-windows-tunneler/_new_id.jwt

