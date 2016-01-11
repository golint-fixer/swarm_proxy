package lonely

import (
	"crypto/tls"
	//"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/discovery"
	"github.com/docker/swarm/scheduler"
	"github.com/docker/swarm/scheduler/node"
	"github.com/samalba/dockerclient"
)

var virueEngine = cluster.NewEngine("unix:///var/run/docker.sock", 0)

type pendingContainer struct {
	Config *cluster.ContainerConfig
	Name   string
	Engine *cluster.Engine
}

func (p *pendingContainer) ToContainer() *cluster.Container {

	return nil
}

// Cluster is exported
type Cluster struct {
	sync.RWMutex

	engine            dockerclient.Client
	pendingContainers map[string]*pendingContainer

	TLSConfig *tls.Config
}

// NewCluster is exported
func NewCluster(_ *scheduler.Scheduler, TLSConfig *tls.Config, _ discovery.Discovery, options cluster.DriverOpts) (cluster.Cluster, error) {
	log.WithFields(log.Fields{"name": "lonely"}).Debug("Initializing cluster")
	client, err := dockerclient.NewDockerClientTimeout("unix:///var/run/docker.sock", TLSConfig, time.Second*5)
	if err != nil {
		return nil, err
	}
	return simpleCluster(client), nil
}

func simpleCluster(client dockerclient.Client) cluster.Cluster {
	return &Cluster{
		engine:            client,
		TLSConfig:         nil,
		pendingContainers: make(map[string]*pendingContainer),
	}
}

// Handle callbacks for the events
func (c *Cluster) Handle(_ *cluster.Event) error {
	return nil
}

// RegisterEventHandler registers an event handler.
func (c *Cluster) RegisterEventHandler(_ cluster.EventHandler) error {
	return nil
}

// Generate a globally (across the cluster) unique ID.
func (c *Cluster) generateUniqueID() string {
	return stringid.GenerateRandomID()
}

// CreateContainer aka schedule a brand new container into the cluster.
func (c *Cluster) CreateContainer(config *cluster.ContainerConfig, name string) (*cluster.Container, error) {
	container, err := c.createContainer(config, name, false)

	//  fails with image not found, then try to reschedule with soft-image-affinity
	if err != nil {
		bImageNotFoundError, _ := regexp.MatchString(`image \S* not found`, err.Error())
		if bImageNotFoundError && !config.HaveNodeConstraint() {
			// Check if the image exists in the cluster
			// If exists, retry with a soft-image-affinity
			if image := c.Image(config.Image); image != nil {
				container, err = c.createContainer(config, name, true)
			}
		}
	}
	return container, err
}

func (c *Cluster) createContainer(config *cluster.ContainerConfig, name string, withSoftImageAffinity bool) (*cluster.Container, error) {
	//append ip process
	dockerConfig := &config.ContainerConfig
	id, err := c.engine.CreateContainer(dockerConfig, name)
	if err != nil {
		return nil, err
	}
	clusterC := &cluster.Container{}
	clusterC.Id = id
	return clusterC, nil
}

// RemoveContainer aka Remove a container from the cluster.
func (c *Cluster) RemoveContainer(container *cluster.Container, force, volumes bool) error {
	err := c.engine.RemoveContainer(container.Id, force, volumes)
	if err != nil {
		return err
	}
	return nil
}

// RemoveNetwork removes a network from the cluster
func (c *Cluster) RemoveNetwork(network *cluster.Network) error {
	err := c.engine.RemoveNetwork(network.ID)
	return err
}

// Images returns all the images in the cluster.
func (c *Cluster) Images() cluster.Images {
	cImages := []*cluster.Image{}

	images, err := c.engine.ListImages(true)
	if err != nil {
		return cluster.Images(cImages)
	}

	for _, image := range images {
		cImages = append(cImages, &cluster.Image{Image: *image, Engine: nil})
	}

	return cluster.Images(cImages)
}

