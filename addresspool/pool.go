package addresspool

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-chassis/openlog"

	"github.com/go-chassis/cari/discovery"
)

// EnvCheckScInterval sc instance health check interval in second
const EnvCheckScInterval = "CHASSIS_SC_HEALTH_CHECK_INTERVAL"

const (
	statusAvailable   string = "available"
	statusUnavailable string = "unavailable"

	defaultCheckScIntervalInSecond = 25 // default sc instance health check interval in second
	healthProbeTimeout             = time.Second
)

// Pool cloud server address pool
type Pool struct {
	mutex          sync.RWMutex
	defaultAddress []string
	// used when the server has addresses of multiple az
	// when we get available address, the priority is sameAzAddress > diffAzAddress > defaultAddress
	sameAzAddress []string
	diffAzAddress []string

	status      map[string]string
	onceMonitor sync.Once
	quit        chan struct{}
	onceQuit    sync.Once
}

func (p *Pool) Close() {
	p.onceQuit.Do(func() {
		close(p.quit)
	})
}

// NewPool Get registry pool instance
func NewPool(addresses []string) *Pool {
	p := &Pool{
		defaultAddress: removeDuplicates(addresses),
		status:         make(map[string]string),
	}
	p.appendAddressToStatus(addresses)
	p.monitor()
	return p
}

func (p *Pool) appendAddressToStatus(addresses []string) {
	for _, v := range addresses {
		if _, ok := p.status[v]; ok {
			continue
		}
		p.status[v] = statusAvailable
	}
}

func (p *Pool) ResetAddress(addresses []string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.defaultAddress = removeDuplicates(addresses)
	p.diffAzAddress = []string{}
	p.sameAzAddress = []string{}
	p.status = make(map[string]string)
	p.appendAddressToStatus(addresses)
}

func (p *Pool) SetAddressByInstances(instances []*discovery.MicroServiceInstance) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	azAddrMap := getAzAddressMap(instances)
	if len(azAddrMap) == 0 {
		return fmt.Errorf("sync endpoints failed")
	}

	for _, addrList := range azAddrMap {
		uniqueAddrList := removeDuplicates(addrList)
		if p.isSameAzAddr(uniqueAddrList) {
			p.sameAzAddress = uniqueAddrList
			p.appendAddressToStatus(uniqueAddrList)
			openlog.Info(fmt.Sprintf("sync same az endpoints: %s", uniqueAddrList))
			continue
		}
		p.diffAzAddress = uniqueAddrList
		p.appendAddressToStatus(uniqueAddrList)
		openlog.Info(fmt.Sprintf("sync different az endpoints: %s", addrList))
	}
	return nil
}

func (p *Pool) isSameAzAddr(addrList []string) bool {
	defaultAddrMap := make(map[string]struct{}, len(p.defaultAddress))
	for _, addr := range p.defaultAddress {
		defaultAddrMap[addr] = struct{}{}
	}
	for _, addr := range addrList {
		if _, exist := defaultAddrMap[addr]; exist {
			return true
		}
	}
	return false
}

// GetAvailableAddress Get an available address from pool by roundrobin
func (p *Pool) GetAvailableAddress() string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	addrs := p.getAvailableAddressList()
	if len(addrs) == 0 {
		addrs = p.defaultAddress
	}

	next := RoundRobin(addrs)
	addr, err := next()
	if err != nil {
		return ""
	}
	return addr
}

func (p *Pool) getAvailableAddressList() []string {
	if addrs := p.filterAvailableAddress(p.sameAzAddress); len(addrs) > 0 {
		return addrs
	}
	if addrs := p.filterAvailableAddress(p.diffAzAddress); len(addrs) > 0 {
		return addrs
	}
	if addrs := p.filterAvailableAddress(p.defaultAddress); len(addrs) > 0 {
		return addrs
	}

	return nil
}

func (p *Pool) filterAvailableAddress(addresses []string) []string {
	if len(addresses) == 0 {
		return nil
	}
	result := make([]string, 0)
	for _, v := range addresses {
		if p.status[v] == statusAvailable {
			result = append(result, v)
		}
	}
	return result
}

func (p *Pool) checkConnectivity() {
	toCheckedAddressList := make([]string, 0, len(p.defaultAddress)+len(p.sameAzAddress)+len(p.diffAzAddress))
	toCheckedAddressList = append(toCheckedAddressList, p.defaultAddress...)
	toCheckedAddressList = append(toCheckedAddressList, p.sameAzAddress...)
	toCheckedAddressList = append(toCheckedAddressList, p.diffAzAddress...)

	status := make(map[string]string) // create new map, to clear dirty address
	for _, v := range toCheckedAddressList {
		if _, exist := status[v]; exist {
			continue
		}
		conn, err := net.DialTimeout("tcp", v, healthProbeTimeout)
		if err != nil {
			openlog.Error("connectivity unavailable: " + v)
			status[v] = statusUnavailable
		} else {
			status[v] = statusAvailable
			conn.Close()
		}
	}

	p.mutex.Lock()
	p.status = status
	p.mutex.Unlock()
}

func (p *Pool) monitor() {
	p.onceMonitor.Do(func() {
		var interval time.Duration
		v, isExist := os.LookupEnv(EnvCheckScInterval)
		if !isExist {
			interval = defaultCheckScIntervalInSecond
		} else {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				interval = defaultCheckScIntervalInSecond
			} else {
				interval = time.Duration(i)
			}
		}
		ticker := time.NewTicker(interval * time.Second)
		p.quit = make(chan struct{})

		go func() {
			for {
				select {
				case <-ticker.C:
					p.checkConnectivity()
				case <-p.quit:
					ticker.Stop()
					return
				}
			}
		}()
	})
}
