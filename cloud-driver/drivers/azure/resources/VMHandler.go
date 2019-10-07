// Proof of Concepts of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is a Cloud Driver Example for PoC Test.
//
// by hyokyung.kim@innogrid.co.kr, 2019.07.

package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/Azure/go-autorest/autorest/to"
	cblog "github.com/cloud-barista/cb-log"
	idrv "github.com/cloud-barista/poc-cb-spider/cloud-driver/interfaces"
	irs "github.com/cloud-barista/poc-cb-spider/cloud-driver/interfaces/resources"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("CB-SPIDER")
}

type AzureVMHandler struct {
	Region idrv.RegionInfo
	Ctx    context.Context
	Client *compute.VirtualMachinesClient
}

func (vmHandler *AzureVMHandler) StartVM(vmReqInfo irs.VMReqInfo) (irs.VMInfo, error) {
	// Set VM Create Information
	imageId := vmReqInfo.ImageInfo.Id
	imageIdArr := strings.Split(imageId, ":")

	//sshKeyPath := "test"
	//sshKeyData := ""

	rootPath := os.Getenv("CBSPIDER_PATH")
	sshPublicKeyPath := rootPath + "/key/mcb-test-key.pub"
	var sshKeyData string
	if _, err := os.Stat(sshPublicKeyPath); err == nil {
		sshBytes, err := ioutil.ReadFile(sshPublicKeyPath)
		if err != nil {
			//log.Fatalf("failed to read SSH key data: %v", err)
			return irs.VMInfo{}, err
		}
		sshKeyData = string(sshBytes)
	}

	vmName := vmReqInfo.Name
	vmNameArr := strings.Split(vmName, ":")

	// Check VM Exists
	vm, err := vmHandler.Client.Get(vmHandler.Ctx, vmNameArr[0], vmNameArr[1], compute.InstanceView)
	if vm.ID != nil {
		errMsg := fmt.Sprintf("VirtualMachine with name %s already exist", vmNameArr[1])
		createErr := errors.New(errMsg)
		return irs.VMInfo{}, createErr
	}

	vmOpts := compute.VirtualMachine{
		Location: &vmHandler.Region.Region,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypes(vmReqInfo.SpecID),
			},
			StorageProfile: &compute.StorageProfile{
				ImageReference: &compute.ImageReference{
					Publisher: &imageIdArr[0],
					Offer:     &imageIdArr[1],
					Sku:       &imageIdArr[2],
					Version:   &imageIdArr[3],
				},
			},
			OsProfile: &compute.OSProfile{
				ComputerName:  &vmNameArr[1],
				AdminUsername: &vmReqInfo.LoginInfo.AdminUsername,
				//AdminPassword: &vmReqInfo.LoginInfo.AdminPassword,
				LinuxConfiguration: &compute.LinuxConfiguration{
					SSH: &compute.SSHConfiguration{
						PublicKeys: &[]compute.SSHPublicKey{
							{
								Path:    to.StringPtr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", vmReqInfo.LoginInfo.AdminUsername)),
								KeyData: &sshKeyData,
							},
						},
					},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: &vmReqInfo.VNetworkInfo.Id,
						NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
							Primary: to.BoolPtr(true),
						},
					},
				},
			},
		},
	}

	future, err := vmHandler.Client.CreateOrUpdate(vmHandler.Ctx, vmNameArr[0], vmNameArr[1], vmOpts)
	if err != nil {
		cblogger.Error(err)
		return irs.VMInfo{}, err
	}
	err = future.WaitForCompletionRef(vmHandler.Ctx, vmHandler.Client.Client)
	if err != nil {
		cblogger.Error(err)
		return irs.VMInfo{}, err
	}

	vm, err = vmHandler.Client.Get(vmHandler.Ctx, vmNameArr[0], vmNameArr[1], compute.InstanceView)
	if err != nil {
		cblogger.Error(err)
	}
	vmInfo := mappingServerInfo(vm)

	return vmInfo, nil
}

func (vmHandler *AzureVMHandler) SuspendVM(vmID string) {
	vmIdArr := strings.Split(vmID, ":")

	future, err := vmHandler.Client.PowerOff(vmHandler.Ctx, vmIdArr[0], vmIdArr[1])
	if err != nil {
		cblogger.Error(err)
	}
	err = future.WaitForCompletionRef(vmHandler.Ctx, vmHandler.Client.Client)
	if err != nil {
		cblogger.Error(err)
	}
}

func (vmHandler *AzureVMHandler) ResumeVM(vmID string) {
	vmIdArr := strings.Split(vmID, ":")

	future, err := vmHandler.Client.Start(vmHandler.Ctx, vmIdArr[0], vmIdArr[1])
	if err != nil {
		cblogger.Error(err)
	}
	err = future.WaitForCompletionRef(vmHandler.Ctx, vmHandler.Client.Client)
	if err != nil {
		cblogger.Error(err)
	}
}

func (vmHandler *AzureVMHandler) RebootVM(vmID string) {
	vmIdArr := strings.Split(vmID, ":")

	future, err := vmHandler.Client.Restart(vmHandler.Ctx, vmIdArr[0], vmIdArr[1])
	if err != nil {
		cblogger.Error(err)
	}
	err = future.WaitForCompletionRef(vmHandler.Ctx, vmHandler.Client.Client)
	if err != nil {
		cblogger.Error(err)
	}
}

