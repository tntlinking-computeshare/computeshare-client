// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.2
// source: api/compute/v1/computepower.proto

package v1

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

const (
	Computepower_RunPythonPackage_FullMethodName = "/api.compute.v1.Computepower/RunPythonPackage"
	Computepower_RunBenchmarks_FullMethodName    = "/api.compute.v1.Computepower/RunBenchmarks"
)

// ComputepowerClient is the client API for Computepower service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ComputepowerClient interface {
	RunPythonPackage(ctx context.Context, in *RunPythonPackageRequest, opts ...grpc.CallOption) (*RunPythonPackageReply, error)
	RunBenchmarks(ctx context.Context, in *RunBenchmarksRequest, opts ...grpc.CallOption) (*RunBenchmarksReply, error)
}

type computepowerClient struct {
	cc grpc.ClientConnInterface
}

func NewComputepowerClient(cc grpc.ClientConnInterface) ComputepowerClient {
	return &computepowerClient{cc}
}

func (c *computepowerClient) RunPythonPackage(ctx context.Context, in *RunPythonPackageRequest, opts ...grpc.CallOption) (*RunPythonPackageReply, error) {
	out := new(RunPythonPackageReply)
	err := c.cc.Invoke(ctx, Computepower_RunPythonPackage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *computepowerClient) RunBenchmarks(ctx context.Context, in *RunBenchmarksRequest, opts ...grpc.CallOption) (*RunBenchmarksReply, error) {
	out := new(RunBenchmarksReply)
	err := c.cc.Invoke(ctx, Computepower_RunBenchmarks_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ComputepowerServer is the server API for Computepower service.
// All implementations must embed UnimplementedComputepowerServer
// for forward compatibility
type ComputepowerServer interface {
	RunPythonPackage(context.Context, *RunPythonPackageRequest) (*RunPythonPackageReply, error)
	RunBenchmarks(context.Context, *RunBenchmarksRequest) (*RunBenchmarksReply, error)
	mustEmbedUnimplementedComputepowerServer()
}

// UnimplementedComputepowerServer must be embedded to have forward compatible implementations.
type UnimplementedComputepowerServer struct {
}

func (UnimplementedComputepowerServer) RunPythonPackage(context.Context, *RunPythonPackageRequest) (*RunPythonPackageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunPythonPackage not implemented")
}
func (UnimplementedComputepowerServer) RunBenchmarks(context.Context, *RunBenchmarksRequest) (*RunBenchmarksReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunBenchmarks not implemented")
}
func (UnimplementedComputepowerServer) mustEmbedUnimplementedComputepowerServer() {}

// UnsafeComputepowerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ComputepowerServer will
// result in compilation errors.
type UnsafeComputepowerServer interface {
	mustEmbedUnimplementedComputepowerServer()
}

func RegisterComputepowerServer(s grpc.ServiceRegistrar, srv ComputepowerServer) {
	s.RegisterService(&Computepower_ServiceDesc, srv)
}

func _Computepower_RunPythonPackage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunPythonPackageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ComputepowerServer).RunPythonPackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Computepower_RunPythonPackage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ComputepowerServer).RunPythonPackage(ctx, req.(*RunPythonPackageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Computepower_RunBenchmarks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunBenchmarksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ComputepowerServer).RunBenchmarks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Computepower_RunBenchmarks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ComputepowerServer).RunBenchmarks(ctx, req.(*RunBenchmarksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Computepower_ServiceDesc is the grpc.ServiceDesc for Computepower service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Computepower_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.compute.v1.Computepower",
	HandlerType: (*ComputepowerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunPythonPackage",
			Handler:    _Computepower_RunPythonPackage_Handler,
		},
		{
			MethodName: "RunBenchmarks",
			Handler:    _Computepower_RunBenchmarks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/compute/v1/computepower.proto",
}