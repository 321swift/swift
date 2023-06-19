-------Overview--------
The Swift application is designed to transfer files from one computer to another using a secure and efficient method. The application consists of two parts - a sender and a receiver - and uses a state machine architecture to manage the transfer process.

------Architecture-------
The Swift application uses a state machine architecture to manage the file transfer process. The state machine consists of a struct or object containing variables that define the current state of the application at any given time. It also includes data channels for sending data out of the core of the application and receiving input from outside the application.

--------Installation-------
Download the Swift application installer from the official website
Run the installer on both the sender and receiver computers
Follow the prompts to complete the installation process
Usage
To use the Swift application, follow these steps:

-Start up the sender and receiver computers
-On the sender computer, open the Swift application and generate a random port number
-Set up the websocket on the generated port number
-Send the public encryption keys and the encrypted port number to the receiver computer
-On the receiver computer, decrypt the port number using the encryption keys
-Connect to the sender computer using the decrypted port number
-Start the file transfer process
