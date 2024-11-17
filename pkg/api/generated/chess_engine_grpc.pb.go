// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: pkg/api/chess_engine.proto

package generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ChessEngine_CalculateBestMove_FullMethodName = "/ChessEngine/CalculateBestMove"
)

// ChessEngineClient is the client API for ChessEngine service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChessEngineClient interface {
	CalculateBestMove(ctx context.Context, in *MoveRequest, opts ...grpc.CallOption) (*MoveResponse, error)
}

type chessEngineClient struct {
	cc grpc.ClientConnInterface
}

func NewChessEngineClient(cc grpc.ClientConnInterface) ChessEngineClient {
	return &chessEngineClient{cc}
}

func (c *chessEngineClient) CalculateBestMove(ctx context.Context, in *MoveRequest, opts ...grpc.CallOption) (*MoveResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MoveResponse)
	err := c.cc.Invoke(ctx, ChessEngine_CalculateBestMove_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChessEngineServer is the server API for ChessEngine service.
// All implementations must embed UnimplementedChessEngineServer
// for forward compatibility.
type ChessEngineServer interface {
	CalculateBestMove(context.Context, *MoveRequest) (*MoveResponse, error)
	mustEmbedUnimplementedChessEngineServer()
}

// UnimplementedChessEngineServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChessEngineServer struct{}

func (UnimplementedChessEngineServer) CalculateBestMove(context.Context, *MoveRequest) (*MoveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculateBestMove not implemented")
}
func (UnimplementedChessEngineServer) mustEmbedUnimplementedChessEngineServer() {}
func (UnimplementedChessEngineServer) testEmbeddedByValue()                     {}

// UnsafeChessEngineServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChessEngineServer will
// result in compilation errors.
type UnsafeChessEngineServer interface {
	mustEmbedUnimplementedChessEngineServer()
}

func RegisterChessEngineServer(s grpc.ServiceRegistrar, srv ChessEngineServer) {
	// If the following call pancis, it indicates UnimplementedChessEngineServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ChessEngine_ServiceDesc, srv)
}

func _ChessEngine_CalculateBestMove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChessEngineServer).CalculateBestMove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChessEngine_CalculateBestMove_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChessEngineServer).CalculateBestMove(ctx, req.(*MoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChessEngine_ServiceDesc is the grpc.ServiceDesc for ChessEngine service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChessEngine_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ChessEngine",
	HandlerType: (*ChessEngineServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CalculateBestMove",
			Handler:    _ChessEngine_CalculateBestMove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/chess_engine.proto",
}