// Image returns an image with IDOrName in the cluster
func (c *Cluster) Image(IDOrName string) *cluster.Image {
	// Abort immediately if the name is empty.
	//if len(IDOrName) == 0 {
	//	return nil
	//}
	//c.engine.
	//if image, err := c.engine.InspectImage(IDOrName); image != nil && err != nil {
	//	return &cluster.Image{Image: *image, Engine: nil}
	//}

	return nil
}

// RemoveImages removes all the images that match `name` from the cluster
func (c *Cluster) RemoveImages(name string, force bool) ([]*dockerclient.ImageDelete, error) {

	return nil, nil
}

func (c *Cluster) refreshNetworks() {

}

// CreateNetwork creates a network in the cluster
func (c *Cluster) CreateNetwork(request *dockerclient.NetworkCreate) (response *dockerclient.NetworkCreateResponse, err error) {
	var (
		parts = strings.SplitN(request.Name, "/", 2)
	)

	if len(parts) == 2 {
		// a node was specified, create the container only on this node
		request.Name = parts[1]

	}

	resp, err := c.engine.CreateNetwork(request)
	return resp, err

}

// CreateVolume creates a volume in the cluster
func (c *Cluster) CreateVolume(request *dockerclient.VolumeCreateRequest) (*cluster.Volume, error) {
	return nil, nil
}

// RemoveVolumes removes all the volumes that match `name` from the cluster
func (c *Cluster) RemoveVolumes(name string) (bool, error) {
	return false, nil
}

// Pull is exported
func (c *Cluster) Pull(name string, authConfig *dockerclient.AuthConfig, callback func(where, status string, err error)) {

}

// Load image
func (c *Cluster) Load(imageReader io.Reader, callback func(where, status string, err error)) {

}

// Import image
func (c *Cluster) Import(source string, repository string, tag string, imageReader io.Reader, callback func(what, status string, err error)) {

}

// Containers returns all the containers in the cluster.
func (c *Cluster) Containers() cluster.Containers {
	return nil
}

func (c *Cluster) checkNameUniqueness(name string) bool {

	return true
}

// Container returns the container with IDOrName in the cluster
func (c *Cluster) Container(IDOrName string) *cluster.Container {
	cl := &cluster.Container{}
	cl.Container.Id = IDOrName
	cl.Engine = virueEngine
	return cl
}

// Networks returns all the networks in the cluster.
func (c *Cluster) Networks() cluster.Networks {

	return nil
}

// Volumes returns all the volumes in the cluster.
func (c *Cluster) Volumes() []*cluster.Volume {

	return nil
}

// Volume returns the volume name in the cluster
func (c *Cluster) Volume(name string) *cluster.Volume {
	return nil
}

// listNodes returns all the engines in the cluster.
func (c *Cluster) listNodes() []*node.Node {
	return nil
}

// listEngines returns all the engines in the cluster.
func (c *Cluster) listEngines() []*cluster.Engine {
	return nil
}

// TotalMemory return the total memory of the cluster
func (c *Cluster) TotalMemory() int64 {
	info, err := c.engine.Info()
	if err != nil {
		return 0
	}
	return info.MemTotal
}

// TotalCpus return the total memory of the cluster
func (c *Cluster) TotalCpus() int64 {
	info, err := c.engine.Info()
	if err != nil {
		return 0
	}
	return info.NCPU
}

// Info returns some info about the cluster, like nb or containers / images
func (c *Cluster) Info() [][]string {
	return nil
}

// RANDOMENGINE returns a random engine.
func (c *Cluster) RANDOMENGINE() (*cluster.Engine, error) {
	return nil, nil
}

// RenameContainer rename a container
func (c *Cluster) RenameContainer(container *cluster.Container, newName string) error {
	return nil
}

// BuildImage build an image
func (c *Cluster) BuildImage(buildImage *dockerclient.BuildImage, out io.Writer) error {
	return nil
}

// TagImage tag an image
func (c *Cluster) TagImage(IDOrName string, repo string, tag string, force bool) error {
	return nil
}
