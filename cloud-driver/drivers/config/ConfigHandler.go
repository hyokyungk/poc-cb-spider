package config

import (
	cblog "github.com/cloud-barista/cb-log"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("CB-SPIDER")
}

type Config struct {
	Openstack struct {
		DomainName       string `yaml:"domain_name"`
		IdentityEndpoint string `yaml:"identity_endpoint"`
		Password         string `yaml:"password"`
		ProjectID        string `yaml:"project_id"`
		Username         string `yaml:"username"`
		Region           string `yaml:"region"`
		VMName           string `yaml:"vm_name"`
		ImageId          string `yaml:"image_id"`
		FlavorId         string `yaml:"flavor_id"`
		NetworkId        string `yaml:"network_id"`
		SecurityGroups   string `yaml:"security_groups"`
		KeypairName      string `yaml:"keypair_name"`

		ServerId string `yaml:"server_id"`
	} `yaml:"openstack"`
	Azure struct {
		SubscriptionID string `yaml:"subscription_id"`
		GroupName      string `yaml:"group_name"`
		VMName         string `yaml:"vm_name"`

		Location string `yaml:"location"`
		VMSize   string `yaml:"vm_size"`
		Image    struct {
			Publisher string `yaml:"publisher"`
			Offer     string `yaml:"offer"`
			Sku       string `yaml:"sku"`
			Version   string `yaml:"version"`
		} `yaml:"image"`
		Os struct {
			ComputeName   string `yaml:"compute_name"`
			AdminUsername string `yaml:"admin_username"`
			AdminPassword string `yaml:"admin_password"`
		} `yaml:"os"`
		Network struct {
			ID      string `yaml:"id"`
			Primary bool   `yaml:"primary"`
		} `yaml:"network"`

		ServerId string `yaml:"server_id"`
	} `yaml:"azure"`
}

func ReadConfigFile() Config {
	// Set Environment Value of Project Root Path
	rootPath := os.Getenv("CBSPIDER_PATH")
	data, err := ioutil.ReadFile(rootPath + "/config/config.yaml")
	if err != nil {
		cblogger.Error(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		cblogger.Error(err)
	}

	return config
}

func GetServiceClient() (*gophercloud.ServiceClient, error) {

	// read configuration YAML file
	config := ReadConfigFile()

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: config.Openstack.IdentityEndpoint,
		Username:         config.Openstack.Username,
		Password:         config.Openstack.Password,
		DomainName:       config.Openstack.DomainName,
		TenantID:         config.Openstack.ProjectID,
		/*Scope: &gophercloud.AuthScope{
			ProjectID: config.Openstack.ProjectID,
		},*/
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.Openstack.Region,
	})
	if err != nil {
		return nil, err
	}

	return client, err
}

/*func getVMClient() (compute.VirtualMachinesClient, autorest.Authorizer) {
	config := ReadConfigFile()

	vmClient := compute.NewVirtualMachinesClient(config.Azure.SubscriptionID)
	iam, _ := iam.GetResourceManagementAuthorizer()
	return vmClient, iam
}*/

/*func GetAuthorizer() (autorest.Authorizer, error) {
	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	return authorizer, err
}*/
