
REPO_SRC = $(shell pwd)
FABRIC_TEST_NETWORK_SRC = $(REPO_SRC)/fabric-samples/test-network
CONTRACT_SRC = $(REPO_SRC)/drp-storage/chaincode-go
CLIENT_SRC = $(REPO_SRC)/drp-client
TEST_CONTRACT_SRC = $(REPO_SRC)/fabric-samples/asset-transfer-basic/chaincode-go

prerequisites: check-prerequisites
	@echo "Installing prerequisites"
	chmod +x pre-requisites.sh
	./pre-requisites.sh


check-prerequisites:
	@echo "Make sure you have the following installed:"
	@echo "1. git"
	@echo "2. curl"
	@echo "3. docker/docker-compose"
	@echo "4. go (for chaincode)"
	@echo "5. jq"

# Download a script to the current path
download_script:
	curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh

install: download_script
	@echo "Installing Fabric"
	./install-fabric.sh d s b
	@echo "Installation complete"

network_up:
	cd $(FABRIC_TEST_NETWORK_SRC) && ./network.sh up createChannel

network_up_couchdb:
	cd $(FABRIC_TEST_NETWORK_SRC) && ./network.sh up createChannel -ca -s couchdb

# Ignore the couchdb setting for now, check the performance later
drp_deploy: down network_up
	cd $(FABRIC_TEST_NETWORK_SRC) && ./network.sh deployCC -ccn basic -ccp $(CONTRACT_SRC) -ccl go

test_basic_deploy: down network_up
	cd $(FABRIC_TEST_NETWORK_SRC) && ./network.sh deployCC -ccn basic -ccp $(TEST_CONTRACT_SRC) -ccl go

drp_couchdb_deploy: down network_up_couchdb
	cd $(FABRIC_TEST_NETWORK_SRC) && ./network.sh deployCC -ccn basic -ccp $(CONTRACT_SRC) -ccl go 

api_server:
	cd $(CLIENT_SRC) && go run main.go

down:
	cd $(FABRIC_TEST_NETWORK_SRC) && ./network.sh down

# Clean command to remove all materials
clean:
	rm -rf fabric-samples
	rm -f install-fabric.sh