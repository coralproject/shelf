

# web
`import "github.com/coralproject/shelf/cmd/sponge/web"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func Request(cmd *cobra.Command, verb, path string, body io.Reader) (string, error)](#Request)


#### <a name="pkg-files">Package files</a>
[web.go](/src/github.com/coralproject/shelf/cmd/sponge/web/web.go) 



## <a name="pkg-variables">Variables</a>
``` go
var DefaultClient request.Client
```
DefaultClient is the default client to perform downstream requests to.



## <a name="Request">func</a> [Request](/src/target/web.go?s=1861:1944#L56)
``` go
func Request(cmd *cobra.Command, verb, path string, body io.Reader) (string, error)
```
Request provides support for executing commands against the
web service.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