func (vmHandler *AzureVMHandler) TerminateVM(vmID string) {
	vmIdArr := strings.Split(vmID, ":")

	future, err := vmHandler.Client.Delete(vmHandler.Ctx, vmIdArr[0], vmIdArr[1])
	//future, err := vmHandler.Client.Deallocate(vmHandler.Ctx, vmIdArr[0], vmIdArr[1])
	if err != nil {
		cblogger.Error(err)
	}
	err = future.WaitForCompletionRef(vmHandler.Ctx, vmHandler.Client.Client)
	if err != nil {
		cblogger.Error(err)
	}
}

func (vmHandler *AzureVMHandler) ListVMStatus() []*irs.VMStatusInfo {
	//serverList, err := vmHandler.Client.ListAll(vmHandler.Ctx)
	serverList, err := vmHandler.Client.List(vmHandler.Ctx, vmHandler.Region.ResourceGroup)
	if err != nil {
		cblogger.Error(err)
	}

	var vmStatusList []*irs.VMStatusInfo
	for _, s := range serverList.Values() {
		if s.InstanceView != nil {
			statusStr := getVmStatus(*s.InstanceView)
			status := irs.VMStatus(statusStr)
			vmStatusInfo := irs.VMStatusInfo{
				VmId:     *s.ID,
				VmStatus: status,
			}
			vmStatusList = append(vmStatusList, &vmStatusInfo)
		} else {
			vmIdArr := strings.Split(*s.ID, "/")
			vmId := vmIdArr[4] + ":" + vmIdArr[8]
			status := vmHandler.GetVMStatus(vmId)
			vmStatusInfo := irs.VMStatusInfo{
				VmId:     *s.ID,
				VmStatus: status,
			}
			vmStatusList = append(vmStatusList, &vmStatusInfo)
		}
	}

	return vmStatusList
}

func (vmHandler *AzureVMHandler) GetVMStatus(vmID string) irs.VMStatus {
	vmIdArr := strings.Split(vmID, ":")
	instanceView, err := vmHandler.Client.InstanceView(vmHandler.Ctx, vmIdArr[0], vmIdArr[1])
	if err != nil {
		cblogger.Error(err)
	}

	// Get powerState, provisioningState
	vmStatus := getVmStatus(instanceView)
	return irs.VMStatus(vmStatus)
}

func (vmHandler *AzureVMHandler) ListVM() []*irs.VMInfo {
	//serverList, err := vmHandler.Client.ListAll(vmHandler.Ctx)
	serverList, err := vmHandler.Client.List(vmHandler.Ctx, vmHandler.Region.ResourceGroup)
	if err != nil {
		cblogger.Error(err)
	}

	var vmList []*irs.VMInfo
	for _, server := range serverList.Values() {
		vmInfo := mappingServerInfo(server)
		vmList = append(vmList, &vmInfo)
	}

	return vmList
}

func (vmHandler *AzureVMHandler) GetVM(vmID string) irs.VMInfo {
	vmIdArr := strings.Split(vmID, ":")
	vm, err := vmHandler.Client.Get(vmHandler.Ctx, vmIdArr[0], vmIdArr[1], compute.InstanceView)
	if err != nil {
		cblogger.Error(err)
	}

	vmInfo := mappingServerInfo(vm)
	return vmInfo
}

func getVmStatus(instanceView compute.VirtualMachineInstanceView) string {
	var powerState, provisioningState string

	for _, stat := range *instanceView.Statuses {
		statArr := strings.Split(*stat.Code, "/")

		if statArr[0] == "PowerState" {
			powerState = statArr[1]
		} else if statArr[0] == "ProvisioningState" {
			provisioningState = statArr[1]
		}
	}

	// Set VM Status Info
	var vmState string
	if powerState != "" && provisioningState != "" {
		vmState = powerState + "(" + provisioningState + ")"
	} else if powerState != "" && provisioningState == "" {
		vmState = powerState
	} else if powerState == "" && provisioningState != "" {
		vmState = provisioningState
	} else {
		vmState = "-"
	}
	return vmState
}

func mappingServerInfo(server compute.VirtualMachine) irs.VMInfo {

	// Get Default VM Info
	vmInfo := irs.VMInfo{
		Name: *server.Name,
		Id:   *server.ID,
		Region: irs.RegionInfo{
			Region: *server.Location,
		},
		SpecID: string(server.VirtualMachineProperties.HardwareProfile.VMSize),
	}

	// Set VM Zone
	if server.Zones != nil {
		vmInfo.Region.Zone = (*server.Zones)[0]
	}

	// Set VM Image Info
	imageRef := server.VirtualMachineProperties.StorageProfile.ImageReference
	imageId := *imageRef.Publisher + ":" + *imageRef.Offer + ":" + *imageRef.Sku + ":" + *imageRef.Version
	vmInfo.ImageID = imageId

	// Set VNic Info
	niList := *server.NetworkProfile.NetworkInterfaces
	for _, ni := range niList {
		if ni.NetworkInterfaceReferenceProperties != nil {
			vmInfo.VNIC = *ni.ID
		}
	}

	// Set GuestUser Id/Pwd
	if server.VirtualMachineProperties.OsProfile.AdminUsername != nil {
		vmInfo.GuestUserID = *server.VirtualMachineProperties.OsProfile.AdminUsername
	}
	if server.VirtualMachineProperties.OsProfile.AdminPassword != nil {
		vmInfo.GuestUserID = *server.VirtualMachineProperties.OsProfile.AdminPassword
	}

	// Set BootDisk
	if server.VirtualMachineProperties.StorageProfile.OsDisk.Name != nil {
		vmInfo.GuestBootDisk = *server.VirtualMachineProperties.StorageProfile.OsDisk.Name
	}

	return vmInfo
}
