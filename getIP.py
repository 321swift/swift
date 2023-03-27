import socket
import os

def send_broadcast_message(message, port=5000):
    # Get the IP address of the local machine
    ip_address = os.popen('hostname -I').read().strip()

    # Set the broadcast address for the current network
    broadcast_address = "<your_broadcast_address>"

    # Create a UDP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    # Enable broadcasting on the socket
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)

    # Send the message to the broadcast address
    sock.sendto(message.encode(), (broadcast_address, port))

    # Close the socket
    sock.close()

    return ip_address
