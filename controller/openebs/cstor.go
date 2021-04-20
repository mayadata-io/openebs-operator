/*
Copyright 2020 The MayaData Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openebs

import (
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/k8s"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"strings"
)

const (
	// ContainerOpenEBSCSIPluginName is the name of the container openebs csi plugin
	ContainerOpenEBSCSIPluginName string = "openebs-csi-plugin"
	// ContainerCSTORCSIPluginName is the name of the container cstor csi plugin
	ContainerCSTORCSIPluginName string = "cstor-csi-plugin"
	// ContainerCSIResizerName is the name of csi-resizer container
	ContainerCSIResizerName string = "csi-resizer"
	// ContainerCSISnapshotterName is the name of csi-snapshotter container
	ContainerCSISnapshotterName string = "csi-snapshotter"
	// ContainerCSISnapshotControllerName is the name of snapshot-controller container
	ContainerCSISnapshotControllerName string = "snapshot-controller"
	// ContainerCSIProvisionerName is the name of csi-provisioner container
	ContainerCSIProvisionerName string = "csi-provisioner"
	// ContainerCSIAttacherName is the name of csi-attacher container
	ContainerCSIAttacherName string = "csi-attacher"
	// ContainerCSIClusterDriverRegistrarName is the name of csi-cluster-driver-registrar container
	ContainerCSIClusterDriverRegistrarName string = "csi-cluster-driver-registrar"
	// ContainerCSINodeDriverRegistrarName is the name of csi-node-driver-registrar container
	ContainerCSINodeDriverRegistrarName string = "csi-node-driver-registrar"
	// ContainerCSIDriverRegistrarName is the name of csi-driver-registrar container
	ContainerCSIDriverRegistrarName string = "csi-driver-registrar"
	// EnvOpenEBSNamespaceKey is the env key for openebs namespace
	EnvOpenEBSNamespaceKey string = "OPENEBS_NAMESPACE"
	// EnvDriverRegSocketPathKey is the env key for driver registration socket path.
	EnvDriverRegSocketPathKey string = "DRIVER_REG_SOCK_PATH"
	// DefaultCSPCOperatorReplicaCount is the default replica count for
	// cspc-operator.
	DefaultCSPCOperatorReplicaCount int32 = 1
	// DefaultCVCOperatorReplicaCount is the default replica count for
	// cvc-operator.
	DefaultCVCOperatorReplicaCount int32 = 1
	// DefaultCStorAdmissionServerReplicaCount is the default replica count for
	// AdmissionServer.
	DefaultCStorAdmissionServerReplicaCount int32 = 1
)

var (
	// List of images which are by default fetched from quay.io/k8scsi registry.
	CSIResizerImage                       string
	CSISnapshotterImage                   string
	CSISnapshotControllerImage            string
	CSIProvisionerForCSIControllerImage   string
	CSIAttacherForCSIControllerImage      string
	CSIClusterDriverRegistrarImage        string
	CSINodeDriverRegistrarForCSINodeImage string
	KubeletPath                           string
)

// SupportedCSIResizerVersionForOpenEBSVersion stores the mapping for
// CSI resizer to OpenEBS version.
var SupportedCSIResizerVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion190:    types.CSIResizerVersion010,
	types.OpenEBSVersion190EE:  types.CSIResizerVersion010,
	types.OpenEBSVersion1100:   types.CSIResizerVersion040,
	types.OpenEBSVersion1100EE: types.CSIResizerVersion040,
	types.OpenEBSVersion1110:   types.CSIResizerVersion040,
	types.OpenEBSVersion1110EE: types.CSIResizerVersion040,
	types.OpenEBSVersion1120:   types.CSIResizerVersion040,
	types.OpenEBSVersion1120EE: types.CSIResizerVersion040,
	types.OpenEBSVersion200:    types.CSIResizerVersion040,
	types.OpenEBSVersion200EE:  types.CSIResizerVersion040,
	types.OpenEBSVersion210:    types.CSIResizerVersion040,
	types.OpenEBSVersion210EE:  types.CSIResizerVersion040,
	types.OpenEBSVersion220:    types.CSIResizerVersion040,
	types.OpenEBSVersion220EE:  types.CSIResizerVersion040,
	types.OpenEBSVersion240:    types.CSIResizerVersion040,
	types.OpenEBSVersion250:    types.CSIResizerVersion110,
	types.OpenEBSVersion260:    types.CSIResizerVersion110,
	types.OpenEBSVersion270:    types.CSIResizerVersion110,
	types.OpenEBSVersion280:    types.CSIResizerVersion110,
}

// SupportedCSISnapshotterVersionForOpenEBSVersion stores the mapping for
// CSI snapshotter to OpenEBS version.
var SupportedCSISnapshotterVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion190:    types.CSISnapshotterVersion201,
	types.OpenEBSVersion190EE:  types.CSISnapshotterVersion201,
	types.OpenEBSVersion1100:   types.CSISnapshotterVersion201,
	types.OpenEBSVersion1100EE: types.CSISnapshotterVersion201,
	types.OpenEBSVersion1110:   types.CSISnapshotterVersion201,
	types.OpenEBSVersion1110EE: types.CSISnapshotterVersion201,
	types.OpenEBSVersion1120:   types.CSISnapshotterVersion201,
	types.OpenEBSVersion1120EE: types.CSISnapshotterVersion201,
	types.OpenEBSVersion200:    types.CSISnapshotterVersion201,
	types.OpenEBSVersion200EE:  types.CSISnapshotterVersion201,
	types.OpenEBSVersion210:    types.CSISnapshotterVersion201,
	types.OpenEBSVersion210EE:  types.CSISnapshotterVersion201,
	types.OpenEBSVersion220:    types.CSISnapshotterVersion201,
	types.OpenEBSVersion220EE:  types.CSISnapshotterVersion201,
	types.OpenEBSVersion240:    types.CSISnapshotterVersion201,
	types.OpenEBSVersion250:    types.CSISnapshotterVersion303,
	types.OpenEBSVersion260:    types.CSISnapshotterVersion303,
	types.OpenEBSVersion270:    types.CSISnapshotterVersion303,
	types.OpenEBSVersion280:    types.CSISnapshotterVersion303,
}

// SupportedCSISnapshotControllerVersionForOpenEBSVersion stores the mapping for
// CSI snapshot controller to OpenEBS version.
var SupportedCSISnapshotControllerVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion190:    types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion190EE:  types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion1100:   types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion1100EE: types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion1110:   types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion1110EE: types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion1120:   types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion1120EE: types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion200:    types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion200EE:  types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion210:    types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion210EE:  types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion220:    types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion220EE:  types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion240:    types.CSISnapshotControllerVersion201,
	types.OpenEBSVersion250:    types.CSISnapshotControllerVersion303,
	types.OpenEBSVersion260:    types.CSISnapshotControllerVersion303,
	types.OpenEBSVersion270:    types.CSISnapshotControllerVersion303,
	types.OpenEBSVersion280:    types.CSISnapshotControllerVersion303,
}

// SupportedCSIProvisionerVersionForCSIControllerVersion stores the mapping for
// CSI provisioner to csi-controller version.
var SupportedCSIProvisionerVersionForCSIControllerVersion = map[string]string{
	types.OpenEBSVersion190:    types.CSIProvisionerVersion150,
	types.OpenEBSVersion190EE:  types.CSIProvisionerVersion150,
	types.OpenEBSVersion1100:   types.CSIProvisionerVersion150,
	types.OpenEBSVersion1100EE: types.CSIProvisionerVersion150,
	types.OpenEBSVersion1110:   types.CSIProvisionerVersion160,
	types.OpenEBSVersion1110EE: types.CSIProvisionerVersion160,
	types.OpenEBSVersion1120:   types.CSIProvisionerVersion160,
	types.OpenEBSVersion1120EE: types.CSIProvisionerVersion160,
	types.OpenEBSVersion200:    types.CSIProvisionerVersion160,
	types.OpenEBSVersion200EE:  types.CSIProvisionerVersion160,
	types.OpenEBSVersion210:    types.CSIProvisionerVersion160,
	types.OpenEBSVersion210EE:  types.CSIProvisionerVersion160,
	types.OpenEBSVersion220:    types.CSIProvisionerVersion160,
	types.OpenEBSVersion220EE:  types.CSIProvisionerVersion160,
	types.OpenEBSVersion240:    types.CSIProvisionerVersion160,
	types.OpenEBSVersion250:    types.CSIProvisionerVersion210,
	types.OpenEBSVersion260:    types.CSIProvisionerVersion210,
	types.OpenEBSVersion270:    types.CSIProvisionerVersion210,
	types.OpenEBSVersion280:    types.CSIProvisionerVersion210,
}

// SupportedCSIAttacherVersionForCSIControllerVersion stores the mapping for
// CSI provisioner to CSIController version.
var SupportedCSIAttacherVersionForCSIControllerVersion = map[string]string{
	types.OpenEBSVersion190:    types.CSIAttacherVersion200,
	types.OpenEBSVersion190EE:  types.CSIAttacherVersion200,
	types.OpenEBSVersion1100:   types.CSIAttacherVersion200,
	types.OpenEBSVersion1100EE: types.CSIAttacherVersion200,
	types.OpenEBSVersion1110:   types.CSIAttacherVersion200,
	types.OpenEBSVersion1110EE: types.CSIAttacherVersion200,
	types.OpenEBSVersion1120:   types.CSIAttacherVersion200,
	types.OpenEBSVersion1120EE: types.CSIAttacherVersion200,
	types.OpenEBSVersion200:    types.CSIAttacherVersion200,
	types.OpenEBSVersion200EE:  types.CSIAttacherVersion200,
	types.OpenEBSVersion210:    types.CSIAttacherVersion200,
	types.OpenEBSVersion210EE:  types.CSIAttacherVersion200,
	types.OpenEBSVersion220:    types.CSIAttacherVersion200,
	types.OpenEBSVersion220EE:  types.CSIAttacherVersion200,
	types.OpenEBSVersion240:    types.CSIAttacherVersion200,
	types.OpenEBSVersion250:    types.CSIAttacherVersion310,
	types.OpenEBSVersion260:    types.CSIAttacherVersion310,
	types.OpenEBSVersion270:    types.CSIAttacherVersion310,
	types.OpenEBSVersion280:    types.CSIAttacherVersion310,
}

// SupportedCSIClusterDriverRegistrarVersionForOpenEBSVersion stores the mapping for
// CSIClusterDriverRegistrar to OpenEBS version.
var SupportedCSIClusterDriverRegistrarVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion190:    types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion190EE:  types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion1100:   types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion1100EE: types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion1110:   types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion1110EE: types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion1120:   types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion1120EE: types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion200:    types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion200EE:  types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion210:    types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion210EE:  types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion220:    types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion220EE:  types.CSIClusterDriverRegistrarVersion101,
	types.OpenEBSVersion240:    types.CSIClusterDriverRegistrarVersion101,
}

// SupportedCSINodeDriverRegistrarVersionForCSINodeVersion stores the mapping for
// CSINodeDriverRegistrar to CSI node version.
var SupportedCSINodeDriverRegistrarVersionForCSINodeVersion = map[string]string{
	types.OpenEBSVersion190:    types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion190EE:  types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion1100:   types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion1100EE: types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion1110:   types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion1110EE: types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion1120:   types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion1120EE: types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion200:    types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion200EE:  types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion210:    types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion210EE:  types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion220:    types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion220EE:  types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion240:    types.CSINodeDriverRegistrarVersion101,
	types.OpenEBSVersion250:    types.CSINodeDriverRegistrarVersion210,
	types.OpenEBSVersion260:    types.CSINodeDriverRegistrarVersion210,
	types.OpenEBSVersion270:    types.CSINodeDriverRegistrarVersion210,
	types.OpenEBSVersion280:    types.CSINodeDriverRegistrarVersion210,
}

// Set the default values for Cstor if not already given.
func (p *Planner) setCStorDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.CstorConfig == nil {
		p.ObservedOpenEBS.Spec.CstorConfig = &types.CstorConfig{}
	}
	// form the cstor-pool image
	if p.ObservedOpenEBS.Spec.CstorConfig.Pool.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.Pool.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	p.ObservedOpenEBS.Spec.CstorConfig.Pool.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-pool:" + p.ObservedOpenEBS.Spec.CstorConfig.Pool.ImageTag

	// form the cstor-pool-mgmt image
	if p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-pool-mgmt:" + p.ObservedOpenEBS.Spec.CstorConfig.PoolMgmt.ImageTag

	// form the cstor-istgt image
	if p.ObservedOpenEBS.Spec.CstorConfig.Target.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.Target.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	p.ObservedOpenEBS.Spec.CstorConfig.Target.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-istgt:" + p.ObservedOpenEBS.Spec.CstorConfig.Target.ImageTag

	// form the cstor-volume-mgmt image
	if p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"cstor-volume-mgmt:" + p.ObservedOpenEBS.Spec.CstorConfig.VolumeMgmt.ImageTag
	// form the cstor-volume-manager image
	volumeManagerImageName := "cstor-volume-manager-amd64:"
	if p.ObservedOpenEBS.Spec.CstorConfig.VolumeManager.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.VolumeManager.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	if p.ObservedOpenEBS.Spec.Version == types.OpenEBSVersion190 {
		volumeManagerImageName = "cstor-volume-mgmt:"
	} else if OpenEBSVersionAbove240 {
		volumeManagerImageName = "cstor-volume-manager:"
	}
	p.ObservedOpenEBS.Spec.CstorConfig.VolumeManager.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		volumeManagerImageName + p.ObservedOpenEBS.Spec.CstorConfig.VolumeManager.ImageTag
	// form the cspi-mgmt image(CSPI_MGMT)
	cspiImageName := "cstor-pool-manager-amd64:"
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	if p.ObservedOpenEBS.Spec.Version == types.OpenEBSVersion190 {
		cspiImageName = "cspi-mgmt:"
	} else if OpenEBSVersionAbove240 {
		cspiImageName = "cstor-pool-manager:"
	}
	p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		cspiImageName + p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.ImageTag

	// set the CSPC operator defaults
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator = &types.CSPCOperator{}
	}
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Enabled = true
	}
	// form the CSPC image
	cspcImage := "cspc-operator-amd64:"
	// set the name with which cspc-operator will be deployed
	if len(p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Name) == 0 {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Name = types.CSPCOperatorNameKey
	}
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	if p.ObservedOpenEBS.Spec.Version == types.OpenEBSVersion190 {
		cspcImage = "cspc-operator:"
	} else if OpenEBSVersionAbove240 {
		cspcImage = "cspc-operator:"
	}
	// form the container image as per the image prefix and image tag.
	p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		cspcImage + p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.ImageTag
	if p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Replicas == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Replicas = DefaultCSPCOperatorReplicaCount
	}
	// set the CVC operator defaults
	if p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator = &types.CVCOperator{}
	}
	if p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Enabled = true
	}
	// form the CVC image
	cvcImage := "cvc-operator-amd64:"
	// set the name with which cvc-operator will be deployed
	if len(p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Name) == 0 {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Name = types.CVCOperatorNameKey
	}
	if p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	if p.ObservedOpenEBS.Spec.Version == types.OpenEBSVersion190 {
		cvcImage = "cvc-operator:"
	} else if OpenEBSVersionAbove240 {
		cvcImage = "cvc-operator:"
	}
	// form the container image as per the image prefix and image tag.
	p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		cvcImage + p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.ImageTag
	if p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Replicas == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Replicas = DefaultCVCOperatorReplicaCount
	}
	if p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Service == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Service = &types.CVCOperatorService{}
	}
	if len(p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Service.Name) == 0 {
		p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Service.Name = types.CVCOperatorServiceNameKey
	}
	// set the admission server defaults
	if p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer = &types.CStorAdmissionServer{}
	}
	if p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Enabled = true
	}

	// set the name with which openebs-cstor-admission-server will be deployed
	if len(p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Name) == 0 {
		p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Name = types.CStorAdmissionServerNameKey
	}
	// form the CStor admission server image
	if p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.ImageTag == "" {
		p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	// form the cstor-webhook image
	cstorWebhookImage := "cstor-webhook-amd64:"
	if OpenEBSVersionAbove240 {
		cstorWebhookImage = "cstor-webhook:"
	}
	// form the container image as per the image prefix and image tag.
	p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		cstorWebhookImage + p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.ImageTag
	if p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Replicas == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Replicas = DefaultCStorAdmissionServerReplicaCount
	}

	err := p.setCSIDefaultsIfNotSet()
	if err != nil {
		return err
	}

	return nil
}

func (p *Planner) setCSIDefaultsIfNotSet() error {
	var (
		// List of images which are by default fetched from quay.io/k8scsi registry.
		CSIResizerImageTag                       string
		CSISnapshotterImageTag                   string
		CSISnapshotControllerImageTag            string
		CSIProvisionerForCSIControllerImageTag   string
		CSIAttacherForCSIControllerImageTag      string
		CSIClusterDriverRegistrarImageTag        string
		CSINodeDriverRegistrarForCSINodeImageTag string
	)
	isCSISupported, err := p.isCSISupported()
	// Do not return the error as not to block installing other components.
	if err != nil {
		isCSISupported = false
		glog.Errorf("Failed to set CSI defaults, error: %v", err)
	}

	if !isCSISupported {
		glog.V(5).Infof("Skipping CSI installation.")
	}
	// update the kubeletPath default value
	if len(p.ObservedOpenEBS.Spec.KubeletRootDirectory) == 0 {
		if p.ObservedOpenEBS.Spec.K8sDistribution == types.KeyMicroK8s {
			KubeletPath = "/var/snap/microk8s/common/var/lib/kubelet"
		} else {
			KubeletPath = "/var/lib/kubelet"
		}
	} else {
		KubeletPath = strings.TrimRight(p.ObservedOpenEBS.Spec.KubeletRootDirectory, "/")
	}
	// Set the default values for cstor csi controller statefulset in configuration.
	if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled = true
	}
	// set the name with which csi-controller will be deployed
	if len(p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Name) == 0 {
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Name = types.CStorCSIControllerNameKey
	}
	if !isCSISupported && *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled == true {
		*p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled = false
	}

	if *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled == true {
		if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.ImageTag == "" {
			p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.ImageTag = p.ObservedOpenEBS.Spec.Version +
				p.ObservedOpenEBS.Spec.ImageTagSuffix
		}
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
			"cstor-csi-driver:" + p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.ImageTag
	}

	// Set the default values for cstor csi node daemonset in configuration.
	if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled == nil {
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled = true
	}
	// set the name with which csi-node will be deployed
	if len(p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Name) == 0 {
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Name = types.CStorCSINodeNameKey
	}
	if !isCSISupported && *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled == true {
		*p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled = false
	}

	if *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled == true {
		if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ImageTag == "" {
			p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ImageTag = p.ObservedOpenEBS.Spec.Version +
				p.ObservedOpenEBS.Spec.ImageTagSuffix
		}
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
			"cstor-csi-driver:" + p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ImageTag

		if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ISCSIPath == "" {
			p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ISCSIPath = "/sbin/iscsiadm"
		}
	}
	// If CStor csi-controller or csi-node is enabled then check and delete the CSI components
	// if not installed in OpenEBS namespace after OpenEBS version 2.0.0.
	if *p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Enabled == true ||
		*p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Enabled == true {
		err = p.deleteCSIComponentsIfRequired()
		if err != nil {
			// only log the error, do not stop the flow.
			glog.Error(err)
		}
	}
	// form the csi-resizer image
	if csiResizerVersion, exist := SupportedCSIResizerVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
		CSIResizerImageTag = "csi-resizer:" + csiResizerVersion
	} else {
		return errors.Errorf("Failed to get csi-resizer version for the given OpenEBS version: %s",
			p.ObservedOpenEBS.Spec.Version)
	}

	// form the csi-snapshotter image
	if csiSnapshotterVersion, exist := SupportedCSISnapshotterVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
		CSISnapshotterImageTag = "csi-snapshotter:" + csiSnapshotterVersion
	} else {
		return errors.Errorf("Failed to get csi-snapshotter version for the given OpenEBS version: %s",
			p.ObservedOpenEBS.Spec.Version)
	}

	// form the CSI snapshot-controller image
	if csiSnapshotControllerVersion, exist := SupportedCSISnapshotControllerVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
		CSISnapshotControllerImageTag = "snapshot-controller:" + csiSnapshotControllerVersion
	} else {
		return errors.Errorf("Failed to get snapshot-controller version for the given OpenEBS version: %s",
			p.ObservedOpenEBS.Spec.Version)
	}

	// form the CSI provisioner image for the CSI controller
	if csiProvisionerForCSIController, exist :=
		SupportedCSIProvisionerVersionForCSIControllerVersion[p.ObservedOpenEBS.Spec.Version]; exist {
		CSIProvisionerForCSIControllerImageTag = "csi-provisioner:" +
			csiProvisionerForCSIController
	} else {
		return errors.Errorf(
			"Failed to get csi-provisioner version for csi-controller for the given OpenEBS version: %s",
			p.ObservedOpenEBS.Spec.Version)
	}

	// form the CSI attacher for CSI controller
	if csiAttacherForCSIController, exist :=
		SupportedCSIAttacherVersionForCSIControllerVersion[p.ObservedOpenEBS.Spec.Version]; exist {
		CSIAttacherForCSIControllerImageTag = "csi-attacher:" +
			csiAttacherForCSIController
	} else {
		return errors.Errorf(
			"Failed to get csi-attacher version for csi-controller for the given OpenEBS version: %s",
			p.ObservedOpenEBS.Spec.Version)
	}

	// csi-cluster-driver-registrar container is present in cstor-csi-controller till
	// OpenEBS version 2.4.0 only.
	if !OpenEBSVersionAbove240 {
		// form the csi-cluster-driver-registrar image for the given OpenEBS version
		if csiClusterDriverRegistrar, exist :=
			SupportedCSIClusterDriverRegistrarVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
			CSIClusterDriverRegistrarImageTag = "csi-cluster-driver-registrar:" +
				csiClusterDriverRegistrar
		} else {
			return errors.Errorf(
				"Failed to get csi-cluster-driver-registrar version for the given OpenEBS version: %s",
				p.ObservedOpenEBS.Spec.Version)
		}
	}
	// form the csi-node-driver-registrar image for CSI node for the given OpenEBS version
	if csiNodeDriverRegistrar, exist :=
		SupportedCSINodeDriverRegistrarVersionForCSINodeVersion[p.ObservedOpenEBS.Spec.Version]; exist {
		CSINodeDriverRegistrarForCSINodeImageTag = "csi-node-driver-registrar:" +
			csiNodeDriverRegistrar
	} else {
		return errors.Errorf(
			"Failed to get csi-node-driver-registrar version for csi-node for the given OpenEBS version: %s",
			p.ObservedOpenEBS.Spec.Version)
	}

	// set the name of cstor-csi-iscsiadm configMap.
	if len(p.ObservedOpenEBS.Spec.CstorConfig.CSI.ISCSIADMConfigmap.Name) == 0 {
		p.ObservedOpenEBS.Spec.CstorConfig.CSI.ISCSIADMConfigmap.Name = types.CStorCSIISCSIADMConfigmapNameKey
	}

	// check if the image registry is the default ones i.e., quay.io/openebs/, openebs/ or mayadataio/,
	// if not then form the k8s repositories related images also so that they can also be pulled from
	// the specified repository only.
	if !(p.ObservedOpenEBS.Spec.ImagePrefix == types.QUAYIOOPENEBSREGISTRY ||
		p.ObservedOpenEBS.Spec.ImagePrefix == types.MAYADATAIOREGISTRY ||
		p.ObservedOpenEBS.Spec.ImagePrefix == types.OPENEBSREGISTRY) {
		CSIResizerImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSIResizerImageTag
		CSISnapshotterImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSISnapshotterImageTag
		CSISnapshotControllerImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSISnapshotControllerImageTag
		CSIProvisionerForCSIControllerImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSIProvisionerForCSIControllerImageTag
		CSIAttacherForCSIControllerImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSIAttacherForCSIControllerImageTag
		CSIClusterDriverRegistrarImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSIClusterDriverRegistrarImageTag
		CSINodeDriverRegistrarForCSINodeImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSINodeDriverRegistrarForCSINodeImageTag
	} else {
		CSISnapshotterImageRegistry := types.QUAYIOK8SCSI
		CSISnapshotControllerImageRegistry := types.QUAYIOK8SCSI
		CSINodeDriverRegistrarForCSINodeImageRegistry := types.QUAYIOK8SCSI
		// For OpenEBS version 2.5.0 or greater, csi-snapshotter and snapshot-controller images
		// are pulled from k8s.gcr.io/sig-storage registry instead of quay.io/k8scsi registry.
		if OpenEBSVersionAbove240 {
			CSISnapshotterImageRegistry = types.K8SGCRSIGSTORAGE
			CSISnapshotControllerImageRegistry = types.K8SGCRSIGSTORAGE
			CSINodeDriverRegistrarForCSINodeImageRegistry = types.K8SGCRSIGSTORAGE
		}
		CSIResizerImage = types.QUAYIOK8SCSI + CSIResizerImageTag
		CSIResizerImage = types.QUAYIOK8SCSI + CSIResizerImageTag
		CSISnapshotterImage = CSISnapshotterImageRegistry + CSISnapshotterImageTag
		CSISnapshotControllerImage = CSISnapshotControllerImageRegistry + CSISnapshotControllerImageTag
		CSIProvisionerForCSIControllerImage = types.QUAYIOK8SCSI + CSIProvisionerForCSIControllerImageTag
		CSIAttacherForCSIControllerImage = types.QUAYIOK8SCSI + CSIAttacherForCSIControllerImageTag
		CSIClusterDriverRegistrarImage = types.QUAYIOK8SCSI + CSIClusterDriverRegistrarImageTag
		CSINodeDriverRegistrarForCSINodeImage = CSINodeDriverRegistrarForCSINodeImageRegistry +
			CSINodeDriverRegistrarForCSINodeImageTag
	}

	return nil
}

// Check the OpenEBS version if it is greater than 2.0.0, if yes then check if CSI components
// are already installed at kube-system, if yes then install CSI components in openebs namespace
// and delete from kube-system.
func (p *Planner) deleteCSIComponentsIfRequired() error {
	// check if CSI is supported or not for this version, if not then do not delete the existing ones.
	isCSISupported, err := p.isCSISupported()
	if err != nil {
		return errors.Errorf("Error checking if CSI is supported or not: %+v", err)
	}
	// get the namespace where CSI components are going to be installed.
	csiNamespace, err := p.getCSIComponentsNamespace()
	if err != nil {
		return errors.Errorf(
			"Error getting the namespace where CSI components will be installed: %+v", err)
	}
	if isCSISupported && csiNamespace != types.NamespaceKubeSystem {
		// check if csi-components are already installed in kube-system namespace.
		for _, observedOpenEBSComp := range p.ObservedOpenEBSComponents {
			if observedOpenEBSComp.GetKind() == types.KindStatefulset ||
				observedOpenEBSComp.GetKind() == types.KindDaemonSet ||
				observedOpenEBSComp.GetKind() == types.KindServiceAccount {
				if observedOpenEBSComp.GetName() == types.CStorCSIControllerNameKey ||
					observedOpenEBSComp.GetName() == types.CStorCSINodeNameKey ||
					observedOpenEBSComp.GetName() == types.CStorCSIControllerSANameKey ||
					observedOpenEBSComp.GetName() == types.CStorCSINodeSANameKey {
					if observedOpenEBSComp.GetNamespace() == types.NamespaceKubeSystem {
						p.ExplicitDeletes = append(p.ExplicitDeletes, observedOpenEBSComp)
					}
				}
			}
		}
	}
	return nil
}

// isCSISupported checks if csi is supported or not in the current kubernetes cluster, if not it will
// return false else true.
func (p *Planner) isCSISupported() (bool, error) {
	// comp stores the result for comparing 2 versions
	var comp int
	// get the kubernetes version.
	k8sVersion, err := k8s.GetK8sVersion()
	if err != nil {
		return false, errors.Errorf("Unable to find kubernetes version, error: %v", err)
	}

	// compare the kubernetes version with the supported version of csi.
	comp, err = compareVersion(k8sVersion, types.CSISupportedVersion)
	if err != nil {
		return false, errors.Errorf("Error comparing versions, error: %v", err)
	}
	if comp < 0 {
		glog.Warningf("CSI is not supported in %s Kubernetes version. "+
			"CSI is supported from kubernetes version %s.", k8sVersion, types.CSISupportedVersion)
		return false, nil
	}

	return true, nil
}

// getCSIComponentsNamespace returns the namespace in which CSI components are going to be installed.
func (p *Planner) getCSIComponentsNamespace() (string, error) {
	var (
		// csiNamespace is the namespace where CSI components should be installed
		csiNamespace string
		// comp stores the result for comparing 2 versions
		comp int
	)
	// get the kubernetes version.
	k8sVersion, err := k8s.GetK8sVersion()
	if err != nil {
		return csiNamespace, errors.Errorf("Unable to find kubernetes version, error: %v", err)
	}
	// Check if the given OpenEBS version is greater than or less than OpenEBS version 2.0.0.
	// For OpenEBS version 2.0.0 or greater, CSI components can be installed in openebs namespace
	// if k8s version is greater than or equal to 1.17.0.
	res, err := compareVersion(p.ObservedOpenEBS.Spec.Version, types.OpenEBSVersion200)
	if err != nil {
		return csiNamespace, errors.Errorf(
			"Error comparing versions while determining namespace for CSI components[v1: %s, v2: %s], error: %v",
			p.ObservedOpenEBS.Spec.Version, types.OpenEBSVersion200, err)
	}
	if res >= 0 {
		// compare the kubernetes version with the supported version of csi.
		comp, err = compareVersion(k8sVersion, types.CSISupportedVersionFromOpenEBS200)
		if err != nil {
			return csiNamespace, errors.Errorf("Error comparing versions, error: %v", err)
		}
		if comp < 0 {
			csiNamespace = types.NamespaceKubeSystem
		} else {
			// set the namespace in which OpenEBS components are going to be installed.
			csiNamespace = p.ObservedOpenEBS.Namespace
		}
	} else {
		// compare the kubernetes version with the supported version of csi.
		comp, err = compareVersion(k8sVersion, types.CSISupportedVersion)
		if err != nil {
			return csiNamespace, errors.Errorf("Error comparing versions, error: %v", err)
		}
		if comp < 0 {
			return csiNamespace, errors.Errorf("CSI is not supported in %s Kubernetes version. "+
				"CSI is supported from kubernetes version %s.", k8sVersion, types.CSISupportedVersion)
		} else {
			csiNamespace = types.NamespaceKubeSystem
		}
	}
	return csiNamespace, nil
}

// updateOpenEBSCStorCSINode updates the values of openebs-cstor-csi-node daemonset as per given configuration.
func (p *Planner) updateOpenEBSCStorCSINode(daemonset *unstructured.Unstructured) error {
	var (
		extraVolumes      []interface{}
		extraVolumeMounts []interface{}
	)
	daemonset.SetName(p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Name)
	// set the namespace in which CSI components should be installed.
	csiNamespace, err := p.getCSIComponentsNamespace()
	if err != nil {
		return err
	}
	daemonset.SetNamespace(csiNamespace)

	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := daemonset.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor CSI Node:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-node
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorCSINodeNameKey
	// set the desired labels
	daemonset.SetLabels(desiredLabels)

	// this will get the extra volumes and volume mounts required to be added in the csi node daemonset
	// for the csi to work for different OS distributions/versions.
	// This volumes and volume mounts will be added in the openebs-csi-plugin container.
	comp, err := compareVersion(p.ObservedOpenEBS.Spec.Version, types.OpenEBSVersion200)
	if err != nil {
		return err
	}
	if comp < 0 {
		extraVolumes, extraVolumeMounts, err = p.getOSSpecificVolumeMounts()
		if err != nil {
			return err
		}
	}

	volumes, err := unstruct.GetNestedSliceOrError(daemonset, "spec", "template", "spec", "volumes")
	if err != nil {
		return err
	}
	// updateVolume updates the volume path of openebs-csi-plugin container.
	updateVolume := func(obj *unstructured.Unstructured) error {
		volumeName, err := unstruct.GetString(obj, "spec", "name")
		if err != nil {
			return err
		}
		if volumeName == "iscsiadm-bin" {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ISCSIPath, "spec", "hostPath", "path")
			if err != nil {
				return err
			}
		} else if volumeName == "registration-dir" {
			err = unstructured.SetNestedField(obj.Object,
				KubeletPath+"/plugins_registry/", "spec", "hostPath", "path")
			if err != nil {
				return err
			}
		} else if volumeName == "plugin-dir" {
			err = unstructured.SetNestedField(obj.Object,
				KubeletPath+"/plugins/cstor.csi.openebs.io/", "spec", "hostPath", "path")
			if err != nil {
				return err
			}
		} else if volumeName == "pods-mount-dir" {
			err = unstructured.SetNestedField(obj.Object,
				KubeletPath+"/", "spec", "hostPath", "path")
			if err != nil {
				return err
			}
		}

		return nil
	}
	err = unstruct.SliceIterator(volumes).ForEachUpdate(updateVolume)
	if err != nil {
		return err
	}

	// Append the new extra volumes with the existing volumes, required for the csi to work.
	volumes = append(volumes, extraVolumes...)

	err = unstructured.SetNestedSlice(daemonset.Object, volumes,
		"spec", "template", "spec", "volumes")
	if err != nil {
		return err
	}

	containers, err := unstruct.GetNestedSliceOrError(daemonset, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	// update the env value of openebs-csi-plugin container
	updateOpenEBSCSIPluginEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == EnvOpenEBSNamespaceKey {
			unstructured.SetNestedField(env.Object, p.ObservedOpenEBS.Namespace, "spec", "value")
		}
		return nil
	}

	// update the env value of csi-node-driver-registrar container
	updateCSINodeDriverRegistrarEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == EnvDriverRegSocketPathKey {
			unstructured.SetNestedField(env.Object, KubeletPath+"/plugins/cstor.csi.openebs.io/csi.sock",
				"spec", "value")
		}
		return nil
	}

	// updateOpenEBSCSIPluginVolumeMount updates the volumeMounts path of openebs-csi-plugin container.
	updateOpenEBSCSIPluginVolumeMount := func(vm *unstructured.Unstructured) error {
		vmName, _, err := unstructured.NestedString(vm.Object, "spec", "name")
		if err != nil {
			return err
		}
		if vmName == "iscsiadm-bin" {
			err = unstructured.SetNestedField(vm.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ISCSIPath, "spec", "mountPath")
			if err != nil {
				return err
			}
		} else if vmName == "pods-mount-dir" {
			err = unstructured.SetNestedField(vm.Object,
				KubeletPath+"/", "spec", "mountPath")
			if err != nil {
				return err
			}
		}
		return nil
	}

	// update the containers
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		volumeMounts, _, err := unstruct.GetSlice(obj, "spec", "volumeMounts")
		if err != nil {
			return err
		}

		if containerName == ContainerOpenEBSCSIPluginName ||
			containerName == ContainerCSTORCSIPluginName {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Image,
				"spec", "image")
			if err != nil {
				return err
			}
			// Set the environments of the container.
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateOpenEBSCSIPluginEnv)
			if err != nil {
				return err
			}
			envs, err = p.ignoreUpdatingImmutableEnvs(p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ENV, envs)
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(volumeMounts).ForEachUpdate(updateOpenEBSCSIPluginVolumeMount)
			if err != nil {
				return err
			}
		} else if containerName == ContainerCSINodeDriverRegistrarName {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, CSINodeDriverRegistrarForCSINodeImage,
				"spec", "image")
			if err != nil {
				return err
			}
			// Set the environments of the container.
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateCSINodeDriverRegistrarEnv)
			if err != nil {
				return err
			}
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}

		// Append the new extra volume mounts with the existing volume mounts, required for the csi to work.
		volumeMounts = append(volumeMounts, extraVolumeMounts...)
		err = unstructured.SetNestedSlice(obj.Object, volumeMounts, "spec", "volumeMounts")
		if err != nil {
			return err
		}

		// Set the resource of the containers.
		if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Resources != nil {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.Resources,
				"spec", "resources")
		} else if p.ObservedOpenEBS.Spec.Resources != nil {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.Resources, "spec", "resources")
		}
		if err != nil {
			return err
		}

		return nil
	}

	// Update the containers.
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}

	// Set back the value of the containers.
	err = unstructured.SetNestedSlice(daemonset.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	// add init-container to check for presence of ISCSI client on the node on which
	// the pod of this resource is running.
	err = addISCSIClientInitContainer(daemonset)
	if err != nil {
		return err
	}

	return nil
}

// addISCSIClientInitContainer adds ISCSI client init-container to a given
// resource so that the resource first checks for presence of ISCSI client before
// initialization of the pod i.e., a resource adding this init-container will run
// if and only if ISCSI client is running on the node on which the pod of this
// resource is running.
func addISCSIClientInitContainer(resource *unstructured.Unstructured) error {
	var (
		isOSSupported bool
		err           error
	)
	// check if the underlying OS is supported or not for ISCSI client setup.
	// get the OS image running on the underlying node
	osImage, err := k8s.GetOSImage()
	if err != nil {
		return errors.Errorf("[ISCSI client initContainer]Error getting OS image of a node, error: %+v", err)
	}
	// make an array of supported OSes
	supportedOSes := []string{"ubuntu", "Red Hat Enterprise Linux", "centos", "amazon linux"}
	for _, supportedOS := range supportedOSes {
		if strings.Contains(strings.ToLower(osImage), supportedOS) {
			isOSSupported = true
			break
		}
	}
	// ISCSI client init-container will not be added to the CSI node components if the underlying
	// OS is not supported for setting up ISCSI client.
	if !isOSSupported {
		return nil
	}
	// get the existing init-container if any
	initContainers, _ := unstruct.GetNestedSliceOrEmpty(resource, "spec", "template", "spec", "initContainers")

	// define and add the ISCSI client init-container.
	ISCSIInitContainer := map[string]interface{}{
		"name":  "init-node",
		"image": "alpine:3.7",
		"securityContext": map[string]interface{}{
			"privileged": true,
		},
		"command": []interface{}{
			"nsenter",
			"--mount=/proc/1/ns/mnt",
			"--",
			"sh",
			"-c",
			"until sudo systemctl status iscsid; do echo waiting for ISCSI client; sleep 2; done;",
		},
	}
	initContainers = append(initContainers, ISCSIInitContainer)
	// Set back the value of the containers.
	err = unstructured.SetNestedSlice(resource.Object,
		initContainers, "spec", "template", "spec", "initContainers")
	if err != nil {
		return err
	}
	// set hostPID: true as it is required to run the above init container.
	err = unstructured.SetNestedField(resource.Object,
		true, "spec", "template", "spec", "hostPID")
	if err != nil {
		return err
	}
	return nil
}

// updateCStorCSIISCSIADMConfig updates/sets the default values for cstor-csi-iscsiadm
// configmap as per the values provided in the OpenEBS CR.
func (p *Planner) updateCStorCSIISCSIADMConfig(configmap *unstructured.Unstructured) error {
	configmap.SetName(p.ObservedOpenEBS.Spec.CstorConfig.CSI.ISCSIADMConfigmap.Name)
	// set the namespace in which CSI components should be installed.
	csiNamespace, err := p.getCSIComponentsNamespace()
	if err != nil {
		return err
	}
	configmap.SetNamespace(csiNamespace)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := configmap.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for openebs-ndm-config configmap
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-iscsiadm
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorCSIISCSIADMConfigmapNameKey
	// set the desired labels
	configmap.SetLabels(desiredLabels)

	return nil
}

// getOSSpecificVolumeMounts returns the volume and volume mounts based on the specific OS distribution/version.
// This volume and volume mounts are for the specific container i.e openebs-csi-plugin.
// This function will get the OS Image and the version for the ubuntu distribution and will return the
// volumes and volume mounts accordngly.
func (p *Planner) getOSSpecificVolumeMounts() ([]interface{}, []interface{}, error) {
	volumes := make([]interface{}, 0)
	volumeMounts := make([]interface{}, 0)

	osImage, err := k8s.GetOSImage()
	if err != nil {
		return volumes, volumeMounts, errors.Errorf("Error getting OS Image of a Node, error: %+v", err)
	}

	ubuntuVersion, err := k8s.GetUbuntuVersion()
	if err != nil {
		return volumes, volumeMounts, errors.Errorf("Error getting Ubuntu Version of a Node, error: %+v", err)
	}

	switch true {
	case strings.Contains(strings.ToLower(osImage), strings.ToLower(types.OSImageSLES12)):
		volumes, volumeMounts = p.getSUSE12VolumeMounts()
	case strings.Contains(strings.ToLower(osImage), strings.ToLower(types.OSImageSLES15)):
		volumes, volumeMounts = p.getSUSE15VolumeMounts()
	case strings.Contains(strings.ToLower(osImage), strings.ToLower(types.OSImageUbuntu1804)) ||
		((ubuntuVersion != 0) && ubuntuVersion >= 18.04):
		volumes, volumeMounts = p.getUbuntu1804VolumeMounts()
	}

	return volumes, volumeMounts, nil
}

// getSUSE12VolumeMounts returns the volumes and volume mounts for suse 12.
func (p *Planner) getSUSE12VolumeMounts() ([]interface{}, []interface{}) {
	volumes := make([]interface{}, 0)
	volumeMounts := make([]interface{}, 0)

	// Create new volumes for suse 12.
	libCryptoVolume := map[string]interface{}{
		"name": "iscsiadm-lib-crypto",
		"hostPath": map[string]interface{}{
			"type": "File",
			"path": "/lib64/libcrypto.so.1.0.0",
		},
	}
	libOpeniscsiusrVolume := map[string]interface{}{
		"name": "iscsiadm-lib-openiscsiusr",
		"hostPath": map[string]interface{}{
			"type": "File",
			"path": "/usr/lib64/libopeniscsiusr.so.0.2.0",
		},
	}
	volumes = append(volumes, libCryptoVolume, libOpeniscsiusrVolume)

	// Create new volume mounts for suse 12.
	libCryptoVolumeMount := map[string]interface{}{
		"name":      "iscsiadm-lib-crypto",
		"mountPath": "/lib/x86_64-linux-gnu/libcrypto.so.1.0.0",
	}
	libOpeniscsiusrVolumeMount := map[string]interface{}{
		"name":      "iscsiadm-lib-openiscsiusr",
		"mountPath": "/lib/x86_64-linux-gnu/libopeniscsiusr.so.0.2.0",
	}
	volumeMounts = append(volumeMounts, libCryptoVolumeMount, libOpeniscsiusrVolumeMount)

	return volumes, volumeMounts
}

// getSUSE15VolumeMounts returns the volumes and volume mounts for suse 15.
func (p *Planner) getSUSE15VolumeMounts() ([]interface{}, []interface{}) {
	volumes := make([]interface{}, 0)
	volumeMounts := make([]interface{}, 0)

	// Create new volumes for suse 15.
	libCryptoVolume := map[string]interface{}{
		"name": "iscsiadm-lib-crypto",
		"hostPath": map[string]interface{}{
			"type": "File",
			"path": "/usr/lib64/libcrypto.so.1.1",
		},
	}
	libOpeniscsiusrVolume := map[string]interface{}{
		"name": "iscsiadm-lib-openiscsiusr",
		"hostPath": map[string]interface{}{
			"type": "File",
			"path": "/usr/lib64/libopeniscsiusr.so.0.2.0",
		},
	}
	volumes = append(volumes, libCryptoVolume, libOpeniscsiusrVolume)

	// Create new volume mounts for suse 15.
	libCryptoVolumeMount := map[string]interface{}{
		"name":      "iscsiadm-lib-crypto",
		"mountPath": "/lib/x86_64-linux-gnu/libcrypto.so.1.1",
	}
	libOpeniscsiusrVolumeMount := map[string]interface{}{
		"name":      "iscsiadm-lib-openiscsiusr",
		"mountPath": "/lib/x86_64-linux-gnu/libopeniscsiusr.so.0.2.0",
	}
	volumeMounts = append(volumeMounts, libCryptoVolumeMount, libOpeniscsiusrVolumeMount)

	return volumes, volumeMounts
}

// getUbuntu1804VolumeMounts returns the volumes and volume mounts for ubuntu 18.04 and above.
func (p *Planner) getUbuntu1804VolumeMounts() ([]interface{}, []interface{}) {
	volumes := make([]interface{}, 0)
	volumeMounts := make([]interface{}, 0)

	// Create new volume for ubuntu 18.04 and above.
	volume := map[string]interface{}{
		"name": "iscsiadm-lib-isns-nocrypto",
		"hostPath": map[string]interface{}{
			"type": "File",
			"path": "/lib/x86_64-linux-gnu/libisns-nocrypto.so.0",
		},
	}
	volumes = append(volumes, volume)

	// Create new volume mount for ubuntu 18.04 and above.
	volumeMount := map[string]interface{}{
		"name":      "iscsiadm-lib-isns-nocrypto",
		"mountPath": "/lib/x86_64-linux-gnu/libisns-nocrypto.so.0",
	}
	volumeMounts = append(volumeMounts, volumeMount)

	return volumes, volumeMounts
}

// updateOpenEBSCStorCSIController updates the values of openebs-cstor-csi-controller statefulset as per given configuration.
func (p *Planner) updateOpenEBSCStorCSIController(statefulset *unstructured.Unstructured) error {
	statefulset.SetName(p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Name)
	// set the namespace in which CSI components should be installed.
	csiNamespace, err := p.getCSIComponentsNamespace()
	if err != nil {
		return err
	}
	statefulset.SetNamespace(csiNamespace)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := statefulset.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor CSI controller:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cstor-csi
	// 2. openebs-upgrade.dao.mayadata.io/component-name: openebs-cstor-csi-controller
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSCStorCSIComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] =
		types.CStorCSIControllerNameKey
	// set the desired labels
	statefulset.SetLabels(desiredLabels)

	containers, err := unstruct.GetNestedSliceOrError(statefulset, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	// update the containers
	err = unstruct.SliceIterator(containers).ForEachUpdate(p.updateOpenEBSCStorCSIControllerContainers())
	if err != nil {
		return err
	}

	err = unstructured.SetNestedSlice(statefulset.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// updateOpenEBSCSIControllerCSIPluginEnv updates the env value of openebs-csi-plugin container of openebs-cstor-csi-controller
func (p *Planner) updateOpenEBSCSIControllerCSIPluginEnv() func(obj *unstructured.Unstructured) error {
	return func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == EnvOpenEBSNamespaceKey {
			unstructured.SetNestedField(env.Object, p.ObservedOpenEBS.Namespace, "spec", "value")
		}
		return nil
	}
}

// updateOpenEBSCStorCSIControllerContainers updates the containers of openebs-cstor-csi-controller
func (p *Planner) updateOpenEBSCStorCSIControllerContainers() func(obj *unstructured.Unstructured) error {
	return func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		switch containerName {
		case ContainerOpenEBSCSIPluginName:
			err = p.updateOpenEBSCSIControllerCSIPluginContainer(obj)
		case ContainerCSTORCSIPluginName:
			err = p.updateOpenEBSCSIControllerCSIPluginContainer(obj)
		case ContainerCSIResizerName:
			err = p.updateOpenEBSCSIControllerCSIResizerContainer(obj)
		case ContainerCSISnapshotterName:
			err = p.updateOpenEBSCSIControllerCSISnapshotterContainer(obj)
		case ContainerCSISnapshotControllerName:
			err = p.updateOpenEBSCSIControllerCSISnapshotControllerContainer(obj)
		case ContainerCSIAttacherName:
			err = p.updateOpenEBSCSIControllerCSIAttacherContainer(obj)
		case ContainerCSIProvisionerName:
			err = p.updateOpenEBSCSIControllerCSIProvisionerContainer(obj)
		case ContainerCSIClusterDriverRegistrarName:
			err = p.updateOpenEBSCSIControllerCSIClusterRegistrarDriverContainer(obj)
		}
		if err != nil {
			return err
		}

		// Set the resource of the containers.
		if p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Resources != nil {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Resources,
				"spec", "resources")
		} else if p.ObservedOpenEBS.Spec.Resources != nil {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.Resources, "spec", "resources")
		}
		if err != nil {
			return err
		}

		return nil
	}
}

// updateOpenEBSCSIControllerCSIResizerContainer updates the csi-resizer container such as the image,
// env, etc.
func (p *Planner) updateOpenEBSCSIControllerCSIResizerContainer(obj *unstructured.Unstructured) error {
	// Set the image of the container.
	err := unstructured.SetNestedField(obj.Object, CSIResizerImage,
		"spec", "image")
	if err != nil {
		return err
	}

	return nil
}

// updateOpenEBSCSIControllerCSISnapshotterContainer updates the csi-snapshotter container such as the image,
// env, etc.
func (p *Planner) updateOpenEBSCSIControllerCSISnapshotterContainer(obj *unstructured.Unstructured) error {
	// Set the image of the container.
	err := unstructured.SetNestedField(obj.Object, CSISnapshotterImage,
		"spec", "image")
	if err != nil {
		return err
	}

	return nil
}

// updateOpenEBSCSIControllerCSISnapshotControllerContainer updates the snapshot-controller container such as the image,
// env, etc.
func (p *Planner) updateOpenEBSCSIControllerCSISnapshotControllerContainer(obj *unstructured.Unstructured) error {
	// Set the image of the container.
	err := unstructured.SetNestedField(obj.Object, CSISnapshotControllerImage,
		"spec", "image")
	if err != nil {
		return err
	}

	return nil
}

// updateOpenEBSCSIControllerCSIProvisionerContainer updates the csi-provisioner container such as the image,
// env, etc.
func (p *Planner) updateOpenEBSCSIControllerCSIProvisionerContainer(obj *unstructured.Unstructured) error {
	// Set the image of the container.
	err := unstructured.SetNestedField(obj.Object, CSIProvisionerForCSIControllerImage,
		"spec", "image")
	if err != nil {
		return err
	}

	return nil
}

// updateOpenEBSCSIControllerCSIAttacherContainer updates the csi-attacher container such as the image,
// env, etc.
func (p *Planner) updateOpenEBSCSIControllerCSIAttacherContainer(obj *unstructured.Unstructured) error {
	// Set the image of the container.
	err := unstructured.SetNestedField(obj.Object, CSIAttacherForCSIControllerImage,
		"spec", "image")
	if err != nil {
		return err
	}

	return nil
}

// updateOpenEBSCSIControllerCSIClusterRegistrarDriverContainer updates the csi-cluster-registrar-driver
// container such as the image, env, etc.
func (p *Planner) updateOpenEBSCSIControllerCSIClusterRegistrarDriverContainer(obj *unstructured.Unstructured) error {
	// Set the image of the container.
	err := unstructured.SetNestedField(obj.Object, CSIClusterDriverRegistrarImage,
		"spec", "image")
	if err != nil {
		return err
	}

	return nil
}

// updateOpenEBSCSIControllerCSIPluginContainer updates the openebs-csi-plugin container such as the image,
// env, etc of openebs-cstor-csi-controller.
func (p *Planner) updateOpenEBSCSIControllerCSIPluginContainer(obj *unstructured.Unstructured) error {
	envs, _, err := unstruct.GetSlice(obj, "spec", "env")
	if err != nil {
		return err
	}
	// Set the image of the container.
	err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.Image,
		"spec", "image")
	if err != nil {
		return err
	}
	// Set the environments of the container.
	err = unstruct.SliceIterator(envs).ForEachUpdate(p.updateOpenEBSCSIControllerCSIPluginEnv())
	if err != nil {
		return err
	}
	envs, err = p.ignoreUpdatingImmutableEnvs(p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.ENV, envs)
	if err != nil {
		return err
	}
	err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
	if err != nil {
		return err
	}

	return nil
}

// updateCSPCOperator updates the CSPC operator manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateCSPCOperator(deploy *unstructured.Unstructured) error {
	deploy.SetName(p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for cspc-operator deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cspc
	// 2. openebs-upgrade.dao.mayadata.io/component-name: cspc-operator
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CSPCComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CSPCOperatorNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	// get the containers of the cspc-operator and update the desired fields
	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	// update the env value of cspc-operator container
	updateCSPCOperatorEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == "OPENEBS_IO_BASE_DIR" {
			err = unstructured.SetNestedField(env.Object, p.ObservedOpenEBS.Spec.DefaultStoragePath,
				"spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_POOL_SPARSE_DIR" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath+"/sparse", "spec", "value")
		} else if envName == "OPENEBS_IO_CSPI_MGMT_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.CSPIMgmt.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_POOL_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.Pool.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_POOL_EXPORTER_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.Policies.Monitoring.Image, "spec", "value")
		}
		if err != nil {
			return err
		}

		return nil
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		// update the envs of cspc-operator container
		// In order to update envs of other containers, just write an updateEnv
		// function for specific containers.
		if containerName == "cspc-operator" {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.Image,
				"spec", "image")
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateCSPCOperatorEnv)
			if err != nil {
				return err
			}
			envs, err = p.ignoreUpdatingImmutableEnvs(p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.ENV, envs)
			if err != nil {
				return err
			}
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	err = unstructured.SetNestedSlice(deploy.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// updateCVCOperator updates the CVC operator manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateCVCOperator(deploy *unstructured.Unstructured) error {
	deploy.SetName(p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for cvc-operator deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cvc
	// 2. openebs-upgrade.dao.mayadata.io/component-name: cvc-operator
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CVCComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CVCOperatorNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	// get the containers of the cvc-operator and update the desired fields
	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	// update the env value of cvc-operator container
	updateCVCOperatorEnv := func(env *unstructured.Unstructured) error {
		envName, _, err := unstructured.NestedString(env.Object, "spec", "name")
		if err != nil {
			return err
		}
		if envName == "OPENEBS_IO_BASE_DIR" {
			err = unstructured.SetNestedField(env.Object, p.ObservedOpenEBS.Spec.DefaultStoragePath,
				"spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_TARGET_DIR" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.DefaultStoragePath+"/sparse", "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_TARGET_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.Target.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_CSTOR_VOLUME_MGMT_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.VolumeManager.Image, "spec", "value")
		} else if envName == "OPENEBS_IO_VOLUME_MONITOR_IMAGE" {
			err = unstructured.SetNestedField(env.Object,
				p.ObservedOpenEBS.Spec.Policies.Monitoring.Image, "spec", "value")
		}
		if err != nil {
			return err
		}

		return nil
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		// update the envs of cvc-operator container
		// In order to update envs of other containers, just write an updateEnv
		// function for specific containers.
		if containerName == "cvc-operator" {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Image,
				"spec", "image")
			if err != nil {
				return err
			}
			err = unstruct.SliceIterator(envs).ForEachUpdate(updateCVCOperatorEnv)
			if err != nil {
				return err
			}
			envs, err = p.ignoreUpdatingImmutableEnvs(p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.ENV, envs)
			if err != nil {
				return err
			}
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	err = unstructured.SetNestedSlice(deploy.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// updateCStorAdmissionServer updates the CStorAdmissionServer manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateCStorAdmissionServer(deploy *unstructured.Unstructured) error {
	deploy.SetName(p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for CStor admissionServer deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-name: cstor-admission-webhook
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CStorAdmissionServerComponentNameLabelValue
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	// update the containers
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		if containerName == "admission-webhook" {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.Image,
				"spec", "image")
			if err != nil {
				return err
			}
			envs, err = p.ignoreUpdatingImmutableEnvs(p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.ENV, envs)
			if err != nil {
				return err
			}
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}
		return nil
	}
	// Update the containers.
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	// Set back the value of the containers.
	err = unstructured.SetNestedSlice(deploy.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// updateCVCOperatorService updates the cvc-operator-service manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateCVCOperatorService(svc *unstructured.Unstructured) error {
	svc.SetName(p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.Service.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := svc.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for cvc-operator-service service
	// 1. openebs-upgrade.dao.mayadata.io/component-group: cvc
	// 2. openebs-upgrade.dao.mayadata.io/component-name: cvc-operator-service
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.CVCComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.CVCOperatorServiceNameKey
	// set the desired labels
	svc.SetLabels(desiredLabels)

	return nil
}

func (p *Planner) fillCSPCOperatorExistingValues(observedComponentDetails ObservedComponentDesiredDetails) error {
	var (
		containerName string
		err           error
	)
	p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.MatchLabels = observedComponentDetails.MatchLabels
	p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.PodTemplateLabels = observedComponentDetails.PodTemplateLabels
	if len(p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.ContainerName) > 0 {
		containerName = p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.ContainerName
	} else {
		containerName = types.CSPCOperatorContainerKey
	}
	p.ObservedOpenEBS.Spec.CstorConfig.CSPCOperator.ENV, err = fetchExistingContainerEnvs(
		observedComponentDetails.Containers, containerName)
	if err != nil {
		return err
	}

	return nil
}

func (p *Planner) fillCVCOperatorExistingValues(observedComponentDetails ObservedComponentDesiredDetails) error {
	var (
		containerName string
		err           error
	)
	p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.MatchLabels = observedComponentDetails.MatchLabels
	p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.PodTemplateLabels = observedComponentDetails.PodTemplateLabels
	if len(p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.ContainerName) > 0 {
		containerName = p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.ContainerName
	} else {
		containerName = types.CVCOperatorContainerKey
	}
	p.ObservedOpenEBS.Spec.CstorConfig.CVCOperator.ENV, err = fetchExistingContainerEnvs(
		observedComponentDetails.Containers, containerName)
	if err != nil {
		return err
	}

	return nil
}

func (p *Planner) fillCStorAdmissionServerExistingValues(observedComponentDetails ObservedComponentDesiredDetails) error {
	var (
		containerName string
		err           error
	)
	p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.MatchLabels = observedComponentDetails.MatchLabels
	p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.PodTemplateLabels = observedComponentDetails.PodTemplateLabels
	if len(p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.ContainerName) > 0 {
		containerName = p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.ContainerName
	} else {
		containerName = types.AdmissionServerContainerKey
	}
	p.ObservedOpenEBS.Spec.CstorConfig.AdmissionServer.ENV, err = fetchExistingContainerEnvs(
		observedComponentDetails.Containers, containerName)
	if err != nil {
		return err
	}

	return nil
}

func (p *Planner) fillCStorCSINodeExistingValues(observedComponentDetails ObservedComponentDesiredDetails) error {
	var (
		containerName string
		err           error
		envs          []interface{}
	)
	p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.MatchLabels = observedComponentDetails.MatchLabels
	p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.PodTemplateLabels = observedComponentDetails.PodTemplateLabels
	if len(p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ContainerName) > 0 {
		containerName = p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ContainerName
	} else {
		containerName = ContainerCSTORCSIPluginName
	}
	envs, err = fetchExistingContainerEnvs(
		observedComponentDetails.Containers, containerName)
	if err != nil {
		return err
	}
	// If we are unable to fetch env for container cstor-csi-plugin then try to fetch for
	// openebs-csi-plugin.
	if envs == nil || len(envs) == 0 {
		envs, err = fetchExistingContainerEnvs(
			observedComponentDetails.Containers, ContainerOpenEBSCSIPluginName)
		if err != nil {
			return err
		}
	}
	// set the existing envs
	p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSINode.ENV = envs

	return nil
}

func (p *Planner) fillCStorCSIControllerExistingValues(observedComponentDetails ObservedComponentDesiredDetails) error {
	var (
		containerName string
		err           error
		envs          []interface{}
	)
	p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.MatchLabels = observedComponentDetails.MatchLabels
	p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.PodTemplateLabels = observedComponentDetails.PodTemplateLabels
	if len(p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.ContainerName) > 0 {
		containerName = p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.ContainerName
	} else {
		containerName = ContainerCSTORCSIPluginName
	}
	envs, err = fetchExistingContainerEnvs(
		observedComponentDetails.Containers, containerName)
	if err != nil {
		return err
	}
	// If we are unable to fetch env for container cstor-csi-plugin then try to fetch for
	// openebs-csi-plugin.
	if envs == nil || len(envs) == 0 {
		envs, err = fetchExistingContainerEnvs(
			observedComponentDetails.Containers, ContainerOpenEBSCSIPluginName)
		if err != nil {
			return err
		}
	}
	// set the existing envs
	p.ObservedOpenEBS.Spec.CstorConfig.CSI.CSIController.ENV = envs

	return nil
}
