docker run --cap-add=NET_ADMIN --device /dev/net/tun --name zt-tunneler-blue --user root --network docker_ztblue -v docker_zt-fs:/persistent --rm -it hanzozt/quickstart /bin/bash
