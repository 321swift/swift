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


def receive_broadcast():
    # create a UDP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)


    # bind the socket to a specific address and port
    sock.bind((get_broadcast_address(), 5050))

    # receive messages in a loop
    while True:
        data, addr = sock.recvfrom(1024)  # buffer size is 1024 bytes
        print(f"Received message from {addr}: {data.decode()}")
