# Recoil

Recoil is a testing tool for [Ricochet](https://ricochet.im) which allows you
to ping, send contact requests and send arbitrary messages to a ricochet client
either through Tor or locally.


## Authenticated Ping

If you want to know whether a particular client is online, you can use the 
`ping` action. This connects to and completes an authentication with the given 
service.

        ./bin/recoil --target <target_hostname> --hostname <your_host_name> -key <private_key_file> -action ping

## Contact Request

Before you can start sending messages you will need to issue a contact request 
to the target. You can do this by using the `contact-request` action.

        ./bin/recoil --target <target_hostname> --hostname <your_host_name> -key <private_key_file> -action contact-request

## Sending Messages

To send messages; create a message file (see sample-ricochet-messages for an 
example) and use the `send-messages` action.

        ./bin/recoil --target <target_hostname> --hostname <your_host_name> -key <private_key_file> -action send-messages --messageFile <file>
