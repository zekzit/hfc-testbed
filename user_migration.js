var hfc = require("hfc");

console.log(" **** starting iMed <--> Blockchain user migration ****");

var PEER_ADDRESS         = process.env.CORE_PEER_ADDRESS || '192.168.9.88:7051';
var MEMBERSRVC_ADDRESS   = process.env.MEMBERSRVC_ADDRESS || '192.168.9.88:7054';

var chain, chaincodeID;

chain = hfc.newChain("vault_6");

// Initialize the place to store sensitive private key information
chain.setKeyValStore( hfc.newFileKeyValStore('/tmp/keyValStore') );

// Set the URL to membership services and to the peer
console.log("member services address ="+MEMBERSRVC_ADDRESS);
console.log("peer address ="+PEER_ADDRESS);
chain.setMemberServicesUrl("grpc://"+MEMBERSRVC_ADDRESS);
chain.addPeer("grpc://"+PEER_ADDRESS);

var mode = process.env['DEPLOY_MODE'] || 'dev';
console.log("DEPLOY_MODE = " + mode);

if (mode === 'dev') {
    chain.setDevMode(true);
    //Deploy will not take long as the chain should already be running
    chain.setDeployWaitTime(10);
} else {
    chain.setDevMode(false);
    //Deploy will take much longer in network mode
    chain.setDeployWaitTime(120);
}

chain.setInvokeWaitTime(10);

// Begin by enrolling the user
enroll();

// Enroll a user.
function enroll() {
   console.log("enrolling user admin ...");
   // Enroll "admin" which is preregistered in the membersrvc.yaml
   chain.enroll("admin", "Xurw3yU9zI0l", function(err, admin) {
      if (err) {
         console.log("ERROR: failed to register admin: %s",err);
         process.exit(1);
      }
      // Set this user as the chain's registrar which is authorized to register other users.
      chain.setRegistrar(admin);

      var userName = "JohnDoe";
      // registrationRequest
      var registrationRequest = {
          enrollmentID: userName,
          affiliation: "bank_a"
      };
      chain.registerAndEnroll(registrationRequest, function(error, user) {
          if (error) throw Error(" Failed to register and enroll " + userName + ": " + error);
          console.log("Enrolled %s successfully\n", userName);
          deploy(user);
      });
   });
}
