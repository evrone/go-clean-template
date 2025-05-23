// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// TranslationClient is the client API for Translation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TranslationClient interface {
	// RPC method to get translation history.
	GetHistory(ctx context.Context, in *GetHistoryRequest, opts ...grpc.CallOption) (*GetHistoryResponse, error)
}

type translationClient struct {
	cc grpc.ClientConnInterface
}

func NewTranslationClient(cc grpc.ClientConnInterface) TranslationClient {
	return &translationClient{cc}
}

func (c *translationClient) GetHistory(ctx context.Context, in *GetHistoryRequest, opts ...grpc.CallOption) (*GetHistoryResponse, error) {
	out := new(GetHistoryResponse)
	err := c.cc.Invoke(ctx, "/grpc.v1.Translation/GetHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TranslationServer is the server API for Translation service.
// All implementations must embed UnimplementedTranslationServer
// for forward compatibility
type TranslationServer interface {
	// RPC method to get translation history.
	GetHistory(context.Context, *GetHistoryRequest) (*GetHistoryResponse, error)
	mustEmbedUnimplementedTranslationServer()
}

// UnimplementedTranslationServer must be embedded to have forward compatible implementations.
type UnimplementedTranslationServer struct {
}

func (UnimplementedTranslationServer) GetHistory(context.Context, *GetHistoryRequest) (*GetHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHistory not implemented")
}
func (UnimplementedTranslationServer) mustEmbedUnimplementedTranslationServer() {}

// UnsafeTranslationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TranslationServer will
// result in compilation errors.
type UnsafeTranslationServer interface {
	mustEmbedUnimplementedTranslationServer()
}

func RegisterTranslationServer(s grpc.ServiceRegistrar, srv TranslationServer) {
	s.RegisterService(&Translation_ServiceDesc, srv)
}

func _Translation_GetHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TranslationServer).GetHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.v1.Translation/GetHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TranslationServer).GetHistory(ctx, req.(*GetHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Translation_ServiceDesc is the grpc.ServiceDesc for Translation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Translation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.v1.Translation",
	HandlerType: (*TranslationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetHistory",
			Handler:    _Translation_GetHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "docs/proto/v1/translation.history.proto",
}
