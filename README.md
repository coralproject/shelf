# Xenia

Alpha Release

A flexible service layer that publishes endpoints against [mongo aggregation pipeline queries](https://docs.mongodb.org/manual/core/aggregation-introduction/).  

Configuration describing the endpoints and queries are stored in a mongo collection allowing for updates to the service layer without touching Go code or restarting the application.


![Xenia Coral](http://www.101-saltwater-aquarium.com/graphics/xenia.jpg)


### Quickstart

#### Installation

1) Make sure you're have [go 1.5 or later](https://golang.org/dl/) installed and your [environment set up](https://golang.org/doc/install).

2) Make sure your go vendor experiment flag is set (will be set by default in a couple of months...)

```
export GO15VENDOREXPERIMENT=1
```

_We recommend adding this to your ~/.bash_profile or other startup script_ as it will become default go behavior soon.

3) Get the source code:

```
go get github.com/coralproject/xenia
```

4) Tell xenia which database you want to use:

Edit one of the .cfg files in /config/, then:

```
source $GOPATH/src/github.com/coralproject/xenia/config/[thefile].cfg
```

_Be careful not to commit any database passwords back to the repo!!_


#### Run the web server

1) To run the web server, build and run /app/xenia:

```
cd $GOPATH/src/github.com/coralproject/xenia/app/xenia
go build
./xenia
```

2) Xenia is secured via an Authorization token.  If you are using it through an application that provides this token (aka, Trust) then you're good to go.  

If you intend to hit endpoint through a browser, install an Addon/plugin/extension that will add headers to your requests. 

```
Authorization "Basic NmQ3MmU2ZGQtOTNkMC00NDEzLTliNGMtODU0NmQ0ZDM1MTRlOlBDeVgvTFRHWjhOdGZWOGVReXZObkpydm4xc2loQk9uQW5TNFpGZGNFdnc9"
```

#### Run the CLI tool (optional)

Xenia has a CLI tool that allows you to manage endpoints and perform other actions.

1) Build and run /cli/xenia:

```
cd $GOPATH/src/github.com/coralproject/xenia/cmd/xenia
go build
./xenia
```

2) ./xenia --help will take you from there


### Publishing Endpoints

Xenia publishes http endpoints against mongodb aggregation pipeline commands.  These endpoints are read from a mongodb collection called _query\_sets_.

_If you are running xenia on a db for the first time, you will need to use the CLI tool to add endpoints.  Without endpoints, xenia just sort of sits in the corner and rusts._

#### Adding query sets

Querysets can be added through the cli tool like so:

```
./xenia query upsert -p ./scrquery/test_basic_var.json
```

By convention, we store core query scripts in the [/xenia/cmd/xenia/scrquery](https://github.com/CoralProject/xenia/tree/master/cmd/xenia/scrquery) folder.  As we develop Coral features, store the .json files there so other members can use them.  Eventually, groups query_sets will be refactored to elsewhere's yet undefined.


#### Viewing active query sets

You can view a list of query_sets via the cli tool like so:

```
./xenia query list
```

Or you can just look in the db at them in raw form thiswise:

```
mongo [flags to connect to your server]
use coral (or your databasename)
db.query_sets.find()
```

By convention, we store working queries in .json files.  You can find them [here](cmd/xenia/scrquery/):

```
cd $GOPATH/src/github.com/coralproject/xenia/cmd/xenia/scrquery
ls
```

#### Example query_set

Here's a basic query_set configuration containing two pipeline calls and using a variable called #station_id#.


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

This query would be published to:

```
http://[server]:[port]/1.0/query/basic?station_id=123123
```

For living documentation of each of the parameters, please see [/pkg/query/main.go](pkg/query/main.go).


*Todo: describe Xenia's Auth paradigm*

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

