// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.2
// source: api/compute/v1/vm.proto

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
	Vm_CreateVm_FullMethodName = "/api.compute.v1.Vm/CreateVm"
	Vm_DeleteVm_FullMethodName = "/api.compute.v1.Vm/DeleteVm"
	Vm_GetVm_FullMethodName    = "/api.compute.v1.Vm/GetVm"
	Vm_ListVm_FullMethodName   = "/api.compute.v1.Vm/ListVm"
	Vm_StartVm_FullMethodName  = "/api.compute.v1.Vm/StartVm"
	Vm_StopVm_FullMethodName   = "/api.compute.v1.Vm/StopVm"
)

// VmClient is the client API for Vm service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VmClient interface {
	CreateVm(ctx context.Context, in *CreateVmRequest, opts ...grpc.CallOption) (*GetVmReply, error)
	DeleteVm(ctx context.Context, in *DeleteVmRequest, opts ...grpc.CallOption) (*DeleteVmReply, error)
	GetVm(ctx context.Context, in *GetVmRequest, opts ...grpc.CallOption) (*GetVmReply, error)
	ListVm(ctx context.Context, in *ListVmRequest, opts ...grpc.CallOption) (*ListVmReply, error)
	StartVm(ctx context.Context, in *GetVmRequest, opts ...grpc.CallOption) (*GetVmReply, error)
	StopVm(ctx context.Context, in *GetVmRequest, opts ...grpc.CallOption) (*GetVmReply, error)
}

type vmClient struct {
	cc grpc.ClientConnInterface
}

func NewVmClient(cc grpc.ClientConnInterface) VmClient {
	return &vmClient{cc}
}

func (c *vmClient) CreateVm(ctx context.Context, in *CreateVmRequest, opts ...grpc.CallOption) (*GetVmReply, error) {
	out := new(GetVmReply)
	err := c.cc.Invoke(ctx, Vm_CreateVm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vmClient) DeleteVm(ctx context.Context, in *DeleteVmRequest, opts ...grpc.CallOption) (*DeleteVmReply, error) {
	out := new(DeleteVmReply)
	err := c.cc.Invoke(ctx, Vm_DeleteVm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vmClient) GetVm(ctx context.Context, in *GetVmRequest, opts ...grpc.CallOption) (*GetVmReply, error) {
	out := new(GetVmReply)
	err := c.cc.Invoke(ctx, Vm_GetVm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vmClient) ListVm(ctx context.Context, in *ListVmRequest, opts ...grpc.CallOption) (*ListVmReply, error) {
	out := new(ListVmReply)
	err := c.cc.Invoke(ctx, Vm_ListVm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vmClient) StartVm(ctx context.Context, in *GetVmRequest, opts ...grpc.CallOption) (*GetVmReply, error) {
	out := new(GetVmReply)
	err := c.cc.Invoke(ctx, Vm_StartVm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vmClient) StopVm(ctx context.Context, in *GetVmRequest, opts ...grpc.CallOption) (*GetVmReply, error) {
	out := new(GetVmReply)
	err := c.cc.Invoke(ctx, Vm_StopVm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VmServer is the server API for Vm service.
// All implementations must embed UnimplementedVmServer
// for forward compatibility
type VmServer interface {
	CreateVm(context.Context, *CreateVmRequest) (*GetVmReply, error)
	DeleteVm(context.Context, *DeleteVmRequest) (*DeleteVmReply, error)
	GetVm(context.Context, *GetVmRequest) (*GetVmReply, error)
	ListVm(context.Context, *ListVmRequest) (*ListVmReply, error)
	StartVm(context.Context, *GetVmRequest) (*GetVmReply, error)
	StopVm(context.Context, *GetVmRequest) (*GetVmReply, error)
	mustEmbedUnimplementedVmServer()
}

// UnimplementedVmServer must be embedded to have forward compatible implementations.
type UnimplementedVmServer struct {
}

func (UnimplementedVmServer) CreateVm(context.Context, *CreateVmRequest) (*GetVmReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateVm not implemented")
}
func (UnimplementedVmServer) DeleteVm(context.Context, *DeleteVmRequest) (*DeleteVmReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVm not implemented")
}
func (UnimplementedVmServer) GetVm(context.Context, *GetVmRequest) (*GetVmReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVm not implemented")
}
func (UnimplementedVmServer) ListVm(context.Context, *ListVmRequest) (*ListVmReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVm not implemented")
}
func (UnimplementedVmServer) StartVm(context.Context, *GetVmRequest) (*GetVmReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartVm not implemented")
}
func (UnimplementedVmServer) StopVm(context.Context, *GetVmRequest) (*GetVmReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopVm not implemented")
}
func (UnimplementedVmServer) mustEmbedUnimplementedVmServer() {}

// UnsafeVmServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VmServer will
// result in compilation errors.
type UnsafeVmServer interface {
	mustEmbedUnimplementedVmServer()
}

func RegisterVmServer(s grpc.ServiceRegistrar, srv VmServer) {
	s.RegisterService(&Vm_ServiceDesc, srv)
}

func _Vm_CreateVm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateVmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VmServer).CreateVm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vm_CreateVm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VmServer).CreateVm(ctx, req.(*CreateVmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vm_DeleteVm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteVmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VmServer).DeleteVm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vm_DeleteVm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VmServer).DeleteVm(ctx, req.(*DeleteVmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vm_GetVm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VmServer).GetVm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vm_GetVm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VmServer).GetVm(ctx, req.(*GetVmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vm_ListVm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListVmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VmServer).ListVm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vm_ListVm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VmServer).ListVm(ctx, req.(*ListVmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vm_StartVm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VmServer).StartVm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vm_StartVm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VmServer).StartVm(ctx, req.(*GetVmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vm_StopVm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VmServer).StopVm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vm_StopVm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VmServer).StopVm(ctx, req.(*GetVmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Vm_ServiceDesc is the grpc.ServiceDesc for Vm service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Vm_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.compute.v1.Vm",
	HandlerType: (*VmServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateVm",
			Handler:    _Vm_CreateVm_Handler,
		},
		{
			MethodName: "DeleteVm",
			Handler:    _Vm_DeleteVm_Handler,
		},
		{
			MethodName: "GetVm",
			Handler:    _Vm_GetVm_Handler,
		},
		{
			MethodName: "ListVm",
			Handler:    _Vm_ListVm_Handler,
		},
		{
			MethodName: "StartVm",
			Handler:    _Vm_StartVm_Handler,
		},
		{
			MethodName: "StopVm",
			Handler:    _Vm_StopVm_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/compute/v1/vm.proto",
}
