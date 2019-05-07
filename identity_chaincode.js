const shim = require('fabric-shim');
const util = require('util');
const ethGSV = require('ethereum-gen-sign-verify');

var Chaincode = class {

  // Initialize the chaincode
  async Init(stub) {
    try {
      await stub.putState("identityAuthority", Buffer.from(stub.getCreator().mspid));
      return shim.success();
    } catch (err) {
      return shim.error(err);
    }
  }

  async Invoke(stub) {
    let ret = stub.getFunctionAndParameters();
    let method = this[ret.fcn];
    
    if (!method) {
      console.log('no method of name:' + ret.fcn + ' found');
      return shim.error();
    }

    try {
      let payload = await method(stub, ret.params);
      return shim.success(payload);
    } catch (err) {
      console.log(err);
      return shim.error(err);
    }
  }

  //get the mspid of the identity authority
  async getCreatorIdentity(stub) {
    let creatorIdentity = await stub.getState('identityAuthority');

    if (!creatorIdentity) {
      throw new Error("Creator identity not found");
    }

    return creatorIdentity;
  }

  async issueIdentity(stub, args) {
    if (args.length < 3) {
      throw new Error('Incorrect number of arguments. Expecting owner public key, meta data hash and unique identifier');
    }

    let id = args[0];
    let publicKey = args[1];
    let metaDataHash = args[2]

    let user = {publicKey, metaDataHash}

    if(!(await stub.getState(`user_${id}`)).toString()) {
      await stub.putState(`user_${id}`, Buffer.from(JSON.stringify(user)))
    } else {
      throw new Error('ID already exists');
    }
  }

  async updateUserMetaDataHash() {
    if (args.length < 1) {
      throw new Error('Incorrect number of arguments. Hash and ID missing');
    }

    let id = args[0]
    let user = (await stub.getState(`user_${id}`)).toString();

    if(user) {
      if((await stub.getState('identityAuthority')).toString() === stub.getCreator().mspid) {
        user = JSON.parse(user)
        user.metaDataHash = args[1]
        await stub.putState(`user_${id}`, Buffer.from(JSON.stringify(user)))
      } else {
        throw new Error('You don\'t have permission to update hash');
      }
    } else {
      throw new Error('User not found');
    }
  }

  async getIdentity(stub, args) {
    if (args.length < 1) {
      throw new Error('User ID missing');
    }

    let user = await stub.getState(`user_${args[0]}`);

    if(!user) {
      throw new Error('User not found');
    }

    return user;
  }

  async addServiceProvider(stub, args) {
    if (args.length < 3) {
      throw new Error('Incorrect number of arguments. Expecting service provider MSP ID, name and public key');
    }

    if((await stub.getState('identityAuthority')).toString() === stub.getCreator().mspid) {
      let id = args[0];
      let name = args[1];
      let publicKey = args[2];

      let sp = {name, publicKey}

      if(!(await stub.getState(`sp_${id}`)).toString()) {
        await stub.putState(`sp_${id}`, Buffer.from(JSON.stringify(sp)))
      } else {
        throw new Error('ID already exists');
      }
    } else {
      throw new Error('You don\'t have permssion to add service provider');
    }
  }

  async getServiceProvider(stub, args) {
    if (args.length < 1) {
      throw new Error('MSP ID missing');
    }

    let sp = await stub.getState(`sp_${args[0]}`);

    if(!sp) {
      throw new Error('Service provider not found');
    }

    return sp;
  }

  async requestAccess(stub, args) {
    if (args.length < 3) {
      throw new Error('Incorrect number of arguments. User ID, message and signature is missing');
    }

    let user = (await stub.getState(`user_${args[0]}`)).toString();
    let mspid = stub.getCreator().mspid;
    let sp = (await stub.getState(`sp_${mspid}`)).toString()
    let msg = JSON.parse(args[1])
    let signature = JSON.parse(args[2])

    

    if(!user) {
      throw new Error('User not found');
    }

    if(msg.action !== 'grantAccess' || msg.to !== mspid) {
      throw new Error('Permission invalid');
    }

    if(sp) {
      user = JSON.parse(user)
      let isValid = ethGSV.verify(JSON.stringify(msg), signature, user.publicKey);

      if(isValid) {
        sp = JSON.parse(sp);
        stub.setEvent('grantAccess', Buffer.from(JSON.stringify({"publicKey": sp.publicKey, "mspid": mspid})))
      } else {
        throw new Error('Permission invalid');
      }
    } else {
      throw new Error('Service provider not found');
    }
  }
};

shim.start(new Chaincode());
