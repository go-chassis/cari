package addresspool

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-chassis/cari/discovery"
)

func TestNewPool(t *testing.T) {
	mockHttpServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		return
	}))

	os.Setenv("CHASSIS_SC_HEALTH_CHECK_INTERVAL", "1")

	defaultAddr := mockHttpServer.Listener.Addr().String()
	// can get an address
	pool := NewPool([]string{defaultAddr})
	addr := pool.GetAvailableAddress()
	assert.Equal(t, defaultAddr, addr)

	// check monitor started
	assert.NotEqual(t, statusAvailable, pool.status[defaultAddr]) // unavailable by default

	time.Sleep(2*time.Second + 100*time.Millisecond)
	assert.Equal(t, statusAvailable, pool.status[defaultAddr])

	mockHttpServer.Close()
	time.Sleep(2*time.Second + 100*time.Millisecond)
	assert.NotEqual(t, statusAvailable, pool.status[defaultAddr]) // the status should be unavailable again
}

func TestAddressPool_GetAvailableAddress_priority(t *testing.T) {
	p := NewPool([]string{})
	sameAzAddr := "127.0.0.1:30100"
	diffAzAddr := "127.0.0.1:30101"
	defaultAddr := "127.0.0.1:30102"
	tests := []struct {
		name  string
		preDo func()
		want  string
	}{
		{
			name: "no address, return empty",
			preDo: func() {
			},
			want: "",
		},
		{
			name: "same az address available, return same az address",
			preDo: func() {
				p.defaultAddress = []string{defaultAddr}
				p.sameAzAddress = []string{sameAzAddr}
				p.diffAzAddress = []string{diffAzAddr}
				p.status[sameAzAddr] = statusAvailable
				p.status[diffAzAddr] = statusAvailable
				p.status[defaultAddr] = statusAvailable
			},
			want: sameAzAddr,
		},
		{
			name: "diff az address available, return diff az address",
			preDo: func() {
				p.status[sameAzAddr] = statusUnavailable
				p.status[diffAzAddr] = statusAvailable
				p.status[defaultAddr] = statusAvailable
			},
			want: diffAzAddr,
		},
		{
			name: "same az/diff az address unavailable, return default address",
			preDo: func() {
				p.status[sameAzAddr] = statusUnavailable
				p.status[diffAzAddr] = statusUnavailable
				p.status[defaultAddr] = statusAvailable
			},
			want: defaultAddr,
		},
		{
			name: "all address unavailable, return default address",
			preDo: func() {
				p.status[defaultAddr] = statusUnavailable
			},
			want: defaultAddr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preDo()
			assert.Equalf(t, tt.want, p.GetAvailableAddress(), "GetAvailableAddress()")
		})
	}
}

func TestAddressPool_GetAvailableAddress_filter(t *testing.T) {
	unavailableAddr := "127.0.0.1:30100"
	availableAddr1 := "127.0.0.1:30101"
	availableAddr2 := "127.0.0.1:30102"
	p := NewPool([]string{unavailableAddr, availableAddr1, availableAddr2})
	p.status[unavailableAddr] = statusUnavailable
	p.status[availableAddr1] = statusAvailable
	p.status[availableAddr2] = statusAvailable
	// should filter out available address
	for i := 0; i < 10; i++ {
		addr := p.GetAvailableAddress()
		assert.NotEqual(t, unavailableAddr, addr)
		assert.True(t, addr == availableAddr1 || addr == availableAddr2)
	}
	// should do load balance
	assert.NotEqual(t, p.GetAvailableAddress(), p.GetAvailableAddress())
}

func TestAddressPool_SetAddressByInstances(t *testing.T) {
	p := NewPool([]string{"192.168.2.1:30100", "192.168.2.3:30100"}) // default address is of az2

	assert.Error(t, p.SetAddressByInstances(nil))

	instances := []*discovery.MicroServiceInstance{
		{
			Endpoints: []string{"rest://192.168.1.1:30100", "grpc://192.168.1.1:30101"},
			DataCenterInfo: &discovery.DataCenterInfo{
				Name:          "engine1",
				Region:        "cn",
				AvailableZone: "az1",
			},
		},
		{
			Endpoints: []string{"rest://192.168.1.2:30100", "grpc://192.168.1.2:30101"},
			DataCenterInfo: &discovery.DataCenterInfo{
				Name:          "engine1",
				Region:        "cn",
				AvailableZone: "az1",
			},
		},
		{
			Endpoints: []string{"rest://192.168.2.1:30100", "grpc://192.168.2.1:30101"},
			DataCenterInfo: &discovery.DataCenterInfo{
				Name:          "engine2",
				Region:        "cn",
				AvailableZone: "az2",
			},
		},
		{
			Endpoints: []string{"rest://192.168.2.2:30100", "grpc://192.168.2.2:30101"},
			DataCenterInfo: &discovery.DataCenterInfo{
				Name:          "engine2",
				Region:        "cn",
				AvailableZone: "az2",
			},
		},
	}
	err := p.SetAddressByInstances(instances)
	assert.NoError(t, err)
	assert.Equal(t, []string{"192.168.2.1:30100", "192.168.2.2:30100"}, p.sameAzAddress)
	assert.Equal(t, []string{"192.168.1.1:30100", "192.168.1.2:30100"}, p.diffAzAddress)
}

func TestAddressPool_checkConnectivity(t *testing.T) {
	server1 := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		return
	}))
	server1Addr := server1.Listener.Addr().String()

	server2 := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		return
	}))
	server2Addr := server2.Listener.Addr().String()

	// init, all address is available
	defaultAddr := "127.0.0.1:30000"
	p := NewPool([]string{defaultAddr, server1Addr, server2Addr})
	assert.NotEqual(t, statusAvailable, p.status[defaultAddr])
	assert.NotEqual(t, statusAvailable, p.status[server1Addr])
	assert.NotEqual(t, statusAvailable, p.status[server2Addr])

	// check connectivity, default address status should be unavailable, as it is fake
	p.checkConnectivity()
	assert.Equal(t, statusUnavailable, p.status[defaultAddr])
	assert.Equal(t, statusAvailable, p.status[server1Addr])
	assert.Equal(t, statusAvailable, p.status[server2Addr])

	// check connectivity, server addresses should be unavailable, as the servers are closed
	server1.Close()
	server2.Close()
	p.checkConnectivity()
	assert.Equal(t, statusUnavailable, p.status[defaultAddr])
	assert.Equal(t, statusUnavailable, p.status[server1Addr])
	assert.Equal(t, statusUnavailable, p.status[server2Addr])
}
