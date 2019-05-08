package main

import (
	"fmt"
	"encoding/json"
	"encoding/hex"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/golang/protobuf/proto"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type User struct {
	PublicKey	string `json:"publicKey"`
	MetadataHash string `json:"metadataHash"`
	Permissions []string `json:"permissions"`
}

type ServiceProvider struct {
	Name	string `json:"name"`
	PublicKey string `json:"publicKey"`
}

type IdentityChaincode struct {
}

func (t *IdentityChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	var err error
	var identity []byte

	identity, err = stub.GetCreator()
	
	if err != nil {
		return shim.Error("An error occured")
	}

	sId := &msp.SerializedIdentity{}
	err = proto.Unmarshal(identity, sId)
	
	if err != nil {
			return shim.Error("An error occured")
	}

	nodeId := sId.Mspid
	err = stub.PutState("identityAuthority", []byte(nodeId))

	if err != nil {
		return shim.Error("An error occured")
	}

	return shim.Success(nil)
}

func (t *IdentityChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "getCreatorIdentity" {
		return t.getCreatorIdentity(stub, args)
	} else if function == "issueIdentity" {
		return t.issueIdentity(stub, args)
	} else if function == "getIdentity" {
		return t.getIdentity(stub, args)
	} else if function == "updateUserMetadataHash" {
		return t.updateUserMetadataHash(stub, args)
	} else if function == "addServiceProvider" {
		return t.addServiceProvider(stub, args)
	} else if function == "getServiceProvider" {
		return t.getServiceProvider(stub, args)
	} else if function == "requestAccess" {
		return t.requestAccess(stub, args)
	}

	return shim.Error("Invalid function name: " + function)
}

func (t *IdentityChaincode) getCreatorIdentity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	identity, err := stub.GetState("identityAuthority")

	if err != nil {
		return shim.Error("An error occured")
	}

	if identity == nil {
		return shim.Error("Identity not yet stored")
	}

	return shim.Success(identity)
}

func (t *IdentityChaincode) issueIdentity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	var identity, identityAuthority, newUserJson []byte
	var err error

	identityAuthority, err = stub.GetState("identityAuthority")

	if err != nil {
		return shim.Error("An error occured")
	}

	identity, err = stub.GetCreator()

	if err != nil {
		return shim.Error("An error occured")
	}

	sId := &msp.SerializedIdentity{}
	err = proto.Unmarshal(identity, sId)
	
	if err != nil {
			return shim.Error("An error occured")
	}

	nodeId := sId.Mspid

	if string(identityAuthority) != nodeId {
		return shim.Error("You are not authorized")
	}

	var newUser User
	newUser.PublicKey = args[1]
	newUser.MetadataHash = args[2]

	newUserJson, err = json.Marshal(newUser)

	if err != nil {
			return shim.Error("An error occured")
	}

	err = stub.PutState("user_" + args[0], []byte(string(newUserJson)))

	if err != nil {
		return shim.Error("An error occured")
	}

	return shim.Success(nil)
}

func (t *IdentityChaincode) getIdentity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	user, err := stub.GetState("user_" + args[0])

	if err != nil {
		return shim.Error("An error occured")
	}

	return shim.Success([]byte(user))
}

func (t *IdentityChaincode) updateUserMetadataHash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	var identity, identityAuthority []byte
	var err error

	identityAuthority, err = stub.GetState("identityAuthority")

	if err != nil {
		return shim.Error("An error occured")
	}

	identity, err = stub.GetCreator()

	if err != nil {
		return shim.Error("An error occured")
	}

	sId := &msp.SerializedIdentity{}
	err = proto.Unmarshal(identity, sId)
	
	if err != nil {
			return shim.Error("An error occured")
	}

	nodeId := sId.Mspid

	if string(identityAuthority) != nodeId {
		return shim.Error("You are not authorized")
	}

	user, err := stub.GetState("user_" + args[0])

	var userStruct User
	err = json.Unmarshal(user, &userStruct)

	if err != nil {
			return shim.Error("An error occured")
	}

	userStruct.MetadataHash = args[1]

	userJson, err := json.Marshal(userStruct)

	err = stub.PutState("user_" + args[0], []byte(string(userJson)))

	if err != nil {
		return shim.Error("An error occured")
	}

	return shim.Success(nil)
}

func (t *IdentityChaincode) addServiceProvider(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	var identity, identityAuthority, newSPJson []byte
	var err error

	identityAuthority, err = stub.GetState("identityAuthority")

	if err != nil {
		return shim.Error("An error occured")
	}

	identity, err = stub.GetCreator()

	if err != nil {
		return shim.Error("An error occured")
	}

	sId := &msp.SerializedIdentity{}
	err = proto.Unmarshal(identity, sId)
	
	if err != nil {
			return shim.Error("An error occured")
	}

	nodeId := sId.Mspid

	if string(identityAuthority) != nodeId {
		return shim.Error("You are not authorized")
	}

	var newSP ServiceProvider
	newSP.Name = args[1]
	newSP.PublicKey = args[2]

	newSPJson, err = json.Marshal(newSP)

	if err != nil {
			return shim.Error("An error occured")
	}

	err = stub.PutState("sp_" + args[0], []byte(string(newSPJson)))

	if err != nil {
		return shim.Error("An error occured")
	}

	return shim.Success(nil)
}

func (t *IdentityChaincode) getServiceProvider(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	sp, err := stub.GetState("sp_" + args[0])

	if err != nil {
		return shim.Error("An error occured")
	}

	return shim.Success([]byte(sp))
}

func (t *IdentityChaincode) trimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

func (t *IdentityChaincode) requestAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	var err error
	var identity []byte

	user, err := stub.GetState("user_" + args[0])

	var userStruct User
	err = json.Unmarshal(user, &userStruct)

	if err != nil {
		return shim.Error("An error occured")
	}

	identity, err = stub.GetCreator()

	if err != nil {
		return shim.Error("An error occured")
	}

	sId := &msp.SerializedIdentity{}
	err = proto.Unmarshal(identity, sId)
	
	if err != nil {
		return shim.Error("An error occured")
	}

	nodeId := sId.Mspid

	sp, err := stub.GetState("sp_" + nodeId)

	if err != nil {
		return shim.Error("An error occured")
	}

	var spStruct ServiceProvider
	err = json.Unmarshal(sp, &spStruct)

	if err != nil {
		return shim.Error("An error occured")
	}

	publicKeyBytes, err := hex.DecodeString(userStruct.PublicKey)

	if err != nil {
		return shim.Error("An error occured")
	}

	data := []byte("{\"action\":\"grantAccess\",\"to\":\"" + nodeId + "\"}")
	hash := crypto.Keccak256Hash(data)
	signature := hexutil.MustDecode(args[1])
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)

	if err != nil {
		return shim.Error("Failed to generate public key")
	}

	if hex.EncodeToString(publicKeyBytes) != t.trimLeftChars(hex.EncodeToString(sigPublicKey), 2){
		return shim.Error("Signature invalid")
	} 

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(IdentityChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
