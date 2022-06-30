package goapstra

import "strings"

type AgentPackages map[string]string

func (o *AgentPackages) raw() rawAgentPackages {
	// todo: one of these lines causes 'null' in JSON output, while the other causes '[]' ... which one is correct for the API?
	//raw := rawAgentPackages{}
	var raw rawAgentPackages
	for k, v := range *o {
		raw = append(raw, k+apstraSystemAgentPlatformStringSep+v)
	}
	return raw
}

type rawAgentPackages []string

func (o *rawAgentPackages) polish() AgentPackages {
	var polish AgentPackages
	if len(*o) > 0 {
		polish = make(map[string]string)
	}
	for _, s := range *o {
		kv := strings.SplitN(s, apstraSystemAgentPlatformStringSep, 2)
		switch len(kv) {
		case 2:
			polish[kv[0]] = kv[1]
		case 1:
			polish[kv[0]] = ""
		}
	}
	return polish
}
