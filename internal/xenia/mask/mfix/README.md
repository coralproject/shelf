

# mfix
`import "github.com/coralproject/shelf/internal/xenia/mask/mfix"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func Add(db *db.DB, msk mask.Mask) error](#Add)
* [func Get(fileName string) ([]mask.Mask, error)](#Get)
* [func Remove(db *db.DB, collection string) error](#Remove)


#### <a name="pkg-files">Package files</a>
[mfix.go](/src/github.com/coralproject/shelf/internal/xenia/mask/mfix/mfix.go) 





## <a name="Add">func</a> [Add](/src/target/mfix.go?s=830:870#L31)
``` go
func Add(db *db.DB, msk mask.Mask) error
```
Add inserts a mask for testing.



## <a name="Get">func</a> [Get](/src/target/mfix.go?s=519:565#L13)
``` go
func Get(fileName string) ([]mask.Mask, error)
```
Get retrieves a slice of mask documents from the filesystem for testing.



## <a name="Remove">func</a> [Remove](/src/target/mfix.go?s=1198:1245#L46)
``` go
func Remove(db *db.DB, collection string) error
```
Remove is used to clear out all the test masks from the collection.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
