package cmdflags

import "time"

var VersionName, VersionCommit, VersionBuilt = "0.0.0", "0000000000", time.Now().In(time.UTC).Format(time.RFC3339)
