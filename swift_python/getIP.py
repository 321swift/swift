import socket
    
def get_ip_address():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect(("8.8.8.8", 80))
    return s.getsockname()[0]

# ip = ni.ifaddresses('wlp2s0')[ni.AF_INET][0]['addr']
print(get_ip_address())  # should print "192.168.100.37"

