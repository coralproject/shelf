**Note: This is an early stage project under active development. It is going to change quite a bit before it is ready. Consider it pre-alpha.**

# Xenia

Backend applications from The Coral Project

### Welcome!

All software in this repo is Open Source, offered under the MIT license.

For more information about The Coral Project, please visit [our website](https://coralproject.net).  For more information about how this technology is used in our projects, please visit [the reef](https://github.com/coralproject/reef).

Note: For expediency we are focusing all of our server-side development efforts in this single repo during initial development.  Once packages and dependencies become clear, we will implement concise apis and separate the various applications into their final homes. 

Note: The repo is under active development. Please browse our Issues and Pull Requests to get an idea of the state of our work.

If anything is unclear, please ask early and often! [Slack channel]  

## Xenia

![Xenia Coral](http://www.101-saltwater-aquarium.com/graphics/xenia.jpg)

Xenia is a flexible service layer that publishes endpoints against [mongo aggregation pipeline queries](https://docs.mongodb.org/manual/core/aggregation-introduction/).  

Configuration describing the endpoints and queries are stored in a mongo collection allowing for updates to the service layer without touching Go code or restarting the application.

### Quickstart

*Todo: write installation guide, link to Trust deployment in Xenia*

*Todo: provide example query_set document*

*Todo: describe Xenia's Auth paradigm*


### Composition

Aggregation pipelines are _chain-able_ allowing for the output of one endpoint to be fed into the next. Xenia will provide a request syntax to allow for this, giving the requesting application another dimension of flexibility via query control.

Similarly, output documents from multiple pipelines can be _bundled_ together. This is particularly useful in the no-sql/document db paradigm in which joins are not natively supported.

### Restructuring of Team Dynamics

Xenia moves 100% of the query logic out of the application code. Front end devs, data analysis, and anyone else familiar with the simple, declarative mongo aggregation syntax can alter the service behavior. This removes the requirement for back end engineering and devops expertise from the process of refining the data requests.

Xenia's CLI tools allow anyone with a basic understanding of document database concepts and aggregation pipeline syntax to create or update endpoints.  (Once the web UI is complete updates to the pipelines will be even more convenient.) 



## Trust Service Layer

The application publishes all endpoints that cannot be accomplished via Xenia Aggregation Pipelines.  

The primary job of the service layer is to expose CRUD functionality for all Trust specific data types:


### Users

Collection: users

A User is a member of the community.  They are the sole _actor_ type insomuch as they can create _content_ (aka, write comments) and perform _actions_ on _content_.

*Todo: define count caching strategy for all _actors_.*

*Todo: provide link to user model*

### Comments

Collection: comments

A comment is a basic type of _content_. As _content_ it can be _acted on_.  

*Todo: define count caching strategy for all _actionable_ content.*

*Todo: provide link to user model*

### Assets

Collection: assets

Addressable pieces of content that live outside of the Coral Ecosystem. The classic asset is an article.  _Content_ such as comments may be _on_ an asset.  In this sense, Assets are used to index conversation threads.

### Actions

Actions are carried out by _actors_ (aka, a User) on _content_ aka a comment.  

Actions include:

* A type (aka, 'like', 'recommend', 'flag')
* An actor (aka, user_id: 23456)
* A Target (aka, comment_id: 3242342)
* (optional) A value (aka: 5 (as in 5 stars))

### User Lists 

Collection: coral_lists

The primary use case of Trust is to allow community curators to build lists of users based on formulas.  

Schema:

* _id - bson id - required - provided by mongo
* name - string - required - identifying name, also used in api url
* curator - bson id - optional - placeholder for user
* formula - CoralFormula - instructions for calculating metrics (to be defined)

### Tags 

Collections: coral_tags, coral_tags_ref(?)

Curators may apply Tags to User and Comments. The list of available tags can be managed in the Settings interface.

(To Be Designed) - The association between a tag and an object with either be stored in a reference collection or a subdocument.  

### Notes (on users and comments)

Collections: coral_notes(?)

Curators may write Notes on Users and Comments. The list of available tags can be managed in the Settings interface.

(To Be Designed) - The association between a tag and an object with either be stored in a reference collection or a subdocument.  
	
