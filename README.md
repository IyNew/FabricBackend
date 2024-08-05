# FabricBackend
A hyperledger fabric based blockchain backend storage

## To make it work
1. Clone this repo to the destination of your choice, make sure you have make and docker/docker compose installed.
> The following commands are assumed to initiated from the FabricBackend/ directory.
2. Check `instal.sh` or run `make prerequisites` to install prerequisites. 
3. Run `make install` to install the necessary components, including fabric, docker images, and chaincode.
4. Use `make deploy` to bring up the test network and deploy the api server.
5. Use `make down` to turn the network down
6. Use `make clean` to clean up the downloaded repos and docker images.

By default all the data will be cleaned when the network is down.