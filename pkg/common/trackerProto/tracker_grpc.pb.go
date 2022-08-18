// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: tracker.proto

package trackerProto

import (
	context "context"
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
	RegisterNewPackage(ctx context.Context, in *PackageInfo, opts ...grpc.CallOption) (*RegisterResponse, error)
	FindPackageLocation(ctx context.Context, in *LocationRequest, opts ...grpc.CallOption) (*PackageLocation, error)
}

type trackerClient struct {
	cc grpc.ClientConnInterface
}

func NewTrackerClient(cc grpc.ClientConnInterface) TrackerClient {
	return &trackerClient{cc}
}

func (c *trackerClient) RegisterNewPackage(ctx context.Context, in *PackageInfo, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/trackerProto.Tracker/RegisterNewPackage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackerClient) FindPackageLocation(ctx context.Context, in *LocationRequest, opts ...grpc.CallOption) (*PackageLocation, error) {
	out := new(PackageLocation)
	err := c.cc.Invoke(ctx, "/trackerProto.Tracker/FindPackageLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrackerServer is the server API for Tracker service.
// All implementations must embed UnimplementedTrackerServer
// for forward compatibility
type TrackerServer interface {
	RegisterNewPackage(context.Context, *PackageInfo) (*RegisterResponse, error)
	FindPackageLocation(context.Context, *LocationRequest) (*PackageLocation, error)
	mustEmbedUnimplementedTrackerServer()
}

// UnimplementedTrackerServer must be embedded to have forward compatible implementations.
type UnimplementedTrackerServer struct {
}

func (UnimplementedTrackerServer) RegisterNewPackage(context.Context, *PackageInfo) (*RegisterResponse, error) {
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
		FullMethod: "/trackerProto.Tracker/RegisterNewPackage",
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
		FullMethod: "/trackerProto.Tracker/FindPackageLocation",
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
	ServiceName: "trackerProto.Tracker",
	HandlerType: (*TrackerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterNewPackage",
			Handler:    _Tracker_RegisterNewPackage_Handler,
		},
		{
			MethodName: "FindPackageLocation",
			Handler:    _Tracker_FindPackageLocation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tracker.proto",
}
