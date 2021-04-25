# example
```
    if env.IDC() == env.DC_HY {
        // ...
    }

    if env.Region() == env.R_CN {
        // ...
    }
```
# env
determine idc by ip:

```
	idcPrefix = map[string][]string{
		DC_HY:  []string{"10.4"},
		DC_LF:  []string{"10.2", "10.3", "10.6", "10.8", "10.9", "10.10", "10.11", "10.12"},
		DC_VA:  []string{"10.100"},
		DC_SG:  []string{"10.101"},
		DC_CA:  []string{"10.106"}, // West America
		DC_ALISG: []string{"10.115"}, // Singapore Aliyun
	}
```

determin region by idc:

```
	regionIDCs = map[string][]string{
		R_CN:  []string{DC_HY, DC_LF},
		R_SG:  []string{DC_SG},
		R_US:  []string{DC_VA},
		R_CA:  []string{DC_CA}, // West America
		R_ALISG: []string{DC_ALISG}, // Singapore Aliyun
	}
```