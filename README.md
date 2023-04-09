# Go-Blockchain
# pidasp - 2022

This is an example application implemented using the Hyperledger Fabric blockchain framework.

The "./fabric-samples" directory mainly contains the code from the official Hyperledger Fabric samples (https://hyperledger-fabric.readthedocs.io/en/release-2.2/install.html). Some of the code has been modified to fit the purposes of this project.

#Using the application

To start network ypu need to enter /fabric-samples/test-network and run "./network.sh createChannel -ca".
The network has 4 organizations, and each has 4 peers.
The option "-ca" ensures that all cryptometerials are generated using the Certificati Authority. If this is needed to overwrite default settings which are using sryptogen tool to generate cryptomaterials.

![alt text](Images/starting.png?raw=true)


#Deploying Chaincode
To deploy cheincode after the network is started, run "./network.sh deployCC -ccn basic -ccp ../MyProject/chaincode/ -ccl go", while in the fabric-samples/test-network directory.This will package, install and commit the chaincode on all 16 peers in the network (4 peers in 4 organizations). After deploying the chaincode, it will be approved by all 4 organizations. The chaincode name is "basic" by default since no name was specified when creating the channel. The source code for the smart contract that makes up the chaincode is written in Golang and is located in the MyProject/chaincode directory.


## Running the client application
To run the client application, enter the MyProject/client, and run "go run .". After this, web application is started and you can write requests in postman or in order to interact with the network. The application communicates with one of the peers from the organization that you choose after starting the application. In order to change which peer this is, open the project/cars-and-persons-application/app_config.json and change the desired fields.

#Online sources that helped in developing this project
https://www.youtube.com/watch?v=85XgeOJHky0&t=8972s
https://www.youtube.com/watch?v=w-d_Uio0jWA
