# Statsig Go Server SDK - Edge Config Adapter

A first party Fastly Integration with Statsig Go SDK

# Setup

1. Get this Go Adapter
```
go get github.com/statsig-io/fastly-data-adapter-go
```
or 
```
require (
  github.com/statsig-io/fastly-data-adapter-go
)
```

2. Initialize Pass in the adapter to SDK

```
import (
  statsig "github.com/statsig-io/go-sdk"
  "github.com/statsig-io/fastly-data-adapter-go" 
)

	adapter := NewFastlyDataAdapter(fastlyKey, storeID, configSpecKey)
  statsigOptions := &Options{
    DataAdapter: &adapter
  }
  statsig.InitializeWithOptions("server-secret-key", statsigOptions)
```

