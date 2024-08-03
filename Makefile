prerequisites:
	@echo "Installing prerequisites"
	chmod +x install.sh
	install.sh


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

install: check-prerequisites prerequisites download_script
	@echo "Installing Fabric"
	./install-fabric.sh d s b
	@echo "Installation complete"

# Clean command to remove all materials
clean:
	rm -rf fabric-samples
	rm -f install-fabric.sh