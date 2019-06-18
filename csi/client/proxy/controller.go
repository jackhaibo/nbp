// Copyright 2018 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxy

import (
	"github.com/golang/glog"

	csi "github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/opensds/nbp/csi/util"
	"golang.org/x/net/context"
)

// Controller define
type Controller struct {
	client csi.ControllerClient
}

////////////////////////////////////////////////////////////////////////////////
//                            Controller Client Init                          //
////////////////////////////////////////////////////////////////////////////////

// GetController return controller struct
func GetController(csiEndpoint string) (client Controller, err error) {
	conn, err := util.GetCSIClientConn(csiEndpoint)
	if err != nil {
		glog.Fatalf("did not connect: %v", err)
	}

	c := csi.NewControllerClient(conn)

	return Controller{client: c}, nil
}

////////////////////////////////////////////////////////////////////////////////
//                            Controller Client Proxy                         //
////////////////////////////////////////////////////////////////////////////////

// CreateVolume proxy
func (c *Controller) CreateVolume(
	ctx context.Context,
	name string,
	capacity *csi.CapacityRange, /*Optional*/
	capabilities []*csi.VolumeCapability,
	params map[string]string, /*Optional*/
	credentials map[string]string /*Optional*/) (volume *csi.Volume, err error) {

	req := &csi.CreateVolumeRequest{
		Name:               name,
		CapacityRange:      capacity,
		VolumeCapabilities: capabilities,
		Parameters:         params,
		Secrets:            credentials,
	}

	rs, err := c.client.CreateVolume(ctx, req)
	if err != nil {
		return nil, err
	}

	return rs.Volume, nil
}

// DeleteVolume proxy
func (c *Controller) DeleteVolume(
	ctx context.Context,
	volumeid string,
	credentials map[string]string /*Optional*/) error {

	req := &csi.DeleteVolumeRequest{
		VolumeId: volumeid,
		Secrets:  credentials,
	}

	_, err := c.client.DeleteVolume(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// ControllerPublishVolume proxy
func (c *Controller) ControllerPublishVolume(
	ctx context.Context,
	volumeid string,
	nodeID string,
	capabilities *csi.VolumeCapability,
	readonly bool,
	credentials map[string]string, /*Optional*/
	volumeattributes map[string]string /*Optional*/) (map[string]string, error) {

	req := &csi.ControllerPublishVolumeRequest{
		VolumeId:         volumeid,
		NodeId:           nodeID,
		VolumeCapability: capabilities,
		Readonly:         readonly,
		Secrets:          credentials,
		VolumeContext:    volumeattributes,
	}

	rs, err := c.client.ControllerPublishVolume(ctx, req)
	if err != nil {
		return nil, err
	}

	return rs.PublishContext, nil
}

// ControllerUnpublishVolume proxy
func (c *Controller) ControllerUnpublishVolume(
	ctx context.Context,
	volumeid string,
	nodeID string, /*Optional*/
	credentials map[string]string /*Optional*/) error {

	req := &csi.ControllerUnpublishVolumeRequest{
		VolumeId: volumeid,
		NodeId:   nodeID,
		Secrets:  credentials,
	}

	_, err := c.client.ControllerUnpublishVolume(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// ValidateVolumeCapabilities proxy
func (c *Controller) ValidateVolumeCapabilities(
	ctx context.Context,
	volumeid string,
	capabilities []*csi.VolumeCapability,
	volumeattributes map[string]string /*Optional*/) (*csi.ValidateVolumeCapabilitiesResponse, error) {

	req := &csi.ValidateVolumeCapabilitiesRequest{
		VolumeId:           volumeid,
		VolumeCapabilities: capabilities,
		VolumeContext:      volumeattributes,
	}

	rs, err := c.client.ValidateVolumeCapabilities(ctx, req)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

// ListVolumes proxy
func (c *Controller) ListVolumes(
	ctx context.Context,
	maxentries int32, /*Optional*/
	startingtoken string /*Optional*/) (entries []*csi.ListVolumesResponse_Entry, nextToken string, err error) {

	req := &csi.ListVolumesRequest{
		MaxEntries:    maxentries,
		StartingToken: startingtoken,
	}

	rs, err := c.client.ListVolumes(ctx, req)
	if err != nil {
		return nil, "", err
	}

	return rs.Entries, rs.NextToken, nil
}

// GetCapacity proxy
func (c *Controller) GetCapacity(
	ctx context.Context,
	capabilities []*csi.VolumeCapability /*Optional*/) (int64, error) {

	req := &csi.GetCapacityRequest{
		VolumeCapabilities: capabilities,
	}

	rs, err := c.client.GetCapacity(ctx, req)
	if err != nil {
		return 0, err
	}

	return rs.AvailableCapacity, nil
}

// ControllerGetCapabilities proxy
func (c *Controller) ControllerGetCapabilities(
	ctx context.Context) (capabilties []*csi.ControllerServiceCapability, err error) {

	req := &csi.ControllerGetCapabilitiesRequest{}

	rs, err := c.client.ControllerGetCapabilities(ctx, req)
	if err != nil {
		return nil, err
	}

	return rs.Capabilities, nil
}
