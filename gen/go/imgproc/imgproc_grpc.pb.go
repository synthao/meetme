// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: imgproc/imgproc.proto

package imgproc_v1

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

// ImageProcessingServiceClient is the client API for ImageProcessingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ImageProcessingServiceClient interface {
	ProcessImage(ctx context.Context, in *ProcessImageRequest, opts ...grpc.CallOption) (*ProcessImageResponse, error)
}

type imageProcessingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewImageProcessingServiceClient(cc grpc.ClientConnInterface) ImageProcessingServiceClient {
	return &imageProcessingServiceClient{cc}
}

func (c *imageProcessingServiceClient) ProcessImage(ctx context.Context, in *ProcessImageRequest, opts ...grpc.CallOption) (*ProcessImageResponse, error) {
	out := new(ProcessImageResponse)
	err := c.cc.Invoke(ctx, "/imgproc.ImageProcessingService/ProcessImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImageProcessingServiceServer is the server API for ImageProcessingService service.
// All implementations must embed UnimplementedImageProcessingServiceServer
// for forward compatibility
type ImageProcessingServiceServer interface {
	ProcessImage(context.Context, *ProcessImageRequest) (*ProcessImageResponse, error)
	mustEmbedUnimplementedImageProcessingServiceServer()
}

// UnimplementedImageProcessingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedImageProcessingServiceServer struct {
}

func (UnimplementedImageProcessingServiceServer) ProcessImage(context.Context, *ProcessImageRequest) (*ProcessImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessImage not implemented")
}
func (UnimplementedImageProcessingServiceServer) mustEmbedUnimplementedImageProcessingServiceServer() {
}

// UnsafeImageProcessingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ImageProcessingServiceServer will
// result in compilation errors.
type UnsafeImageProcessingServiceServer interface {
	mustEmbedUnimplementedImageProcessingServiceServer()
}

func RegisterImageProcessingServiceServer(s grpc.ServiceRegistrar, srv ImageProcessingServiceServer) {
	s.RegisterService(&ImageProcessingService_ServiceDesc, srv)
}

func _ImageProcessingService_ProcessImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageProcessingServiceServer).ProcessImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imgproc.ImageProcessingService/ProcessImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageProcessingServiceServer).ProcessImage(ctx, req.(*ProcessImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ImageProcessingService_ServiceDesc is the grpc.ServiceDesc for ImageProcessingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ImageProcessingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "imgproc.ImageProcessingService",
	HandlerType: (*ImageProcessingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessImage",
			Handler:    _ImageProcessingService_ProcessImage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "imgproc/imgproc.proto",
}
