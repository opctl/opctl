# httprs
A ReadSeeker for http.Response.Body

[![wercker status](https://app.wercker.com/status/b8ab18faefae7d1f88f9f23d642f0847/s/master "wercker status")](https://app.wercker.com/project/bykey/b8ab18faefae7d1f88f9f23d642f0847)

## Usage

```
import "github.com/jfbus/httprs"

resp, err := http.Get(url)

rs := httprs.NewHttpReadSeeker(resp)
defer rs.Close()
io.ReadFull(rs, buf) // reads the first bytes from the response
rs.Seek(1024, 0) // moves the position
io.ReadFull(rs, buf) // does an additional range request and reads the first bytes from the second response
```
if you use a specific http.Client :
```
rs := httprs.NewHttpReadSeeker(resp, client)
```

## Doc

See http://godoc.org/github.com/jfbus/httprs

## LICENSE

MIT - See LICENSE