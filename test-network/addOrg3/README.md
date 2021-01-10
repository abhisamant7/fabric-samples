## Adding Dfarmclient to the test network

You can use the `addOrg3.sh` script to add another organization to the Fabric test network. The `addOrg3.sh` script generates the Dfarmclient crypto material, creates an Dfarmclient organization definition, and adds Dfarmclient to a channel on the test network.

You first need to run `./network.sh up createChannel` in the `test-network` directory before you can run the `addOrg3.sh` script.

```
./network.sh up createChannel
cd addOrg3
./addOrg3.sh up
```

If you used `network.sh` to create a channel other than the default `mychannel`, you need pass that name to the `addorg3.sh` script.
```
./network.sh up createChannel -c channel1
cd addOrg3
./addOrg3.sh up -c channel1
```

You can also re-run the `addOrg3.sh` script to add Dfarmclient to additional channels.
```
cd ..
./network.sh createChannel -c channel2
cd addOrg3
./addOrg3.sh up -c channel2
```

For more information, use `./addOrg3.sh -h` to see the `addOrg3.sh` help text.
