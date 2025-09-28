// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package compare

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/device"
	"github.com/stretchr/testify/require"
)

func DeviceProfile(t testing.TB, req, resp device.Profile, msg ...string) {
	msg = addMsg(msg, "Comparing Rack Type")

	require.Equal(t, req.Selector, resp.Selector, msg)
	require.Equal(t, req.DeviceProfileType, resp.DeviceProfileType, msg)
	require.Equal(t, req.DualRoutingEngine, resp.DualRoutingEngine, msg)
	require.Equal(t, req.SoftwareCapabilities, resp.SoftwareCapabilities, msg)
	require.Equal(t, req.ReferenceDesignCapabilities, resp.ReferenceDesignCapabilities, msg)
	require.Equal(t, req.ChassisProfileID, resp.ChassisProfileID, msg)
	HardwareCapabilities(t, req.HardwareCapabilities, resp.HardwareCapabilities, msg...)
	require.Equal(t, req.Predefined, resp.Predefined, msg)
	require.Equal(t, req.SlotCount, resp.SlotCount, msg)
	require.Equal(t, req.ChassisCount, resp.ChassisCount, msg)
	require.Equal(t, len(req.Ports), len(resp.Ports), msg)
	for i := range len(req.Ports) {
		Port(t, req.Ports[i], resp.Ports[i], addMsg(msg, "Comparing Port %d", i)...)
	}
	require.Equal(t, req.Label, resp.Label, msg)
	if req.ChassisInfo == nil {
		require.Nil(t, resp.ChassisInfo, msg)
	} else {
		require.NotNil(t, resp.ChassisInfo, msg)
		ChassisInfo(t, *req.ChassisInfo, *resp.ChassisInfo, msg...)
	}
	require.Equal(t, len(req.LinecardsInfo), len(resp.LinecardsInfo), msg)
	for i := range len(req.LinecardsInfo) {
		LinecardInfo(t, req.LinecardsInfo[i], resp.LinecardsInfo[i], addMsg(msg, "Comparing Linecard %d", i)...)
	}
	require.Equal(t, len(req.SlotConfiguration), len(resp.SlotConfiguration), msg)
	for i := range len(req.SlotConfiguration) {
		require.Equal(t, req.SlotConfiguration[i], resp.SlotConfiguration[i], addMsg(msg, "Comparing Slot %d Configuration", i))
	}
	require.Equal(t, req.PhysicalDevice, resp.PhysicalDevice, msg)
	if req.ID() != nil && resp.ID() != nil {
		require.Equal(t, req.ID(), resp.ID(), msg)
	}
	if req.CreatedAt() != nil && resp.CreatedAt() != nil {
		require.Equal(t, req.CreatedAt(), resp.CreatedAt(), msg)
	}
	if req.LastModifiedAt() != nil && resp.LastModifiedAt() != nil {
		require.Equal(t, req.LastModifiedAt(), resp.LastModifiedAt(), msg)
	}
}

func HardwareCapabilities(t testing.TB, req, resp device.HardwareCapabilities, msg ...string) {
	msg = addMsg(msg, "Comparing Hardware Capabilities")

	require.Equal(t, req.MaxL3Mtu, resp.MaxL3Mtu, msg)
	require.Equal(t, req.MaxL2Mtu, resp.MaxL2Mtu, msg)
	require.Equal(t, req.FormFactor, resp.FormFactor, msg)
	require.Equal(t, req.VTEPLimit, resp.VTEPLimit, msg)
	require.Equal(t, req.BFDSupported, resp.BFDSupported, msg)
	FeatureVersions(t, req.COPPStrict, resp.COPPStrict, addMsg(msg, "comparing COPPStrict feature support")...)
	require.Equal(t, req.ECMPLimit, resp.ECMPLimit, msg)
	FeatureVersions(t, req.ASNSequencing, resp.ASNSequencing, addMsg(msg, "comparing ASNSequencing feature support")...)
	require.Equal(t, req.RAM, resp.RAM, msg)
	require.Equal(t, req.VTEPFloodLimit, resp.VTEPFloodLimit, msg)
	require.Equal(t, req.BreakoutCapable, resp.BreakoutCapable, msg)
	require.Equal(t, req.Userland, resp.Userland, msg)
	require.Equal(t, req.ASIC, resp.ASIC, msg)
	require.Equal(t, req.VRFLimit, resp.VRFLimit, msg)
	FeatureVersions(t, req.RoutingInstance, resp.RoutingInstance, addMsg(msg, "comparing RoutingInstance feature support")...)
	require.Equal(t, req.VxlanSupported, resp.VxlanSupported, msg)
	require.Equal(t, req.CPU, resp.CPU, msg)
}

