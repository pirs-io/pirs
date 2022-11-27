// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: process-storage.proto

package grpc

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

// StorageClient is the client API for Storage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageClient interface {
	UploadProcess(ctx context.Context, opts ...grpc.CallOption) (Storage_UploadProcessClient, error)
	DownloadProcess(ctx context.Context, in *ProcessDownloadRequest, opts ...grpc.CallOption) (Storage_DownloadProcessClient, error)
}

type storageClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageClient(cc grpc.ClientConnInterface) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) UploadProcess(ctx context.Context, opts ...grpc.CallOption) (Storage_UploadProcessClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[0], "/grpc.Storage/UploadProcess", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageUploadProcessClient{stream}
	return x, nil
}

type Storage_UploadProcessClient interface {
	Send(*ProcessUploadRequest) error
	Recv() (*ProcessUploadResponse, error)
	grpc.ClientStream
}

type storageUploadProcessClient struct {
	grpc.ClientStream
}

func (x *storageUploadProcessClient) Send(m *ProcessUploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storageUploadProcessClient) Recv() (*ProcessUploadResponse, error) {
	m := new(ProcessUploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageClient) DownloadProcess(ctx context.Context, in *ProcessDownloadRequest, opts ...grpc.CallOption) (Storage_DownloadProcessClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[1], "/grpc.Storage/DownloadProcess", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageDownloadProcessClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Storage_DownloadProcessClient interface {
	Recv() (*ProcessDownloadResponse, error)
	grpc.ClientStream
}

type storageDownloadProcessClient struct {
	grpc.ClientStream
}

func (x *storageDownloadProcessClient) Recv() (*ProcessDownloadResponse, error) {
	m := new(ProcessDownloadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StorageServer is the server API for Storage service.
// All implementations must embed UnimplementedStorageServer
// for forward compatibility
type StorageServer interface {
	UploadProcess(Storage_UploadProcessServer) error
	DownloadProcess(*ProcessDownloadRequest, Storage_DownloadProcessServer) error
	mustEmbedUnimplementedStorageServer()
}

// UnimplementedStorageServer must be embedded to have forward compatible implementations.
type UnimplementedStorageServer struct {
}

func (UnimplementedStorageServer) UploadProcess(Storage_UploadProcessServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadProcess not implemented")
}
func (UnimplementedStorageServer) DownloadProcess(*ProcessDownloadRequest, Storage_DownloadProcessServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadProcess not implemented")
}
func (UnimplementedStorageServer) mustEmbedUnimplementedStorageServer() {}

// UnsafeStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServer will
// result in compilation errors.
type UnsafeStorageServer interface {
	mustEmbedUnimplementedStorageServer()
}

func RegisterStorageServer(s grpc.ServiceRegistrar, srv StorageServer) {
	s.RegisterService(&Storage_ServiceDesc, srv)
}

func _Storage_UploadProcess_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServer).UploadProcess(&storageUploadProcessServer{stream})
}

type Storage_UploadProcessServer interface {
	Send(*ProcessUploadResponse) error
	Recv() (*ProcessUploadRequest, error)
	grpc.ServerStream
}

type storageUploadProcessServer struct {
	grpc.ServerStream
}

func (x *storageUploadProcessServer) Send(m *ProcessUploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storageUploadProcessServer) Recv() (*ProcessUploadRequest, error) {
	m := new(ProcessUploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Storage_DownloadProcess_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ProcessDownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServer).DownloadProcess(m, &storageDownloadProcessServer{stream})
}

type Storage_DownloadProcessServer interface {
	Send(*ProcessDownloadResponse) error
	grpc.ServerStream
}

type storageDownloadProcessServer struct {
	grpc.ServerStream
}

func (x *storageDownloadProcessServer) Send(m *ProcessDownloadResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Storage_ServiceDesc is the grpc.ServiceDesc for Storage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Storage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadProcess",
			Handler:       _Storage_UploadProcess_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadProcess",
			Handler:       _Storage_DownloadProcess_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "process-storage.proto",
}
