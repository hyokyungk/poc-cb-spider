// Cloud Driver Interface of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is Resouces interfaces of Cloud Driver.
//
// by powerkim@etri.re.kr, 2019.06.

package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	idrv "github.com/cloud-barista/poc-cb-spider/cloud-driver/interfaces"
	irs "github.com/cloud-barista/poc-cb-spider/cloud-driver/interfaces/resources"
	"github.com/davecgh/go-spew/spew"
)

type AwsPublicIPHandler struct {
	Region idrv.RegionInfo
	Client *ec2.EC2
}

//@TODO : EC2에 Public를 할당하는 Associate함수 필요 함.

func (publicIpHandler *AwsPublicIPHandler) CreatePublicIP(publicIPReqInfo irs.PublicIPReqInfo) (irs.PublicIPInfo, error) {
	cblogger.Info("Start : ", publicIPReqInfo)
	//result, err := vNetworkHandler.Client.DescribeKeyPairs(input)

	//@TODO: 대체해야 함.
	instanceID := publicIPReqInfo.Id

	// Attempt to allocate the Elastic IP address.
	allocRes, err := publicIpHandler.Client.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"), // 범이 : VPC
	})

	if err != nil {
		cblogger.Errorf("Unable to allocate IP address, %v", err)
		return irs.PublicIPInfo{}, err
	}

	spew.Dump(allocRes)
	cblogger.Infof("EIP 생성 성공 - Public IP : [%s], Allocation Id : [%s]", *allocRes.PublicIp, *allocRes.AllocationId)

	cblogger.Infof("[%s] EC2에 [%s] IP 할당 시작", instanceID, *allocRes.PublicIp)
	// EC2에 할당.
	// Associate the new Elastic IP address with an existing EC2 instance.
	assocRes, err := publicIpHandler.Client.AssociateAddress(&ec2.AssociateAddressInput{
		AllocationId: allocRes.AllocationId,
		InstanceId:   aws.String(instanceID),
	})
	if err != nil {
		cblogger.Errorf("Unable to associate IP address with %s, %v", instanceID, err)
		return irs.PublicIPInfo{}, err
	}
	cblogger.Infof("[%s] EC2에 [%s] IP 할당 완료 - Allocation Id : [%s]", instanceID, *allocRes.PublicIp, *assocRes.AssociationId)

	return irs.PublicIPInfo{}, nil
}

func (publicIpHandler *AwsPublicIPHandler) ListPublicIP() ([]*irs.PublicIPInfo, error) {
	return nil, nil
}

func (publicIpHandler *AwsPublicIPHandler) GetPublicIP(publicIPID string) (irs.PublicIPInfo, error) {
	return irs.PublicIPInfo{}, nil
}

func (publicIpHandler *AwsPublicIPHandler) DeletePublicIP(publicIPID string) (bool, error) {
	return false, nil
}
