
FABRIC_TEST_NETWORK_SRC = ./fabric-samples/test-network
CONTRACT_SRC = ./drp-storage/chaincode-go

prerequisites:
	@echo "Installing prerequisites"
	chmod +x install.sh
	./install.sh


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

install: check-prerequisites download_script
	@echo "Installing Fabric"
	./install-fabric.sh d s b
	@echo "Installation complete"

# Ignore the couchdb setting for now, check the performance later
deploy: down
	cd $(FABRIC_TEST_NETWORK_SRC) && ./network.sh up createChannel
	./network.sh deployCC -ccn basic -ccp $(CONTRACT_SRC) -ccl go

down:
	cd $(FABRIC_TEST_NETWORK_SRC) && ./network.sh down

# Clean command to remove all materials
clean:
	rm -rf fabric-samples
	rm -f install-fabric.sh