package addresspool

import (
	"net/url"

	"github.com/go-chassis/openlog"

	"github.com/go-chassis/cari/discovery"
)

func getProtocolMap(eps []string) map[string]string {
	m := make(map[string]string)
	for _, ep := range eps {
		u, err := url.Parse(ep)
		if err != nil {
			openlog.Error("url err: " + err.Error())
			continue
		}
		m[u.Scheme] = u.Host
	}
	return m
}

func getAzAddressMap(instances []*discovery.MicroServiceInstance) map[string][]string {
	azAddrMap := make(map[string][]string) // key: az, value: address list

	for _, instance := range instances {
		azName := "unknown"
		if instance.DataCenterInfo != nil && len(instance.DataCenterInfo.AvailableZone) > 0 {
			azName = instance.DataCenterInfo.AvailableZone
		}

		m := getProtocolMap(instance.Endpoints)
		if ep := m["rest"]; len(ep) > 0 {
			azAddrMap[azName] = append(azAddrMap[azName], ep)
		}
	}
	return azAddrMap
}

func removeDuplicates(input []string) []string {
	if len(input) == 0 {
		return input
	}

	tmpMap := make(map[string]struct{})
	output := make([]string, 0)
	for _, v := range input {
		if _, ok := tmpMap[v]; ok {
			continue
		}
		tmpMap[v] = struct{}{}
		output = append(output, v)

	}
	return output
}
