import socket
import struct
import netifaces
import ipaddress

def get_network_interface_info():
    # Get the list of network interfaces on the local machine
    interfaces = netifaces.interfaces()

    # Iterate over the interfaces and find the WiFi interface
    for iface in interfaces:
        if iface.startswith('w'):
            addresses = netifaces.ifaddresses(iface)
            if netifaces.AF_INET in addresses:
                # Get the IPv4 address and netmask for the interface
                ipv4 = addresses[netifaces.AF_INET][0]
                netmask = ipv4['netmask']

                return ipv4['addr'], netmask

    return None, None

def get_broadcast_address():
    # Get the IP address and netmask for the network interface
    ip_address, netmask = get_network_interface_info()

    if ip_address and netmask:
        # Convert the IP address and netmask to a network address
        network_address = ipaddress.IPv4Network(f'{ip_address}/{netmask}', strict=False)

        # Get the broadcast address from the network address
        broadcast_address = str(network_address.broadcast_address)

        return broadcast_address

    return None

