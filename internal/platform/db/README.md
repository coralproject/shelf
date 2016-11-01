

# db
`import "github.com/coralproject/shelf/internal/platform/db"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
Package db abstracts different database systems we can use.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func RegMasterSession(context interface{}, name string, url string, timeout time.Duration) error](#RegMasterSession)
* [type DB](#DB)
  * [func NewMGO(context interface{}, name string) (*DB, error)](#NewMGO)
  * [func (db *DB) BatchedQueryMGO(context interface{}, colName string, q bson.M) (*mgo.Iter, error)](#DB.BatchedQueryMGO)
  * [func (db *DB) BulkOperationMGO(context interface{}, colName string) (*mgo.Bulk, error)](#DB.BulkOperationMGO)
  * [func (db *DB) CloseCayley(context interface{})](#DB.CloseCayley)
  * [func (db *DB) CloseMGO(context interface{})](#DB.CloseMGO)
  * [func (db *DB) CollectionMGO(context interface{}, colName string) (*mgo.Collection, error)](#DB.CollectionMGO)
  * [func (db *DB) CollectionMGOTimeout(context interface{}, timeout time.Duration, colName string) (*mgo.Collection, error)](#DB.CollectionMGOTimeout)
  * [func (db *DB) ExecuteMGO(context interface{}, colName string, f func(*mgo.Collection) error) error](#DB.ExecuteMGO)
  * [func (db *DB) ExecuteMGOTimeout(context interface{}, timeout time.Duration, colName string, f func(*mgo.Collection) error) error](#DB.ExecuteMGOTimeout)
  * [func (db *DB) GraphHandle(context interface{}) (*cayley.Handle, error)](#DB.GraphHandle)
  * [func (db *DB) NewCayley(context interface{}, name string) error](#DB.NewCayley)


#### <a name="pkg-files">Package files</a>
[cayley.go](/src/github.com/coralproject/shelf/internal/platform/db/cayley.go) [db.go](/src/github.com/coralproject/shelf/internal/platform/db/db.go) [mongo.go](/src/github.com/coralproject/shelf/internal/platform/db/mongo.go) 



## <a name="pkg-variables">Variables</a>
``` go
var ErrGraphHandle = errors.New("Graph handle not initialized.")
```
ErrGraphHandle is returned when a graph handle is not initialized.

``` go
var ErrInvalidDBProvided = errors.New("Invalid DB provided")
```
ErrInvalidDBProvided is returned in the event that an uninitialized db is
used to perform actions against.



## <a name="RegMasterSession">func</a> [RegMasterSession](/src/target/mongo.go?s=1022:1118#L30)
``` go
func RegMasterSession(context interface{}, name string, url string, timeout time.Duration) error
```
RegMasterSession adds a new master session to the set. If no url is provided,
it will default to localhost:27017.




## <a name="DB">type</a> [DB](/src/target/db.go?s=391:525#L3)
``` go
type DB struct {
    // contains filtered or unexported fields
}
```
DB is a collection of support for different DB technologies. Currently
only MongoDB has been implemented. We want to be able to access the raw
database support for the given DB so an interface does not work. Each
database is too different.







### <a name="NewMGO">func</a> [NewMGO](/src/target/mongo.go?s=1634:1692#L55)
``` go
func NewMGO(context interface{}, name string) (*DB, error)
```
NewMGO returns a new DB value for use with MongoDB based on a registered
master session.





### <a name="DB.BatchedQueryMGO">func</a> (\*DB) [BatchedQueryMGO](/src/target/mongo.go?s=3185:3280#L113)
``` go
func (db *DB) BatchedQueryMGO(context interface{}, colName string, q bson.M) (*mgo.Iter, error)
```
BatchedQueryMGO returns an iterator capable of iterating over
all the results of a query in batches.




### <a name="DB.BulkOperationMGO">func</a> (\*DB) [BulkOperationMGO](/src/target/mongo.go?s=3537:3623#L125)
``` go
func (db *DB) BulkOperationMGO(context interface{}, colName string) (*mgo.Bulk, error)
```
BulkOperationMGO returns a bulk value that allows multiple orthogonal
changes to be delivered to the server.




### <a name="DB.CloseCayley">func</a> (\*DB) [CloseCayley](/src/target/cayley.go?s=1289:1335#L45)
``` go
func (db *DB) CloseCayley(context interface{})
```
CloseCayley closes a graph handle value.




### <a name="DB.CloseMGO">func</a> (\*DB) [CloseMGO](/src/target/mongo.go?s=2398:2441#L87)
``` go
func (db *DB) CloseMGO(context interface{})
```
CloseMGO closes a DB value being used with MongoDB.




### <a name="DB.CollectionMGO">func</a> (\*DB) [CollectionMGO](/src/target/mongo.go?s=3835:3924#L138)
``` go
func (db *DB) CollectionMGO(context interface{}, colName string) (*mgo.Collection, error)
```
CollectionMGO is used to get a collection value.




### <a name="DB.CollectionMGOTimeout">func</a> (\*DB) [CollectionMGOTimeout](/src/target/mongo.go?s=4116:4235#L147)
``` go
func (db *DB) CollectionMGOTimeout(context interface{}, timeout time.Duration, colName string) (*mgo.Collection, error)
```
CollectionMGOTimeout is used to get a collection value with a timeout.




### <a name="DB.ExecuteMGO">func</a> (\*DB) [ExecuteMGO](/src/target/mongo.go?s=2518:2616#L92)
``` go
func (db *DB) ExecuteMGO(context interface{}, colName string, f func(*mgo.Collection) error) error
```
ExecuteMGO is used to execute MongoDB commands.




### <a name="DB.ExecuteMGOTimeout">func</a> (\*DB) [ExecuteMGOTimeout](/src/target/mongo.go?s=2800:2928#L101)
``` go
func (db *DB) ExecuteMGOTimeout(context interface{}, timeout time.Duration, colName string, f func(*mgo.Collection) error) error
```
ExecuteMGOTimeout is used to execute MongoDB commands with a timeout.




### <a name="DB.GraphHandle">func</a> (\*DB) [GraphHandle](/src/target/cayley.go?s=1080:1150#L36)
``` go
func (db *DB) GraphHandle(context interface{}) (*cayley.Handle, error)
```
GraphHandle returns the Cayley graph handle for graph interactions.




### <a name="DB.NewCayley">func</a> (\*DB) [NewCayley](/src/target/cayley.go?s=462:525#L8)
``` go
func (db *DB) NewCayley(context interface{}, name string) error
```
NewCayley adds support to a DB value for cayley based on a registered
master cayley handle.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
