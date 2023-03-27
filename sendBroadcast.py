import socket
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




def send_broadcast_message(message, port=5050):
    # Set the broadcast address for the current network
    broadcast_address = get_wifi_broadcast_address();

    # Create a UDP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    # Enable broadcasting on the socket
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)

    # Send the message to the broadcast address
    sock.sendto(message.encode(), (broadcast_address, port))

    # Close the socket
    sock.close()

send_broadcast_message("hello there i am David")
