# install the prerequisites
sudo apt-get install git curl -y

# install go and jq
sudo apt-get install golang-go jq -y

# Make sure the Docker daemon is running.
# sudo systemctl start docker

# Add your user to the Docker group.
sudo usermod -a -G docker $USER

# Check version numbers  
docker --version
docker-compose --version
go version
jq --version
