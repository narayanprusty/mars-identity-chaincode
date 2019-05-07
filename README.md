# mars-identity-chaincode

This chaincode allows identity authority to assign public key based unique ID to every person in mars. Every person is provided with an ID-card and the chip on the card carries the private key. Using the ID-card people can proof their identity in an electronic environment.

> The design process of the ID-card and scanning device is out of scope of the PoC. ID-card based identity solutions are already invented so it's not challenging to achieve this.

The meta data of identity owner such as name, age, DOB, address, mobile number and so on will be encrypted and stored in a proxy re-encryption server. 

Every person will be assigned a unique PRE (proxy re-encryption) private key using which the meta data of the person will be encrypted and stored in PRE server. The private key stored in the ID-card is used to proof ownership and grant/revoke data access to other governments/enterprises in the blockchain network. The PRE private key is stored with identity authority.   

