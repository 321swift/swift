import socket
import netifaces
from broadcast2 import get_broadcast_address




def send_broadcast_message(message, port=5050):
    # Set the broadcast address for the current network
    broadcast_address = get_broadcast_address();

    # Create a UDP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    # Enable broadcasting on the socket
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)

    # Send the message to the broadcast address
    sock.sendto(message.encode(), (broadcast_address, port))

    # Close the socket
    sock.close()


send_broadcast_message("hello there i am David")
