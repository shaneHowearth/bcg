# BCG

## Pre build actions

Before build takes place some configuration information needs to be created. A Twilio account, and a Gmail account are required for the default setup.

### Twilio account information required:
- TwilioAccountSID
- TwilioSenderNumber
- TwilioAuthToken

These items need to be stored in `bcg/notify/internal/externalProviders.env`

### GMail

`credentials.json` can be downloaded from `https://developers.google.com/gmail/api/quickstart/go`, click on the "Enable the Gmail API" and click "DOWNLOAD CLIENT CONFIGURATION".
The file needs to be saved to `bcg/notify/internal/emailprovider/gmail/credentials.json`,
The user then needs to `cd bcg/notify/internal/emailprovider/gmail/getToken` and run `go run main.go`.
This will create a `token.json` file in the directory above and is REQUIRED to allow emails to be sent via gmail.


Ensure that no process is currently listening on port 80 of the host machine.

## Build
With those required files the project can be built as follows:
In the project root directory (`bcg`) run `docker-compose up`.
This will create 3 (three) containers, `bcg_restserver`, `bcg_customer`, and `bcg_notify`.

The rest server will be listening on port 80, with 3 endpoints:
POST: /customer
GET: /customer/<id>
POST: /notify

See swagger.yaml for full explanation of endpoints.

## Notes
- For ease of development I have created a single container for the customer data path. Normally I would prefer to have a CQRS architecture where the Create Customer data path is a container that stores data in a (separately hosted) SQL database and uses a Message Queue to populate a NoSQL cache for the Query data path (Get Customer).
- To use a different email provider, create a package under `bcg/notify/internal/emailprovider`. Create a type that implements the `mailer.Mailer` interface (as defined in `bcg/notify/integration/email/email.go`). Then pass an instance of that to the `notify.Server` struct in `bcg/notify/cmd/main.go`
- To use a different notification provider, create a package under `bcg/notify/internal/sms`. Create a type that implements the `messenger.Messenger` interface (as defined in `bcg/notify/integration/sms/sms.go`). Then pass an instance of that to the `notify.Server` struct in `bcg/notify/cmd/main.go`

## Cleanup
The following commands will remove all of the docker containers when you are finished.
```
docker stop {bcg_notify_1,bcg_customer_1,bcg_restserver_1}
docker rm {bcg_notify_1,bcg_customer_1,bcg_restserver_1}
docker rmi {bcg_notify:latest,bcg_customer:latest,bcg_restserver:latest}
# Remove hanging images
# WARNING this will remove all hanging images on the host
# (Uncomment to use)
# docker image prune
```
