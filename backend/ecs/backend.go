package backend_ecs

import (
	"strings"
	"time"
	"strconv"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ecs"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/daryl/cash"
	"gopkg.in/alecthomas/kingpin.v1"
	log "github.com/Sirupsen/logrus"

	"github.com/nickschuch/marco/backend"
	"github.com/nickschuch/marco/handling"
)

var (
	cliECSRegion  = kingpin.Flag("ecs-region", "The region to run the build in.").Default("ap-southeast-2").OverrideDefaultFromEnvar("ECS_REGION").String()
	cliECSCluster = kingpin.Flag("ecs-cluster", "The cluster to run this build against.").Default("default").OverrideDefaultFromEnvar("ECS_CLUSTER").String()
	cliECSPorts   = kingpin.Flag("ecs-ports", "The ports you wish to proxy.").Default("80,8080,2368,8983").String()
)

type BackendECS struct {
	// This is where we store the container addresses
	// for a period of time. It cuts down on API calls
	// and frees up resource.
	cache *cash.Cash
}

func init() {
	backend.Register("ecs", &BackendECS{})
}

func (b *BackendECS) Start() error {
	// Create a brand new cache item that we can use to
	// store all the domains and there associated list of
	// addresses.
	b.cache = cash.New(cash.Conf{
		// Default expiration.
		time.Minute,
		// Clean interval.
		30 * time.Minute,
	})

	return nil
}

func (b *BackendECS) Addresses(domain string) ([]string, error) {
	var list []string

	if v, ok := b.cache.Get(domain); ok {
		// We found a cached item and we should return it's list of urls.
		list = v.([]string)
		return list, nil
	}

	// This call was already at a cost and given we couldn't filter on the domain for results.
	// We might as well set all the associated URLs as well as the one is missing.
	list, err := getListByDomain(domain)
	handling.Check(err)
	b.cache.Set(domain, list, time.Minute)

	return list, nil
}

func getListByDomain(domain string) ([]string, error) {
	var domainList []string

	list, err := getList()
	handling.Check(err)
	if len(list[domain]) > 0 {
		domainList = list[domain]
	}

	return domainList, nil
}

func getList() (map[string][]string, error) {
	client := getECSClient()

	list := make(map[string][]string)
	ips := make(map[string]string)

	tasksInput := &ecs.ListTasksInput{}
	tasks, err := client.ListTasks(tasksInput)
	check(err)

	// We only have one task that we care about. So we are
	// only going to pass this one in the list.
	describeInput := &ecs.DescribeTasksInput{
		Cluster: aws.String(*cliECSCluster),
		Tasks:   tasks.TaskARNs,
	}
	described, err := client.DescribeTasks(describeInput)
	check(err)

	// Get the IP address of each of the container instances.
	// That way we can use these addresses further down on our container urls.
	instancesInput := &ecs.ListContainerInstancesInput{
		Cluster:    aws.String(*cliECSCluster),
	}
	instances, err := client.ListContainerInstances(instancesInput)
	check(err)
	for _, i := range instances.ContainerInstanceARNs {
		containerInstance := getContainerInstance(i)
		ips[*i] = getEc2IP(containerInstance.EC2InstanceID) 
	}

	// Loop over the containers and build a list of urls to hit.
	for _, t := range described.Tasks {
		for _, c := range t.Containers {
			// Loop over all the ports that have been exposed.
			for _, p := range c.NetworkBindings {
				// Check that this container has exposed the port that we require.
				containerPort := strconv.FormatInt(*p.ContainerPort, 10)
				if ! strings.Contains(*cliECSPorts, containerPort) {
					continue
				}

				// Add the port to the list.
				hostIP := ips[*t.ContainerInstanceARN]
				hostPort := strconv.FormatInt(*p.HostPort, 10)
				url := "http://"+hostIP+":"+hostPort
				list[*c.Name] = append(list[*c.Name], url)
			}
		}
	}

	return list, nil
}

func getECSClient() *ecs.ECS {
	client := ecs.New(&aws.Config{Region: *cliECSRegion})
	return client
}

func getEC2Client() *ec2.EC2 {
	client := ec2.New(&aws.Config{Region: *cliECSRegion})
	return client
}

// Helper function to check if something went wrong with
// code run against AWS ECS.
func check(err error) {
	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		log.Info(awserr.Message)
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
	}
}

func getContainerInstance(arn *string) *ecs.ContainerInstance {
	client := getECSClient()

	params := &ecs.DescribeContainerInstancesInput{
		ContainerInstances: []*string{
			aws.String(*arn),
		},
		Cluster: aws.String(*cliECSCluster),
	}
	resp, err := client.DescribeContainerInstances(params)
	check(err)

	return resp.ContainerInstances[0];
}

func getEc2IP(id *string) string {
	client := getEC2Client()

	// Query the EC2 backend for the host that we require.
	params := &ec2.DescribeInstancesInput{
		InstanceIDs: []*string{
			aws.String(*id),
		},
	}
	resp, err := client.DescribeInstances(params)
	check(err)

	// https://github.com/awslabs/aws-sdk-go/blob/master/service/ec2/api.go#L13194
	return *resp.Reservations[0].Instances[0].PublicIPAddress;
}
