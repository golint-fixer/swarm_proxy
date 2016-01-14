package lonely

import (
	"testing"

	"github.com/docker/swarm/cluster"
	"github.com/samalba/dockerclient"
	"github.com/samalba/dockerclient/mockclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockInfo = &dockerclient.Info{
		ID:              "id",
		Name:            "name",
		NCPU:            10,
		MemTotal:        20,
		Driver:          "driver-test",
		ExecutionDriver: "execution-driver-test",
		KernelVersion:   "1.2.3",
		OperatingSystem: "golang",
		Labels:          []string{"foo=bar"},
	}

	mockVersion = &dockerclient.Version{
		Version: "1.6.2",
	}
)

func TestInfo(t *testing.T) {

	// create mock client
	client := mockclient.NewMockClient()
	client.On("Info").Return(mockInfo, nil)
	client.On("Version").Return(mockVersion, nil)
	client.On("StartMonitorEvents", mock.Anything, mock.Anything, mock.Anything).Return()
	client.On("ListContainers", true, false, "").Return([]dockerclient.Container{}, nil).Once()
	client.On("ListImages", mock.Anything).Return([]*dockerclient.Image{}, nil)
	client.On("ListVolumes", mock.Anything).Return([]*dockerclient.Volume{}, nil)
	client.On("ListNetworks", mock.Anything).Return([]*dockerclient.NetworkResource{}, nil)

	// create a cluster from mock client
	c := simpleCluster(client)

	assert.Equal(t, c.TotalCpus(), int64(10))
	assert.Equal(t, c.TotalMemory(), int64(20))
}

func TestCreateContainer(t *testing.T) {

	// create mock client
	config := new(cluster.ContainerConfig)
	config.Labels = make(map[string]string)
	config.Labels["upm.ip"] = "192.168.11.124/24:enp0s25"

	dockinfo := new(dockerclient.ContainerInfo)
	dockinfo.Id = "123456789"
	dockinfo.Config = &dockerclient.ContainerConfig{
		Labels: map[string]string{
			"upm.ip": "192.168.11.124/24:enp0s25",
		},
	}

	client := mockclient.NewMockClient()
	client.On("Info").Return(mockInfo, nil)
	client.On("Version").Return(mockVersion, nil)
	client.On("StartMonitorEvents", mock.Anything, mock.Anything, mock.Anything).Return()
	client.On("ListContainers", true, false, "").Return([]dockerclient.Container{}, nil).Once()
	client.On("ListImages", mock.Anything).Return([]*dockerclient.Image{}, nil)
	client.On("ListVolumes", mock.Anything).Return([]*dockerclient.Volume{}, nil)
	client.On("ListNetworks", mock.Anything).Return([]*dockerclient.NetworkResource{}, nil)
	client.On("ListVolumes", mock.Anything).Return([]*dockerclient.Volume{}, nil)
	client.On("CreateContainer", mock.Anything, mock.Anything).Return("123456789", nil)
	client.On("RemoveContainer", mock.AnythingOfType("string"), true, true).Return(nil)
	client.On("InspectContainer", "123456789").Return(dockinfo, nil)

	// create a cluster from mock client
	c := simpleCluster(client)
	ca, err := c.CreateContainer(config, "lee.test")
	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, ca.Id, "123456789")

	err = c.RemoveContainer(ca, true, true)
	assert.Nil(t, err)
}
