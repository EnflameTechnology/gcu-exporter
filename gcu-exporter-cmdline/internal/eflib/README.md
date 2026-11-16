# Enflame go-eflib user guide

go-eflib use Enflame Management Library (EFML) to manage enflame gcu information.


## go.mod

Add `go-eflib v1.0.0` in `go.mod`, for example:

```
module modulename

go 1.18

require (
	go-eflib v1.0.0
)

```

## vendor/modules.txt

Add `go-eflib` in `vendor/modules.txt`, for example:

```
# go-eflib v1.0.0
## explicit
go-eflib
go-eflib/efml

```


## import

import in go file, for example:

```
import (
	"go-eflib"
	"go-eflib/efml"
)
```