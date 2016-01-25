
# script
    import "github.com/coralproject/xenia/pkg/script"




## Constants
``` go
const (
    Collection        = "query_scripts"
    CollectionHistory = "query_scripts_history"
)
```
Contains the name of Mongo collections.


## Variables
``` go
var (
    ErrNotFound = errors.New("Set Not found")
)
```
Set of error variables.


## func Delete
``` go
func Delete(context interface{}, db *db.DB, name string) error
```
Delete is used to remove an existing Set document.


## func GetByNames
``` go
func GetByNames(context interface{}, db *db.DB, names []string) ([]Script, error)
```
GetByNames retrieves the documents for the specified names.


## func GetNames
``` go
func GetNames(context interface{}, db *db.DB) ([]string, error)
```
GetNames retrieves a list of script names.


## func GetScripts
``` go
func GetScripts(context interface{}, db *db.DB, tags []string) ([]Script, error)
```
GetScripts retrieves a list of scripts.


## func Upsert
``` go
func Upsert(context interface{}, db *db.DB, scr *Script) error
```
Upsert is used to create or update an existing Script document.



## type Script
``` go
type Script struct {
    Name     string   `bson:"name" json:"name" validate:"required,min=3"` // Unique name per Script document
    Commands []string `bson:"commands" json:"commands"`                   // Commands to add to a query.
}
```
Script contain pre and post commands to use per set or per query.









### func GetByName
``` go
func GetByName(context interface{}, db *db.DB, name string) (*Script, error)
```
GetByName retrieves the document for the specified name.


### func GetLastHistoryByName
``` go
func GetLastHistoryByName(context interface{}, db *db.DB, name string) (*Script, error)
```
GetLastHistoryByName gets the last written Script within the history.




### func (\*Script) Validate
``` go
func (scr *Script) Validate() error
```
Validate checks the query value for consistency.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)