package addresspool

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-chassis/foundation/httpclient"
	"github.com/go-chassis/openlog"

	"github.com/go-chassis/cari/discovery"
)

// EnvCheckScInterval sc instance health check interval in second
const EnvCheckScInterval = "CHASSIS_SC_HEALTH_CHECK_INTERVAL"

const (
	statusAvailable   string = "available"
	statusUnavailable string = "unavailable"

	defaultCheckScIntervalInSecond = 15 // default sc instance health check interval in second
	healthProbeTimeout             = time.Second
)

const (
	ReadinessSuccess       = 0 // 连续两次成功
	ReadinessFailed        = 1 // 连续三次失败
	ReadinessIndeterminate = 2 // 其他状态
)

type HttpProbeOptions struct {
	Protocol string
	Path     string
}

type Options struct {
	HttpProbeOptions *HttpProbeOptions // used in available check if set, tcp will be used if not set
	DiffAzEndponits  []string
}

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

	httpProbeOptions *HttpProbeOptions
	httpProbeClient  *httpclient.Requests
	statusHistory    []map[string]string
}

func (p *Pool) Close() {
	p.onceQuit.Do(func() {
		close(p.quit)
	})
}

// NewPool Get registry pool instance
func NewPool(addresses []string, opts ...Options) *Pool {
	p := &Pool{
		defaultAddress: removeDuplicates(addresses),
		status:         make(map[string]string),
		statusHistory:  make([]map[string]string, 0, 4),
	}

	if len(opts) > 0 {
		if opts[0].HttpProbeOptions != nil {
			optCopy := *(opts[0].HttpProbeOptions)
			p.httpProbeOptions = &optCopy
			if len(p.httpProbeOptions.Protocol) == 0 {
				p.httpProbeOptions.Protocol = "http"
			}
			p.httpProbeClient, _ = httpclient.New(&httpclient.Options{
				TLSConfig:      &tls.Config{InsecureSkipVerify: true},
				RequestTimeout: 5 * time.Second,
			})
		}
		if len(opts[0].DiffAzEndponits) != 0 {
			p.diffAzAddress = opts[0].DiffAzEndponits
		}
	}
	
	p.monitor()
	return p
}

func (p *Pool) ResetAddress(addresses []string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.defaultAddress = removeDuplicates(addresses)
	p.diffAzAddress = []string{}
	p.sameAzAddress = []string{}
	p.status = make(map[string]string)
	p.statusHistory = make([]map[string]string, 0, 4)
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
			openlog.Info(fmt.Sprintf("sync same az endpoints: %s", uniqueAddrList))
			continue
		}
		p.diffAzAddress = uniqueAddrList
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

func (p *Pool) CheckReadiness() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	statusHistory := p.statusHistory

	if len(statusHistory) < 2 {
		return ReadinessIndeterminate
	}

	successCnt := 0
	failedCnt := 0
	for _, status := range statusHistory {
		if !existAvailableEndpointInStatus(status) {
			successCnt = 0
			failedCnt++
			continue
		}
		successCnt++
		failedCnt = 0
	}
	if successCnt >= 2 {
		return ReadinessSuccess
	}
	if failedCnt == 3 {
		return ReadinessFailed
	}

	return ReadinessIndeterminate
}

func existAvailableEndpointInStatus(status map[string]string) bool {
	for _, v := range status {
		if v == statusAvailable {
			return true
		}
	}
	return false
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
		err := p.doCheckConnectivity(v)
		if err != nil {
			openlog.Error(fmt.Sprintf("%s connectivity unavailable: %s", v, err))
			status[v] = statusUnavailable
		} else {
			status[v] = statusAvailable
		}
	}

	p.mutex.Lock()
	p.status = status
	p.statusHistory = append(p.statusHistory, status)
	cnt := len(p.statusHistory)
	if cnt > 3 {
		p.statusHistory = p.statusHistory[(cnt - 3):]
	}
	p.mutex.Unlock()
}

func (p *Pool) doCheckConnectivity(endpoint string) error {
	if p.httpProbeOptions != nil {
		return p.doCheckConnectivityWithHttp(endpoint)
	}

	return p.doCheckConnectivityWithTcp(endpoint)
}

func (p *Pool) doCheckConnectivityWithTcp(endpoint string) error {
	conn, err := net.DialTimeout("tcp", endpoint, healthProbeTimeout)
	if err != nil {
		return err
	}

	err = conn.Close()
	if err != nil {
		openlog.Error(fmt.Sprintf("close conn failed when check connectivity: %s", err))
	}

	return nil
}

func (p *Pool) doCheckConnectivityWithHttp(endpoint string) error {
	u := p.httpProbeOptions.Protocol + "://" + endpoint + p.httpProbeOptions.Path
	resp, err := p.httpProbeClient.Get(context.Background(), u, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return nil
	}
	// do tcp check if api not exist, ensure to compatible with old scenes
	if resp.StatusCode == http.StatusNotFound {
		return p.doCheckConnectivityWithTcp(endpoint)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("http status: %s, read resp error: %s", resp.Status, err)
	}
	err = resp.Body.Close()
	if err != nil {
		openlog.Error(fmt.Sprintf("close http resp.Body failed when check connectivity: %s", err))
	}
	return fmt.Errorf("http status: %s, resp: %s", resp.Status, string(body))
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

		p.checkConnectivity()
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
