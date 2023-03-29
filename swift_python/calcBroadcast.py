def calculate_broadcast_address(ip_subnet):
    """
    Calculate the broadcast address of a network given an IP address and subnet mask in the format of "192.168.190.1/24"
    """
    # Split the IP address and subnet mask
    ip_address, subnet_mask_bits = ip_subnet.split('/')
    
    # Convert IP address to integer value
    ip_int = int(''.join([bin(int(x)+256)[3:] for x in ip_address.split('.')]), 2)
    
    # Calculate the subnet mask integer value
    subnet_mask_int = ((1 << int(subnet_mask_bits)) - 1) << (32 - int(subnet_mask_bits))
    
    # Calculate the network address and the host address bitmask
    network_address_int = ip_int & subnet_mask_int
    host_address_bitmask_int = (~subnet_mask_int) & 0xFFFFFFFF
    
    # Calculate the broadcast address
    broadcast_address_int = network_address_int | host_address_bitmask_int
    
    # Convert the broadcast address to a string representation
    broadcast_address_str = '.'.join([str(int(broadcast_address_int >> (i * 8) & 0xFF)) for i in range(4)][::-1])
    
    return broadcast_address_str

def get_broadcast_address():
    return calculate_broadcast_address('10.20.127.147/24')
# print(calculate_broadcast_address('10.20.127.147/24')))

