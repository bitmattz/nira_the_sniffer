# portHandler.go

Currently the local scan is fast, but when it comes to scan another devices, it take a lot of time to scan all the ports. Here is what it needs to be done:

    - Concurrency: Work with routines to scan the ports in parallel
    - One connection per ip: Now every port is opening a connection, reducing the time and increasing uneccessary calls.
    -User SYN scan 
    
