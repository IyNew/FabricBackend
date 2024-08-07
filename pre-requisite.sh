# install the prerequisites
# sudo apt-get install git curl -y

# install go and jq
# sudo apt-get install golang-go jq -y

# Make sure the Docker daemon is running.
# sudo systemctl start docker

# Add your user to the Docker group.
# sudo usermod -a -G docker $USER

# Check version numbers  
if [ -x "$(command -v docker)" ]; then
    docker --version
else
    echo "error: docker is not installed"
fi

# check if docker needs sudo
if docker run hello-world; then
    echo "Docker does not require sudo"
else
    echo "Docker requires sudo"
    echo "Add your user to the Docker group to run Docker without sudo."
fi

if [ -x "$(command -v docker-compose)" ]; then
    docker-compose --version
else
    echo "error: docker-compose is not installed"
fi

if [ -x "$(jq --version)" ]; then
    jq --version
else
    echo "error: jq is not installed"
    echco ""
fi

# get the go version number, and check if it is 1.22.5 or higher
if go version | grep -q "go1.22.[5-9]\|go1.2[3-9][0-9]\|go[2-9][0-9]\|go[1-9][0-9][0-9]"; then
    echo "Go version satisfies the requirement"
    echo "use 'make install' and 'make deploy' to start the network"
else
    echo "error: Go version is not 1.22.5 or higher"
fi