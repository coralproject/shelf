# Xenia

Alpha Release

A flexible service layer that publishes endpoints against [mongo aggregation pipeline queries](https://docs.mongodb.org/manual/core/aggregation-introduction/).  

Configuration describing the endpoints and queries are stored in a mongo collection allowing for updates to the service layer without touching Go code or restarting the application.

![Xenia Coral](http://www.101-saltwater-aquarium.com/graphics/xenia.jpg)

## Installation

### Download source code

1) Make sure you're have [go 1.5 or later](https://golang.org/dl/) installed and your [environment set up](https://golang.org/doc/install).

2) Make sure your go vendor experiment flag is set (will be set by default in a couple of months...)

```
export GO15VENDOREXPERIMENT=1
```

_We recommend adding this to your ~/.bash_profile or other startup script as it will become default go behavior soon._

3) Get the source code:

```
go get github.com/coralproject/xenia
```

4) Tell xenia which database you want to use:

Make a copy of the dev.cfg and then edit your version to set the appropriate values. Finally source your edited cfg file to create and set the environment variables: 

```
source $GOPATH/src/github.com/coralproject/xenia/config/[thefile].cfg
```

These environment variables must be set before running any of the code.

_Be careful not to commit any database passwords back to the repo!!_

### Build the CLI tool

Xenia has a CLI tool that allows you to manage endpoints and perform other actions.

1) Build the tool:

```
cd $GOPATH/src/github.com/coralproject/xenia/cmd/xenia
go build
```

_Note: It is best to run with logging level 0 when using the xenia command:_

```
export XENIA_LOGGING_LEVEL=0
```

### Creating a Xenia database for the first time

If you are running Xenia on a Mongo database for the first time you will need the Xenia command line tool to perform these functions.

```
cd $GOPATH/src/github.com/coralproject/xenia/app/xenia
```

1) Run the `db create` command:

```
./xenia db create -f ./srcdb/database.json
```

_You must run this on a new database to create the collections and the proper set of indexes._

2) Load all the exiting queries:

```
./xenia query upsert -p ./scrquery
```

## Run the web service

1) To run the web service, build and run `/app/xenia`:

```
go build
./xenia
```

2) Use a proper web token: 

```
Authorization "Basic NmQ3MmU2ZGQtOTNkMC00NDEzLTliNGMtODU0NmQ0ZDM1MTRlOlBDeVgvTFRHWjhOdGZWOGVReXZObkpydm4xc2loQk9uQW5TNFpGZGNFdnc9"
```

Xenia is secured via an authorization token.  If you are using it through an application that provides this token (aka, Trust) then you're good to go.  

If you intend to hit endpoints through a browser, install an Addon/plugin/extension that will allow you to add headers to your requests.

### API calls

If you set the authorization header properly in your browser you can run the following endpoints:

1) Get a list of configured queries:

```
http://localhost:4000/1.0/query
```

2) Get the query set document for the `basic` query set:

```
http://localhost:4000/1.0/query/basic
```

3) Execute the query for the `basic` query set:

```
http://localhost:5000/1.0/query/basic/exec
```

4) Execute the query for the `basic_var` query set with variables:

```
http://localhost:5000/1.0/query/basic_var/exec?station_id=42021
```

## Query management

Using the Xenia command line tool you can manage query sets.

```
cd $GOPATH/src/github.com/coralproject/xenia/app/xenia
```

1) Get a list of saved queries:

```
./xenia query list
```

3) Look at the details of a query:

```
./xenia query get -n top_commenters_by_count
```

4) Execute a query:

```
./xenia query exec -n top_commenters_by_count
```

5) Add or update a query for use:

```
./xenia query upsert -p ./scrquery/test_basic_var.json
```

