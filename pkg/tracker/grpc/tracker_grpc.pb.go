// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: tracker.proto

package grpc

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TrackerClient is the client API for Tracker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TrackerClient interface {
	RegisterTrackerInstance(ctx context.Context, in *TrackerInfo, opts ...grpc.CallOption) (*InstanceRegisterResponse, error)
	GetAllRegisteredInstances(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (Tracker_GetAllRegisteredInstancesClient, error)
	RegisterNewPackage(ctx context.Context, in *PackageInfo, opts ...grpc.CallOption) (*PackageRegisterResponse, error)
	FindPackageLocation(ctx context.Context, in *LocationRequest, opts ...grpc.CallOption) (*PackageLocation, error)
}

type trackerClient struct {
	cc grpc.ClientConnInterface
}

func NewTrackerClient(cc grpc.ClientConnInterface) TrackerClient {
	return &trackerClient{cc}
}

func (c *trackerClient) RegisterTrackerInstance(ctx context.Context, in *TrackerInfo, opts ...grpc.CallOption) (*InstanceRegisterResponse, error) {
	out := new(InstanceRegisterResponse)
	err := c.cc.Invoke(ctx, "/grpc.Tracker/RegisterTrackerInstance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackerClient) GetAllRegisteredInstances(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (Tracker_GetAllRegisteredInstancesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Tracker_ServiceDesc.Streams[0], "/grpc.Tracker/GetAllRegisteredInstances", opts...)
	if err != nil {
		return nil, err
	}
	x := &trackerGetAllRegisteredInstancesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Tracker_GetAllRegisteredInstancesClient interface {
	Recv() (*TrackerInfo, error)
	grpc.ClientStream
}

type trackerGetAllRegisteredInstancesClient struct {
	grpc.ClientStream
}

func (x *trackerGetAllRegisteredInstancesClient) Recv() (*TrackerInfo, error) {
	m := new(TrackerInfo)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *trackerClient) RegisterNewPackage(ctx context.Context, in *PackageInfo, opts ...grpc.CallOption) (*PackageRegisterResponse, error) {
	out := new(PackageRegisterResponse)
	err := c.cc.Invoke(ctx, "/grpc.Tracker/RegisterNewPackage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackerClient) FindPackageLocation(ctx context.Context, in *LocationRequest, opts ...grpc.CallOption) (*PackageLocation, error) {
	out := new(PackageLocation)
	err := c.cc.Invoke(ctx, "/grpc.Tracker/FindPackageLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrackerServer is the server API for Tracker service.
// All implementations must embed UnimplementedTrackerServer
// for forward compatibility
type TrackerServer interface {
	RegisterTrackerInstance(context.Context, *TrackerInfo) (*InstanceRegisterResponse, error)
	GetAllRegisteredInstances(*empty.Empty, Tracker_GetAllRegisteredInstancesServer) error
	RegisterNewPackage(context.Context, *PackageInfo) (*PackageRegisterResponse, error)
	FindPackageLocation(context.Context, *LocationRequest) (*PackageLocation, error)
	mustEmbedUnimplementedTrackerServer()
}

// UnimplementedTrackerServer must be embedded to have forward compatible implementations.
type UnimplementedTrackerServer struct {
}

func (UnimplementedTrackerServer) RegisterTrackerInstance(context.Context, *TrackerInfo) (*InstanceRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterTrackerInstance not implemented")
}
func (UnimplementedTrackerServer) GetAllRegisteredInstances(*empty.Empty, Tracker_GetAllRegisteredInstancesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAllRegisteredInstances not implemented")
}
func (UnimplementedTrackerServer) RegisterNewPackage(context.Context, *PackageInfo) (*PackageRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterNewPackage not implemented")
}
func (UnimplementedTrackerServer) FindPackageLocation(context.Context, *LocationRequest) (*PackageLocation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindPackageLocation not implemented")
}
func (UnimplementedTrackerServer) mustEmbedUnimplementedTrackerServer() {}

// UnsafeTrackerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TrackerServer will
// result in compilation errors.
type UnsafeTrackerServer interface {
	mustEmbedUnimplementedTrackerServer()
}

func RegisterTrackerServer(s grpc.ServiceRegistrar, srv TrackerServer) {
	s.RegisterService(&Tracker_ServiceDesc, srv)
}

func _Tracker_RegisterTrackerInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrackerInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackerServer).RegisterTrackerInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Tracker/RegisterTrackerInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackerServer).RegisterTrackerInstance(ctx, req.(*TrackerInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tracker_GetAllRegisteredInstances_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TrackerServer).GetAllRegisteredInstances(m, &trackerGetAllRegisteredInstancesServer{stream})
}

type Tracker_GetAllRegisteredInstancesServer interface {
	Send(*TrackerInfo) error
	grpc.ServerStream
}

type trackerGetAllRegisteredInstancesServer struct {
	grpc.ServerStream
}

func (x *trackerGetAllRegisteredInstancesServer) Send(m *TrackerInfo) error {
	return x.ServerStream.SendMsg(m)
}

func _Tracker_RegisterNewPackage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PackageInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackerServer).RegisterNewPackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Tracker/RegisterNewPackage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackerServer).RegisterNewPackage(ctx, req.(*PackageInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tracker_FindPackageLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackerServer).FindPackageLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Tracker/FindPackageLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackerServer).FindPackageLocation(ctx, req.(*LocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Tracker_ServiceDesc is the grpc.ServiceDesc for Tracker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tracker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Tracker",
	HandlerType: (*TrackerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterTrackerInstance",
			Handler:    _Tracker_RegisterTrackerInstance_Handler,
		},
		{
			MethodName: "RegisterNewPackage",
			Handler:    _Tracker_RegisterNewPackage_Handler,
		},
		{
			MethodName: "FindPackageLocation",
			Handler:    _Tracker_FindPackageLocation_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAllRegisteredInstances",
			Handler:       _Tracker_GetAllRegisteredInstances_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "tracker.proto",
}
