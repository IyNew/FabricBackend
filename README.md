# FabricBackend
A hyperledger fabric based blockchain backend storage

## To make it work
1. Clone this repo to the destination of your choice, make sure you have make and docker/docker compose installed.
> The following commands are assumed to initiated from the FabricBackend/ directory.
2. Check `instal.sh` or run `make prerequisites` to install prerequisites. 
3. Run `make install` to install the necessary components, including fabric, docker images, and chaincode.
4. Use `make drp_deploy` to bring up the test network and deploy the api server.
> This may fail at contract installation in rare occasions, simply trying the command again should resolve the problem.
5. Use `make api_server` to bring up the api service

## Initialize the ledger with csv data
Navigate to [FabricBackend/data-import-helper](./data-import-helper) and use `go run import_file.go` to put the data in ds1.csv to the ledger.
> This helper relies on the api server.

## To bring it down
1. Use `make down` to turn the network down
2. Use `make clean` to clean up the downloaded repos and docker images.

By default all the data will be cleaned when the network is down.