By convention, we store core query scripts in the [/xenia/cmd/xenia/scrquery](https://github.com/CoralProject/xenia/tree/master/cmd/xenia/scrquery) folder.  As we develop Coral features, store the .json files there so other members can use them.  Eventually, groups of query sets will be refactored to elsewhere's yet undefined.

```
cd $GOPATH/src/github.com/coralproject/xenia/cmd/xenia/scrquery
ls
```

#### Direct Mongo access (optional)

You can look in the db at existing queries:

```
mongo [flags to connect to your server]
use coral (or your databasename)
db.query_sets.find()
```

#### Example query_set

Here's a basic query set document containing two pipeline calls and using a variable called #station_id#:

```
{
   "name":"basic",
   "desc":"Shows a basic multi result query.",
   "enabled":true,
   "params":null,
   "queries":[
      {
         "name":"Basic",
         "type":"pipeline",
         "collection":"test_bill",
         "return":true,
         "scripts":[
            "{\"$match\": {\"station_id\" : \"#station_id#\"}}",
            "{\"$project\": {\"_id\": 0, \"name\": 1}}"
         ]
      },
      {
         "name":"Time",
         "type":"pipeline",
         "collection":"test_bill",
         "return":true,
         "has_date":true,
         "scripts":[
            "{\"$match\": {\"condition.date\" : {\"$gt\": \"ISODate(\\\"2013-01-01T00:00:00.000Z\\\")\"}}}",
            "{\"$project\": {\"_id\": 0, \"name\": 1}}",
            "{\"$limit\": 2}"
         ]
      }
   ]
}
```

This query once saved can be executed via the API:

```
http://[server]:[port]/1.0/query/basic?station_id=123123
```

For documentation of each field in a query set document please refer to the [models.go](pkg/query/models.go) source code file.

## API Authentication

The [auth](https://github.com/ardanlabs/kit/tree/master/auth) package provides API's for managing users who will be accessing the xenia API. This includes all the CRUD related support for users and authentication. There are two collections in MongoDB called `auth_users` and `auth_sessions` that contain API user information and authentication. The `auth_users` collection contains registered users and `auth_sessions` contain sessions that allows users to be active in the system.

### Users

A User is an entity that can be authenticated on the system and granted rights to the API. A user document has the following form:

```
{ 
    "_id" : ObjectId("5660bc6e16908cae692e0593"), 
    "public_id" : "d648d9d1-f3a7-4586-b64e-f8d61ca986fe", 
    "private_id" : "5d829805-d801-408e-b418-2e9055da244b", 
    "status" : NumberInt(1), 
    "full_name" : "TEST USER DON'T DELETE", 
    "email" : "bill@ardanstudios.com", 
    "password" : "$2a$10$CRoh/8Uex49hviQYDlDvruoQUO10QxVOU7O0UMliqGlXSySK4SZEq", 
    "is_deleted" : false, 
    "date_modified" : ISODate("2015-12-03T22:04:30.117+0000"), 
    "date_created" : ISODate("2015-12-03T22:04:30.117+0000")
}
```

From an authentication standpoint several fields from a User document are important:

**PublicID**  : This is the users public identifier and can be shared with the world. It provides a unique id for each user. It is used to lookup users from the database. This is a randomlu generated UUID.

**PrivateID** : This is the users private identifier and must not be shared with the world. It is used in conjunction with the user supplied password to create an encrypted password. To authenticate with a password you need the users password and this private id. This is a randomly generated UUID.

**Password**  : This is a hash value based on a user provided password string and the user's private identifier. These values are combined and encrypted to create a hash value that is stored in the user document as the password.

### Sessions

A Session is a document in the database tied to a User via their PublicID. Sessions provide a level of security for web tokens by giving them an expiration date and a lookup point for the user accessing the API. The SessionID is what is used to look up the User performing the web call. The SessionID is a randomly generated UUID. If the Session is active, then a PublicID lookup can be performed and authentication can take place. If the Session is expired, authentication failed immediately. A user can have several Session documents, and when this is the case, the latest document is used to check authencation.

```
{ 
    "_id" : ObjectId("5660bc6e16908cae692e0594"), 
    "session_id" : "6d72e6dd-93d0-4413-9b4c-8546d4d3514e", 
    "public_id" : "d648d9d1-f3a7-4586-b64e-f8d61ca986fe", 
    "date_expires" : ISODate("2016-12-02T22:04:30.282+0000"), 
    "date_created" : ISODate("2015-12-03T22:04:30.282+0000")
}
```

### Web tokens

Access to Xenia's web service API requires sending a web token on every request. HTTP `Basic Authorization` is being used:

```
Authorization: Basic WebToken
```

A web token is a value that is not stored in the database for any User but is a value that can be consistently generated by having a User document and a SessionID. It is made up of two parts, a SessionID and a Token which are concatinated together and then base64 encoded for use over HTTP:

```
base64Encode(SessionID:Token)
```

The Token is generated by using the PublicID, PrivateID and Email fields from the User document to create a Salt value that is then combined with the User supplied Password to create a signed SHA256 hash value. This is the Token value that can be consistenly re-created when all the same values are present. If any of the fields used in this Token change, the Token will be invalidated.

### Web token authentication

To make things as secure as possible, database lookups are performed as part of web token authentication. The user must keep their token secure.

Here are the steps to web token authentication:

	* Decode the web token and break it into its parts of SessionID and Token.
	* Retrieve the Session document for the provided SessionID and validate it has not expired.
	* Retrieve the User document from the PublicID in the Session document.
	* Validate the Token is valid by generating a new Token from the retrieved User document.

If any of these steps fail, authorization fails.

### Generating users and tokens

The [Kit](https://github.com/ardanlabs/kit/tree/master/cmd) repo has a command line tool for creating new users. The tooling currently allows you to look up users and add new users to the system.
 
## Concepts and Motivations

### Composition

Aggregation pipelines are _chain-able_ allowing for the output of one endpoint to be fed into the next. Xenia will provide a request syntax to allow for this, giving the requesting application another dimension of flexibility via query control.

Similarly, output documents from multiple pipelines can be _bundled_ together. This is particularly useful in the no-sql/document db paradigm in which joins are not natively supported.

### Restructuring of Team Dynamics

Xenia moves 100% of the query logic out of the application code. Front end devs, data analysis, and anyone else familiar with the simple, declarative mongo aggregation syntax can alter the service behavior. This removes the requirement for back end engineering and devops expertise from the process of refining the data requests.

Xenia's CLI tools allow anyone with a basic understanding of document database concepts and aggregation pipeline syntax to create or update endpoints.  (Once the web UI is complete updates to the pipelines will be even more convenient.) 

### Also, Welcome!

All software in this repo is Open Source, offered under the MIT license.

For more information about The Coral Project, please visit [our website](https://coralproject.net).  For more information about how this technology is used in our projects, please visit [the reef](https://github.com/coralproject/reef).

