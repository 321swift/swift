import socket
from getIP import get_ip_address
from calcBroadcast import calculate_broadcast_address

def send_broadcast_message(message = get_ip_address(), dest_port = 5050):

    # Create a UDP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    # Enable broadcasting mode
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)

    # Set the destination address and port number
    dest_addr = calculate_broadcast_address('10.20.127.147/24')
    dest = (dest_addr, dest_port)

    # Send the message
    sock.sendto(bytes(message, 'utf8'), dest)

    # Close the socket
    sock.close()


send_broadcast_message("hello there i am David")