func Port(t testing.TB, req, resp device.Port, msg ...string) {
	msg = addMsg(msg, "Comparing Port")

	require.Equal(t, req.ConnectorType, resp.ConnectorType, msg)
	require.Equal(t, req.Panel, resp.Panel, msg)
	require.Equal(t, len(req.Transformations), len(resp.Transformations), msg)
	for i := range len(req.Transformations) {
		Transformation(t, req.Transformations[i], resp.Transformations[i], addMsg(msg, "Comparing Transformation %d", i)...)
	}
	require.Equal(t, req.Column, resp.Column, msg)
	require.Equal(t, req.ID, resp.ID, msg)
	require.Equal(t, req.Row, resp.Row, msg)
	require.Equal(t, req.FailureDomain, resp.FailureDomain, msg)
	require.Equal(t, req.Display, resp.Display, msg)
	require.Equal(t, req.Slot, resp.Slot, msg)
}

func Transformation(t testing.TB, req, resp device.Transformation, msg ...string) {
	msg = addMsg(msg, "Comparing Transformation")

	require.Equal(t, req.ID, resp.ID, msg)
	require.Equal(t, req.IsDefault, resp.IsDefault, msg)
	require.Equal(t, len(req.Interfaces), len(resp.Interfaces), msg)
	for i := range len(req.Interfaces) {
		Interface(t, req.Interfaces[i], resp.Interfaces[i], addMsg(msg, "Comparing Interface %d", i)...)
	}
}

func Interface(t testing.TB, req, resp device.TransformationInterface, msg ...string) {
	msg = addMsg(msg, "Comparing Interface")

	require.Equal(t, req.ID, resp.ID, msg)
	require.Equal(t, req.Name, resp.Name, msg)
	require.Equal(t, req.State, resp.State, msg)
	require.Equal(t, req.Setting, resp.Setting, msg)
	require.Equal(t, req.Speed, resp.Speed, msg)
}

func FeatureVersions(t testing.TB, req, resp device.FeatureVersions, msg ...string) {
	msg = addMsg(msg, "Comparing Feature Versions")

	require.NoError(t, req.Validate(), "validating req")
	require.NoError(t, req.Validate(), "validating resp")
	require.Equal(t, len(req), len(resp), addMsg(msg, "different size slices"))
}

func ChassisInfo(t testing.TB, req, resp device.ProfileChassisInfo, msg ...string) {
	msg = addMsg(msg, "Comparing Chassis Info")

	require.Equal(t, req.ID, resp.ID, msg)
	require.Equal(t, req.Selector, resp.Selector, msg)
	HardwareCapabilities(t, req.HardwareCapabilities, resp.HardwareCapabilities, msg...)
	require.Equal(t, req.SoftwareCapabilities, resp.SoftwareCapabilities, msg)
	require.Equal(t, req.DualRoutingEngine, resp.DualRoutingEngine, msg)
	require.Equal(t, req.LinecardSlotIDs, resp.LinecardSlotIDs, msg)
	require.Equal(t, req.PhysicalDevice, resp.PhysicalDevice, msg)
	require.Equal(t, req.ReferenceDesignCapabilities, resp.ReferenceDesignCapabilities, msg)
}

func LinecardInfo(t testing.TB, req, resp device.ProfileLinecardInfo, msg ...string) {
	msg = addMsg(msg, "Comparing Linecard Info")

	require.Equal(t, req.ID, resp.ID, msg)
	require.Equal(t, req.Selector, resp.Selector, msg)
	HardwareCapabilities(t, req.HardwareCapabilities, resp.HardwareCapabilities, msg...)
}
