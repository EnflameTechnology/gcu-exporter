#!/bin/bash
#
# Copyright (c) 2022 Enflame. All Rights Reserved.
#

cp collector/gcu_usage.go collector/$1.go
sed -i s/gcu_usage/$1/g collector/$1.go
sed -i s/usageCollector/$2/g collector/$1.go
sed -i s/UsageCollector/$4/g collector/$1.go
sed -i s/_usage/$3/g collector/$1.go
sed -i "s/Gcu usage as reported by the device/$5/g" collector/$1.go
