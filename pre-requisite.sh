# install the prerequisites
# sudo apt-get install git curl -y

# install go and jq
# sudo apt-get install golang-go jq -y

# Make sure the Docker daemon is running.
# sudo systemctl start docker

# Add your user to the Docker group.
# sudo usermod -a -G docker $USER

if [ -x "$(command -v git --version)" ]; then
    # jq --version
    echo "Passed: git" $(git --version)
else
    echo "Error: git is not installed"
    # stop the script
    exit 1
fi

if [ -x "$(command -v curl --version)" ]; then
    # jq --version
    echo "Passed: curl"
else
    echo "Error: curl is not installed"
    # stop the script
    exit 1
fi


# Check version numbers  
if [ -x "$(command -v docker --version)" ]; then
    echo "Passed: Docker" $(docker --version)
else
    echo "Error: docker is not installed"
    exit 1
fi


if [ -x "$(command -v jq --version)" ]; then
    # jq --version
    echo "Passed: jq" $(jq --version)
else
    echo "Error: jq is not installed"
    exit 1
fi


# get the go version number, and check if it is 1.22.5 or higher
if go version | grep -q "go1.22.[5-9]\|go1.2[3-9][0-9]\|go[2-9][0-9]\|go[1-9][0-9][0-9]"; then
    echo "Passed: Go" $(go version)
else
    echo "error: Go version is not 1.22.5 or higher"
    exit 1
fi



