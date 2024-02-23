# A simple to-do list application with echo, redis and docker 

## This project constists of client side and server side. Server accepts "orders" to complete. It verifies everything and then passes them to client side. Client has an ability to close given orders i.e complete them by recieving http requests.

## Starting with the project 
Clone my repository to your machine with ```git clone github.com/mishaRomanov/redis-project```
## Running the whole thing
Go to cloned directory and do ``docker compose up -d``
## How to use 
Send a `POST` request to `localhost:8080/order` with json body like this `{"description":"test")`.

This creates a new instance in `redis` database which runs in container and then passes it to client side.
Then it returns you the ID of new order so you can interact with `client` and close it.

To close the order you need to send get request to `localhost:3030/close/*id*` where id is the number of order you want to close. You received that number the step before.
#### You can only close the order using client side endpoint `localhost:3030/close/:id`
## How it works
You run the following command `docker compose up -d` which starts 3 containers: server, client and redis.
Server accepts new orders and then passes them to client. Client receives and displays all orders (you can see them in container logs). 
Only client can access the order closure since the whole service runs with JWT token authorization. 