import socket
from calcBroadcast import get_broadcast_address


# def receive_broadcast():
#     # create a UDP socket
#     sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)


#     # bind the socket to a specific address and port
#     sock.bind((get_broadcast_address(), 5050))

#     # receive messages in a loop
#     while True:
#         data, addr = sock.recvfrom(1024)  # buffer size is 1024 bytes
#         print(f"Received message from {addr}: {data.decode()}")


def listen_on_port(broadcast_address=get_broadcast_address(), port=5050):
    # Create a UDP socket and bind it to the broadcast address and port
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    print(broadcast_address)
    sock.bind((broadcast_address, port))

    print(f'Listening for messages on {broadcast_address}:{port}...')

    # Continuously receive and print messages
    while True:
        data, address = sock.recvfrom(1024)
        print(f'Received message from {address}: {data.decode()}')

listen_on_port()
