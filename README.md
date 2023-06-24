# Overview 👀

The Swift application is designed to transfer files from one computer to another.  
For android users, you could think of the application as "Xender for pcs" and for Iphone users, you could also think of this application as "Airdrop for pcs".

## Why use Swift❓:

This application is meant to be a one size fits all when it comes to sending and receiving files.

### **The Problem:**

It is often the case where two people in the same room or on the same wifi network need to transfer files and yet cannot do so because there is no real direct way of doing so without having to setup some network groups or configure some network stuff.  
While that may be easy to do for some, it may not be the case for all people, especially those who are not tech savy, these people however form the majority of the target demographic.

---

The swift application allows said people to transfer files from one computer to another no matter the operating system they are on. All that is required is that the two are on the same wifi network (**technically speaking the same subnet**).

# Architecture

-   The swift application as of now has only one major component in it called a node.
-   A node can do one of two things at a time; send a file or receive a file
-   A node has two components; the backend and the user interface (UI).
-   The user interface of a node opens up in a browser while the backend starts up as a console application.
-   The backend and the User Interface communicate via a web socket.
-   Both the backend and the frontend need each other to work properly.

# How the application works:

-   When the application starts, the backend first starts up and then the backend starts up the user interface.
-   The user interface when opens up in the default browser of the user.
-   The user will at this point have to select a role {sending / receiving}.
    -   note that Selecting the role determines how the backend will connect to the other computer.
    -   If two computers need to share files, selecting send on one means the "receive" must be selected on the other
    -   **_No matter the selected role, the user will still be able to send or receive files_**
-   After selecting the role, the computer will connect to the other node who is waiting to receive or send a file.
-   The user will then be allowed to select the file or files to be sent.

---

-   ⚠⚠ ‼‼ At this point in development, the application does not display how much of the file has been sent or received. so the user will unfortunately have to wait for the app to display about 10 messages containing the text "Writing to file" in the console output of the backend.

-   While we work to fix this in subsequent releases, we do hope users enjoy the application
    -   We plead with users to also give good feedback by submitting issues to the issues section of the github repository.
-   ‼ Issues submitted should contain relevant information including the following:
    -   A screenshot of the the backend console output and the UI at the time the issue occurs.
    -   A brief description of the issue and any other behaviours that occur while using the application that you feel should not be the case.

## **⚠‼ Expected Anomalies.**

-   For windows users, the OS might prompt you to allow or restrict access to public and private wifi networks. When this happens please allow both networks.
-   Please allow the application permissions on your firewall
    -   note that for antivirus software which has a firewall you would need to allow the application to bypass the firewall else the other computer will not be able to connect.
-   In some cases, windows does not prompt the user to allow the app to use public or private networks; in this case the user will need to manually add inbound and outbound rules for the swift application.

# 🎯 Installation:

Download the release from the github repository,  
unzip the file and run the executable... As simple as that 😉
