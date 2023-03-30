# Working Description

1. two machines exist
2. one begins to listen
3. the other sends out a message
4. the listener receives the message  
   now the question comes:  
   once the listener receives the message, how does the sender know to

# idea 2

1. two machines exist
2. one will listen
3. the other will send out the broadcast.

## The breakthrough

- let the machine which will send the broadcast set up a socket before sending the broadcast.
- The sender will also include the a port number in the broadcast message
- The port number included in the broadcast will be the port on which the web socket will be set up on.
- The listener upon receiving the broadcast will then connect to the sender(server) using the received port number.

# The Proposed approach:

1. The listener starts up
2. the Sender starts up.
3. The sender sends the public encryption keys
4. The sender generates a random port number
5. The sender sets up the websocket using the generated port number
6. The sender sends the port number in an encrupted format.
7. The receiver receives the port number
8. The receiver decrypts the port number using the earlier received encryption keys
9. The receiver connects to the websocket on the port number decrypted in the previous step
