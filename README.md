# gorsa

**gorsa** is a library written in Golang that helps you with the **RSA Security Analytics REST API**.

## Usage

```go
import "github.com/fallais/gorsa"
```

Construct a new QRadar client, then use the various services on the client to
access different parts of the QRadar API. For example:

```go
client := gorsa.NewClient(nil)
```

If you want to provide your own `http.Client`, you can do it :

```go
httpClient := &http.Client{}
client := gorsa.NewClient(httpClient)
```
