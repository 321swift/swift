# The broadcaster

In an attempt to establish a connection between the two pcs,  
One of the pcs involved would have to assume the role of the broadcaster.  
The broadcaster role is described as the pc who would announce its presence on the network  
in order for other pcs waiting to connect to know who to connect to.
The pc who assumes this role also eventually becomes the server: the pc to open up a web socet for the other pcs on the network to join.

## Proceedure:

The following is the proceedure for a pc to successfully assume the role of the broadcaster:

1.  - The pc starts up a web socket on a random port
    - The pc also sends a broadcast message on the network.  
      The said message contains the hostname of the pc and the port number of the port

2.  If there are any connections established on the socket, then broadcast is stopped.
3.  if the timeout limit for the broadcast is reached, then the broadcast is stopped.
4.  If the tasks 1 to 3 are successfully completed, then the pc has successfully assumed the role fo the broadcaster.
