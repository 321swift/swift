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

## Sender specifics

1. The sender has two parts:
   - The broadcaster
   - The open socket
2. The open socket and the broadcaster must be linked together via a state variable:
   - The state variable is used to stop the broadcast when a new user connects to the sender
   - The state variable is initially false
   - The state variable be used to cause the broadcaster to stop broadcasting when the socket variable sets it to true

## Receiver specifics

1. The receiver is the part of the application that does two things
   - first, it listens for the senders broadcast message
   - second, it connects to the websocket opened by the sender.

### Proposed Approach for the receiver component:

1. The receiver starts listening for the broadcast message
2. The receiver then receives the broadcast message
3. the receiver uses obtains two things from the broadcast:
   - the sender ip address
   - the port number
4. The receiver connect to the ip address on the port number specified
5. then the file transfer session begins.
