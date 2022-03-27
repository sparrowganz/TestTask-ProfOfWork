# TestTask-ProfOfWork

---
## Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.  
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.  
• The choice of the POW algorithm should be explained.  
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
• Docker file should be provided both for the server and for the client that solves the POW challenge

---
## Starting step by step
Building app: ```make build``` - for build client and server app

Running app (server as daemon): ```make run``` - create network , start server as daemon and client interactively

Restart app: ```make clean && make run``` - remove all dependencies and start new 

Stop and clean: ```make clean``` - remove all dependencies

---
## Application:
This app controlled user activity in tcp-connection. If user starting send many requests - system blocked him.
Limit - 5 requests per 5 second (Demonstration version).
System send question if requests limit exceeding.
And User can send answer up to 3 times.

---
## Client-Server Contract
```GET``` - get random excerption

```STOP``` - exit onto tcp session 
