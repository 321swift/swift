import socket
import ipaddress
import netifaces


def get_wifi_broadcast_address():
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

                # Calculate the broadcast address
                broadcast = socket.inet_ntoa(
                    struct.pack('!I', 
                                (struct.unpack('!I', 
                                               socket.inet_aton(ipv4['addr']))[0] 
                                 | ((1 << 32 - netmask) - 1))
                                )
                    )

                return broadcast
    return None



def receive_broadcast():
    # create a UDP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)


    # bind the socket to a specific address and port
    sock.bind((get_wifi_broadcast_address(), 5050))

    # receive messages in a loop
    while True:
        data, addr = sock.recvfrom(1024)  # buffer size is 1024 bytes
        print(f"Received message from {addr}: {data.decode()}")

receive_broadcast();
