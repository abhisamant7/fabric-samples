#!/bin/bash

if [ -z $SUDO_USER ]
then
    echo "===== Script need to be executed with sudo ===="
    echo "Change directory to 'setup'"
    echo "Usage: sudo ./caserver.sh"
    exit 0
fi

source ./to_absolute_path.sh

echo "=======Set up go (Takes a loooong time - get a coffee :)======"
# apt-get update
# UPDATED: 1.4-1.4  December 15, 2019
# apt-get -y install golang-1.10-go
# mv /usr/lib/go-1.10 /usr/local
# mv /usr/local/go-1.10   /usr/local/go

curl -O https://dl.google.com/go/go1.12.14.linux-amd64.tar.gz
tar -zvxf go1.12.14.linux-amd64.tar.gz -C /usr/local 
rm go1.12.14.linux-amd64.tar.gz

# tar -zxvf archive.tar.gz /folder/in/archive

# If GOROOT already set then DO Not set it again
if [ -z $GOROOT ]
then
    echo "export GOROOT=/usr/local/go" >> ~/.profile
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile

    GOPATH=$PWD/../../gocc
    to-absolute-path $GOPATH
    GOPATH=$ABS_PATH

    echo "export GOPATH=$GOPATH" >> ~/.profile
    echo "======== Updated .profile with GOROOT/GOPATH/PATH===="

    echo "export GOROOT=/usr/local/go" >> ~/.bashrc
    echo "export GOPATH=$GOPATH" >> ~/.bashrc
    echo "======== Updated .profile with GOROOT/GOPATH/PATH===="

    # UPDATED Dec 15, 2019
    echo "export GOCACHE=~/.go-cache" >> ~/.bashrc


else
    echo "======== No Change made to .profile ====="
fi

echo "======= Done. PLEASE LOG OUT & LOG Back In ===="
echo "Then validate by executing    'go version'"