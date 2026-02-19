docker run --cap-add=NET_ADMIN --device /dev/net/tun --name zt-tunneler-red --user root --network docker_ztred -v docker_zt-fs:/persistent --rm -it hanzozt/quickstart /bin/bash
