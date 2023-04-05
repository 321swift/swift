# How do we identify the ip to use for the calculation of the broadcast address?

- we could try a bruteforce approach
- we could filter out the ips based on their interfaces
  - this poses a challenge because it is possible for one to rename the interface

## The way forward:

- in this case the brute force approach is much preferred  
  why:
  - it is possible for one to be connected to both wifi and ethernet at the same time
  - in that case, there ought to be a broadcast to both networks,
  - the rest of the process will be handled by the rest of the application
    since the sender would receive the sender's ip address and use it for the connection.

so in the UI, the broadcast can be configured to send the hostname as part of the message, so that the listener can see which pc to connect to.

# Proposed approach.

1. a list of all the ips of the senders pc is generated
2. the list is filtered out to remove ips of the following interfaces or forms:
   - loopback interface
   - the ipv6s
3. a broadcast is sent to each of the remaining ips
4. the process continues from there
