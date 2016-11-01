# Shelf

Shelf provides a number of open source services that can be used to power community (comments, forms, etc.) on a site (e.g., a blog or newsroom site).  More specifically, [Shelf](https://github.com/coralproject/shelf) is a configurable service layer that publishes endpoints against functionality for: 
- [mongo aggregation pipeline queries](https://docs.mongodb.org/manual/core/aggregation-introduction/).
- importing, formatting, and validation of data.
- management of relationships between pieces of data.
- management and generation of views of imported data.

The individual services providing these functionalities can be built from the `cmd` directory as daemons (powering JSON APIs) or CLI tools.  The services are as follows:

- `xenia`/`xeniad` - allows a user to query a collection or a "view" of data via mongo aggregation pipeline queries
- `sponge`/`sponged` - allows a user to import data into a "items" collection and infer relationships between the imported data and other data in the "items" collection.
- `wire` - allows a user to create a "view" of data based on inferred relationships.
- `askd` - allows a user to create forms and manage submissions to those forms.
- `corald` - proxies some of the endpoints exposed by the above services and serves as a central point for interaction with the platform services (assuming you are running more than one of the above).

All of the Shelf documentation (including installation instructions) can be found in the [Coral Project Documentation](https://coralprojectdocs.herokuapp.com/).
