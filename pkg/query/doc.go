// Package query provides API's for managing querysets which will be used in
// executing different aggregation tests against their respective data collection.
//
// QuerySet
// In query, records are required to follow specific formatting and are at this
// point, only allowed to be in a json serializable format which meet the query.Set
// structure.
//
// The query set execution supports the following types:
//
// - Pipeline
//      Pipeline query set types take advantage of MongoDB's aggregation API
//    (the currently supported data backend), which allows insightful use of its
//    internal query language, in providing context against data sets within the database.
//
//
//
//
// QuerySet Sample:
//
// ```json
// {
//   "name":"spending_advice",
//   "description":"tests against user spending rate and provides adequate advice on saving more",
//   "enabled": true,
//   "params":[
//     {
//       "name":"user_id",
//       "default":"396bc782-6ac6-4183-a671-6e75ca5989a5",
//       "desc":"provides the user_id to check against the collection"
//     }
//   ],
//   "rules": [
//   {
//     "desc":"match spending rate over 20 dollars",
//     "type":"pipeline",
//     "continue": true,
//     "script_options": {
//       "collection":"demo_user_transactions",
//       "has_date":false,
//       "has_objectid": false
//     },
//     "save_options": {
//       "save_as":"high_dollar_users",
//       "variables": true,
//       "to_json": true
//     },
//     "var_options":{},
//     "scripts":[
//       "{ \"$match\" : { \"user_id\" : \"#userId#\", \"category\" : \"gas\" }}",
//       "{ \"$group\" : { \"_id\" : { \"category\" : \"$category\" }, \"amount\" : { \"$sum\" : \"$amount\" }}}",
//       "{ \"$match\" : { \"amount\" : { \"$gt\" : 20.00} }}"
//     ]
//    }]
// }
//```
//
package query
