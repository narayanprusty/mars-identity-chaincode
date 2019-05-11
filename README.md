# mars-identity-chaincode

This chaincode is used to assign unique identity every person in mars. This chaincode will be deployed on the "identity" channel by the identity authority.

## Install and Instantiate 

First ssh into the EC2 that's running the containers. Then access to shell of containers using this command: `docker exec -i -t container_id /bin/bash`. 

Then follow this steps to install and instantiate the chaincode:

1. Create the directory using `mkdir /opt/gopath/src && mkdir /opt/gopath/src/github.com` if doesn't exist.
2. We need to install the chaincode in all three authorities containers. For all authorities, clone the chaincode repo using the command `cd /opt/gopath/src/github.com && git clone https://github.com/narayanprusty/mars-identity-chaincode.git`.
3. Then for all authorities install using this command: `peer chaincode install -n identity -v 1.0 -p github.com/mars-identity-chaincode`
4. Then in identity authority container run the following command to instantiate the chaincode: `peer chaincode instantiate -o $ORDERER_URL -C identity -n identity -v 1.0 -c '{"Args":[]}' --cafile /home/crypto/managedblockchain-tls-chain.pem --tls`

