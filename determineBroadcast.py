import socket
import ipaddress

def get_broadcast_address():
    # get the IP address of the current host
    hostname = socket.gethostname()
    ip_address = socket.gethostbyname(hostname)

    # get the network address and netmask of the current network
    net_address = ipaddress.IPv4Network(f"{ip_address}/24")  # assuming a subnet mask of /24

    # determine the broadcast address
    broadcast_address = net_address.broadcast_address

    return str(broadcast_address)

