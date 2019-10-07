// Proof of Concepts of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is a Cloud Driver Example for PoC Test.
//
// by hyokyung.kim@innogrid.co.kr, 2019.07.

package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	azcon "github.com/cloud-barista/poc-cb-spider/cloud-driver/drivers/azure/connect"
	idrv "github.com/cloud-barista/poc-cb-spider/cloud-driver/interfaces"
	icon "github.com/cloud-barista/poc-cb-spider/cloud-driver/interfaces/connect"
	"time"
)

/*var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("CB-SPIDER")
}*/

type AzureDriver struct{}

func (AzureDriver) GetDriverVersion() string {
	return "AZURE DRIVER Version 1.0"
}

func (AzureDriver) GetDriverCapability() idrv.DriverCapabilityInfo {
	var drvCapabilityInfo idrv.DriverCapabilityInfo

	drvCapabilityInfo.ImageHandler = false
	drvCapabilityInfo.VNetworkHandler = false
	drvCapabilityInfo.SecurityHandler = false
	drvCapabilityInfo.KeyPairHandler = false
	drvCapabilityInfo.VNicHandler = false
	drvCapabilityInfo.PublicIPHandler = false
	drvCapabilityInfo.VMHandler = true

	return drvCapabilityInfo
}

func (driver *AzureDriver) ConnectCloud(connectionInfo idrv.ConnectionInfo) (icon.CloudConnection, error) {
	// 1. get info of credential and region for Test A Cloud from connectionInfo.
	// 2. create a client object(or service  object) of Test A Cloud with credential info.
	// 3. create CloudConnection Instance of "connect/TDA_CloudConnection".
	// 4. return CloudConnection Interface of TDA_CloudConnection.

	Ctx, VMClient, err := getVMClient(connectionInfo.CredentialInfo)
	if err != nil {
		return nil, err
	}
	Ctx, imageClient, err := getImageClient(connectionInfo.CredentialInfo)
	if err != nil {
		return nil, err
	}
	Ctx, publicIPClient, err := getPublicIPClient(connectionInfo.CredentialInfo)
	if err != nil {
		return nil, err
	}
	Ctx, sgClient, err := getSecurityGroupClient(connectionInfo.CredentialInfo)
	if err != nil {
		return nil, err
	}
	Ctx, vNicClient, err := getVNicClient(connectionInfo.CredentialInfo)
	if err != nil {
		return nil, err
	}
	Ctx, SubnetClient, err := getSubnetClient(connectionInfo.CredentialInfo)
	if err != nil {
		return nil, err
	}
	Ctx, VNetClient, err := getVNetworkClient(connectionInfo.CredentialInfo)
	if err != nil {
		return nil, err
	}
	iConn := azcon.AzureCloudConnection{
		Region:              connectionInfo.RegionInfo,
		Ctx:                 Ctx,
		VMClient:            VMClient,
		ImageClient:         imageClient,
		PublicIPClient:      publicIPClient,
		SecurityGroupClient: sgClient,
		VNetClient:          VNetClient,
		VNicClient:          vNicClient,
		SubnetClient:        SubnetClient,
	}
	return &iConn, nil
}

func getVMClient(credential idrv.CredentialInfo) (context.Context, *compute.VirtualMachinesClient, error) {
	/*auth.NewClientCredentialsConfig()
	  authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	  if err != nil {
	      return nil, nil, err
	  }*/
	config := auth.NewClientCredentialsConfig(credential.ClientId, credential.ClientSecret, credential.TenantId)
	authorizer, err := config.Authorizer()
	if err != nil {
		return nil, nil, err
	}

	vmClient := compute.NewVirtualMachinesClient(credential.SubscriptionId)
	vmClient.Authorizer = authorizer
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	return ctx, &vmClient, nil
}

func getImageClient(credential idrv.CredentialInfo) (context.Context, *compute.ImagesClient, error) {
	config := auth.NewClientCredentialsConfig(credential.ClientId, credential.ClientSecret, credential.TenantId)
	authorizer, err := config.Authorizer()
	if err != nil {
		return nil, nil, err
	}

	imageClient := compute.NewImagesClient(credential.SubscriptionId)
	imageClient.Authorizer = authorizer
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	return ctx, &imageClient, nil
}

func getPublicIPClient(credential idrv.CredentialInfo) (context.Context, *network.PublicIPAddressesClient, error) {
	config := auth.NewClientCredentialsConfig(credential.ClientId, credential.ClientSecret, credential.TenantId)
	authorizer, err := config.Authorizer()
	if err != nil {
		return nil, nil, err
	}

	publicIPClient := network.NewPublicIPAddressesClient(credential.SubscriptionId)
	publicIPClient.Authorizer = authorizer
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	return ctx, &publicIPClient, nil
}

func getSecurityGroupClient(credential idrv.CredentialInfo) (context.Context, *network.SecurityGroupsClient, error) {
	config := auth.NewClientCredentialsConfig(credential.ClientId, credential.ClientSecret, credential.TenantId)
	authorizer, err := config.Authorizer()
	if err != nil {
		return nil, nil, err
	}

	sgClient := network.NewSecurityGroupsClient(credential.SubscriptionId)
	sgClient.Authorizer = authorizer
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	return ctx, &sgClient, nil
}

func getVNetworkClient(credential idrv.CredentialInfo) (context.Context, *network.VirtualNetworksClient, error) {
	config := auth.NewClientCredentialsConfig(credential.ClientId, credential.ClientSecret, credential.TenantId)
	authorizer, err := config.Authorizer()
	if err != nil {
		return nil, nil, err
	}

	vNetClient := network.NewVirtualNetworksClient(credential.SubscriptionId)
	vNetClient.Authorizer = authorizer
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	return ctx, &vNetClient, nil
}

func getVNicClient(credential idrv.CredentialInfo) (context.Context, *network.InterfacesClient, error) {
	config := auth.NewClientCredentialsConfig(credential.ClientId, credential.ClientSecret, credential.TenantId)
	authorizer, err := config.Authorizer()
	if err != nil {
		return nil, nil, err
	}

	vNicClient := network.NewInterfacesClient(credential.SubscriptionId)
	vNicClient.Authorizer = authorizer
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	return ctx, &vNicClient, nil
}

func getSubnetClient(credential idrv.CredentialInfo) (context.Context, *network.SubnetsClient, error) {
	config := auth.NewClientCredentialsConfig(credential.ClientId, credential.ClientSecret, credential.TenantId)
	authorizer, err := config.Authorizer()
	if err != nil {
		return nil, nil, err
	}

	subnetClient := network.NewSubnetsClient(credential.SubscriptionId)
	subnetClient.Authorizer = authorizer
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	return ctx, &subnetClient, nil
}

var TestDriver AzureDriver
