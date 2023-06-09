// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 5f7d5bf015
// Version Date: 2021-11-26T09:27:01Z

// Package grpc provides a gRPC client for the Swapsvc service.
package grpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/mises-id/mises-swapsvc/proto"
	"github.com/mises-id/mises-swapsvc/svc"
)

// New returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn, options ...ClientOption) (pb.SwapsvcServer, error) {
	var cc clientConfig

	for _, f := range options {
		err := f(&cc)
		if err != nil {
			return nil, errors.Wrap(err, "cannot apply option")
		}
	}

	clientOptions := []grpctransport.ClientOption{
		grpctransport.ClientBefore(
			contextValuesToGRPCMetadata(cc.headers)),
	}
	var syncswaporderEndpoint endpoint.Endpoint
	{
		syncswaporderEndpoint = grpctransport.NewClient(
			conn,
			"swapsvc.Swapsvc",
			"SyncSwapOrder",
			EncodeGRPCSyncSwapOrderRequest,
			DecodeGRPCSyncSwapOrderResponse,
			pb.SyncSwapOrderResponse{},
			clientOptions...,
		).Endpoint()
	}

	var swaporderpageEndpoint endpoint.Endpoint
	{
		swaporderpageEndpoint = grpctransport.NewClient(
			conn,
			"swapsvc.Swapsvc",
			"SwapOrderPage",
			EncodeGRPCSwapOrderPageRequest,
			DecodeGRPCSwapOrderPageResponse,
			pb.SwapOrderPageResponse{},
			clientOptions...,
		).Endpoint()
	}

	var findswaporderEndpoint endpoint.Endpoint
	{
		findswaporderEndpoint = grpctransport.NewClient(
			conn,
			"swapsvc.Swapsvc",
			"FindSwapOrder",
			EncodeGRPCFindSwapOrderRequest,
			DecodeGRPCFindSwapOrderResponse,
			pb.FindSwapOrderResponse{},
			clientOptions...,
		).Endpoint()
	}

	var listswaptokenEndpoint endpoint.Endpoint
	{
		listswaptokenEndpoint = grpctransport.NewClient(
			conn,
			"swapsvc.Swapsvc",
			"ListSwapToken",
			EncodeGRPCListSwapTokenRequest,
			DecodeGRPCListSwapTokenResponse,
			pb.ListSwapTokenResponse{},
			clientOptions...,
		).Endpoint()
	}

	var getswapapproveallowanceEndpoint endpoint.Endpoint
	{
		getswapapproveallowanceEndpoint = grpctransport.NewClient(
			conn,
			"swapsvc.Swapsvc",
			"GetSwapApproveAllowance",
			EncodeGRPCGetSwapApproveAllowanceRequest,
			DecodeGRPCGetSwapApproveAllowanceResponse,
			pb.GetSwapApproveAllowanceResponse{},
			clientOptions...,
		).Endpoint()
	}

	var approveswaptransactionEndpoint endpoint.Endpoint
	{
		approveswaptransactionEndpoint = grpctransport.NewClient(
			conn,
			"swapsvc.Swapsvc",
			"ApproveSwapTransaction",
			EncodeGRPCApproveSwapTransactionRequest,
			DecodeGRPCApproveSwapTransactionResponse,
			pb.ApproveSwapTransactionResponse{},
			clientOptions...,
		).Endpoint()
	}

	var swaptradeEndpoint endpoint.Endpoint
	{
		swaptradeEndpoint = grpctransport.NewClient(
			conn,
			"swapsvc.Swapsvc",
			"SwapTrade",
			EncodeGRPCSwapTradeRequest,
			DecodeGRPCSwapTradeResponse,
			pb.SwapTradeResponse{},
			clientOptions...,
		).Endpoint()
	}

	var swapquoteEndpoint endpoint.Endpoint
	{
		swapquoteEndpoint = grpctransport.NewClient(
			conn,
			"swapsvc.Swapsvc",
			"SwapQuote",
			EncodeGRPCSwapQuoteRequest,
			DecodeGRPCSwapQuoteResponse,
			pb.SwapQuoteResponse{},
			clientOptions...,
		).Endpoint()
	}

	var testEndpoint endpoint.Endpoint
	{
		testEndpoint = grpctransport.NewClient(
			conn,
			"swapsvc.Swapsvc",
			"Test",
			EncodeGRPCTestRequest,
			DecodeGRPCTestResponse,
			pb.TestResponse{},
			clientOptions...,
		).Endpoint()
	}

	var healthEndpoint endpoint.Endpoint
	{
		healthEndpoint = grpctransport.NewClient(
			conn,
			"swapsvc.Swapsvc",
			"Health",
			EncodeGRPCHealthRequest,
			DecodeGRPCHealthResponse,
			pb.HealthResponse{},
			clientOptions...,
		).Endpoint()
	}

	return svc.Endpoints{
		SyncSwapOrderEndpoint:           syncswaporderEndpoint,
		SwapOrderPageEndpoint:           swaporderpageEndpoint,
		FindSwapOrderEndpoint:           findswaporderEndpoint,
		ListSwapTokenEndpoint:           listswaptokenEndpoint,
		GetSwapApproveAllowanceEndpoint: getswapapproveallowanceEndpoint,
		ApproveSwapTransactionEndpoint:  approveswaptransactionEndpoint,
		SwapTradeEndpoint:               swaptradeEndpoint,
		SwapQuoteEndpoint:               swapquoteEndpoint,
		TestEndpoint:                    testEndpoint,
		HealthEndpoint:                  healthEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCSyncSwapOrderResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC syncswaporder reply to a user-domain syncswaporder response. Primarily useful in a client.
func DecodeGRPCSyncSwapOrderResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SyncSwapOrderResponse)
	return reply, nil
}

// DecodeGRPCSwapOrderPageResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC swaporderpage reply to a user-domain swaporderpage response. Primarily useful in a client.
func DecodeGRPCSwapOrderPageResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SwapOrderPageResponse)
	return reply, nil
}

// DecodeGRPCFindSwapOrderResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC findswaporder reply to a user-domain findswaporder response. Primarily useful in a client.
func DecodeGRPCFindSwapOrderResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.FindSwapOrderResponse)
	return reply, nil
}

// DecodeGRPCListSwapTokenResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC listswaptoken reply to a user-domain listswaptoken response. Primarily useful in a client.
func DecodeGRPCListSwapTokenResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ListSwapTokenResponse)
	return reply, nil
}

// DecodeGRPCGetSwapApproveAllowanceResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC getswapapproveallowance reply to a user-domain getswapapproveallowance response. Primarily useful in a client.
func DecodeGRPCGetSwapApproveAllowanceResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetSwapApproveAllowanceResponse)
	return reply, nil
}

// DecodeGRPCApproveSwapTransactionResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC approveswaptransaction reply to a user-domain approveswaptransaction response. Primarily useful in a client.
func DecodeGRPCApproveSwapTransactionResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ApproveSwapTransactionResponse)
	return reply, nil
}

// DecodeGRPCSwapTradeResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC swaptrade reply to a user-domain swaptrade response. Primarily useful in a client.
func DecodeGRPCSwapTradeResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SwapTradeResponse)
	return reply, nil
}

// DecodeGRPCSwapQuoteResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC swapquote reply to a user-domain swapquote response. Primarily useful in a client.
func DecodeGRPCSwapQuoteResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SwapQuoteResponse)
	return reply, nil
}

// DecodeGRPCTestResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC test reply to a user-domain test response. Primarily useful in a client.
func DecodeGRPCTestResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.TestResponse)
	return reply, nil
}

// DecodeGRPCHealthResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC health reply to a user-domain health response. Primarily useful in a client.
func DecodeGRPCHealthResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.HealthResponse)
	return reply, nil
}

// GRPC Client Encode

// EncodeGRPCSyncSwapOrderRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain syncswaporder request to a gRPC syncswaporder request. Primarily useful in a client.
func EncodeGRPCSyncSwapOrderRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.SyncSwapOrderRequest)
	return req, nil
}

// EncodeGRPCSwapOrderPageRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain swaporderpage request to a gRPC swaporderpage request. Primarily useful in a client.
func EncodeGRPCSwapOrderPageRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.SwapOrderPageRequest)
	return req, nil
}

// EncodeGRPCFindSwapOrderRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain findswaporder request to a gRPC findswaporder request. Primarily useful in a client.
func EncodeGRPCFindSwapOrderRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.FindSwapOrderRequest)
	return req, nil
}

// EncodeGRPCListSwapTokenRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain listswaptoken request to a gRPC listswaptoken request. Primarily useful in a client.
func EncodeGRPCListSwapTokenRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ListSwapTokenRequest)
	return req, nil
}

// EncodeGRPCGetSwapApproveAllowanceRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain getswapapproveallowance request to a gRPC getswapapproveallowance request. Primarily useful in a client.
func EncodeGRPCGetSwapApproveAllowanceRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetSwapApproveAllowanceRequest)
	return req, nil
}

// EncodeGRPCApproveSwapTransactionRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain approveswaptransaction request to a gRPC approveswaptransaction request. Primarily useful in a client.
func EncodeGRPCApproveSwapTransactionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ApproveSwapTransactionRequest)
	return req, nil
}

// EncodeGRPCSwapTradeRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain swaptrade request to a gRPC swaptrade request. Primarily useful in a client.
func EncodeGRPCSwapTradeRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.SwapTradeRequest)
	return req, nil
}

// EncodeGRPCSwapQuoteRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain swapquote request to a gRPC swapquote request. Primarily useful in a client.
func EncodeGRPCSwapQuoteRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.SwapQuoteRequest)
	return req, nil
}

// EncodeGRPCTestRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain test request to a gRPC test request. Primarily useful in a client.
func EncodeGRPCTestRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.TestRequest)
	return req, nil
}

// EncodeGRPCHealthRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain health request to a gRPC health request. Primarily useful in a client.
func EncodeGRPCHealthRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.HealthRequest)
	return req, nil
}

type clientConfig struct {
	headers []string
}

// ClientOption is a function that modifies the client config
type ClientOption func(*clientConfig) error

func CtxValuesToSend(keys ...string) ClientOption {
	return func(o *clientConfig) error {
		o.headers = keys
		return nil
	}
}

func contextValuesToGRPCMetadata(keys []string) grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		var pairs []string
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				pairs = append(pairs, k, v)
			}
		}

		if pairs != nil {
			*md = metadata.Join(*md, metadata.Pairs(pairs...))
		}

		return ctx
	}
}
