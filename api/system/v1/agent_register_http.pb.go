// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.0
// - protoc             v4.23.2
// source: api/system/v1/agent_register.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationAgentRegisterHeartbeat = "/api.system.v1.AgentRegister/Heartbeat"

type AgentRegisterHTTPServer interface {
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatReply, error)
}

func RegisterAgentRegisterHTTPServer(s *http.Server, srv AgentRegisterHTTPServer) {
	r := s.Route("/")
	r.GET("/health", _AgentRegister_Heartbeat0_HTTP_Handler(srv))
}

func _AgentRegister_Heartbeat0_HTTP_Handler(srv AgentRegisterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in HeartbeatRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAgentRegisterHeartbeat)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Heartbeat(ctx, req.(*HeartbeatRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*HeartbeatReply)
		return ctx.Result(200, reply)
	}
}

type AgentRegisterHTTPClient interface {
	Heartbeat(ctx context.Context, req *HeartbeatRequest, opts ...http.CallOption) (rsp *HeartbeatReply, err error)
}

type AgentRegisterHTTPClientImpl struct {
	cc *http.Client
}

func NewAgentRegisterHTTPClient(client *http.Client) AgentRegisterHTTPClient {
	return &AgentRegisterHTTPClientImpl{client}
}

func (c *AgentRegisterHTTPClientImpl) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...http.CallOption) (*HeartbeatReply, error) {
	var out HeartbeatReply
	pattern := "/health"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAgentRegisterHeartbeat))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
