# Kronos-CLI

## Installation
To install the CLI, first download the binaries 
```
curl -LO https://storage.googleapis.com/kronos-cli/$(curl -L -s https://storage.googleapis.com/kronos-cli/stable.txt)/kronos-cli
```
To make the CLI shared between all users, execute the following 
```
chmod +x kronos-cli
sudo cp kronos-cli /usr/local/bin
```
To verify the installation of the CLI, run the following 
```
kronos-cli version
```