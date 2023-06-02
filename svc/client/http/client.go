// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 5f7d5bf015
// Version Date: 2021-11-26T09:27:01Z

// Package http provides an HTTP client for the Swapsvc service.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gogo/protobuf/jsonpb"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"

	// This Service
	pb "github.com/mises-id/mises-swapsvc/proto"
	"github.com/mises-id/mises-swapsvc/svc"
)

var (
	_ = endpoint.Chain
	_ = httptransport.NewClient
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = ioutil.NopCloser
	_ = io.EOF
)

// New returns a service backed by an HTTP server living at the remote
// instance. We expect instance to come from a service discovery system, so
// likely of the form "host:port".
func New(instance string, options ...httptransport.ClientOption) (pb.SwapsvcServer, error) {

	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	_ = u

	var SyncSwapOrderZeroEndpoint endpoint.Endpoint
	{
		SyncSwapOrderZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/swap_order/sync/"),
			EncodeHTTPSyncSwapOrderZeroRequest,
			DecodeHTTPSyncSwapOrderResponse,
			options...,
		).Endpoint()
	}
	var SwapOrderPageZeroEndpoint endpoint.Endpoint
	{
		SwapOrderPageZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/swap_order/page/"),
			EncodeHTTPSwapOrderPageZeroRequest,
			DecodeHTTPSwapOrderPageResponse,
			options...,
		).Endpoint()
	}
	var FindSwapOrderZeroEndpoint endpoint.Endpoint
	{
		FindSwapOrderZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/swap_order/find/"),
			EncodeHTTPFindSwapOrderZeroRequest,
			DecodeHTTPFindSwapOrderResponse,
			options...,
		).Endpoint()
	}
	var ListSwapTokenZeroEndpoint endpoint.Endpoint
	{
		ListSwapTokenZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/swap_token/list/"),
			EncodeHTTPListSwapTokenZeroRequest,
			DecodeHTTPListSwapTokenResponse,
			options...,
		).Endpoint()
	}
	var GetSwapApproveAllowanceZeroEndpoint endpoint.Endpoint
	{
		GetSwapApproveAllowanceZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/swap/approve/allowance/"),
			EncodeHTTPGetSwapApproveAllowanceZeroRequest,
			DecodeHTTPGetSwapApproveAllowanceResponse,
			options...,
		).Endpoint()
	}
	var ApproveSwapTransactionZeroEndpoint endpoint.Endpoint
	{
		ApproveSwapTransactionZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/swap/approve/transaction/"),
			EncodeHTTPApproveSwapTransactionZeroRequest,
			DecodeHTTPApproveSwapTransactionResponse,
			options...,
		).Endpoint()
	}
	var SwapTradeZeroEndpoint endpoint.Endpoint
	{
		SwapTradeZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/swap/trade/"),
			EncodeHTTPSwapTradeZeroRequest,
			DecodeHTTPSwapTradeResponse,
			options...,
		).Endpoint()
	}
	var SwapQuoteZeroEndpoint endpoint.Endpoint
	{
		SwapQuoteZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/swap/quote/"),
			EncodeHTTPSwapQuoteZeroRequest,
			DecodeHTTPSwapQuoteResponse,
			options...,
		).Endpoint()
	}
	var TestZeroEndpoint endpoint.Endpoint
	{
		TestZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/test/"),
			EncodeHTTPTestZeroRequest,
			DecodeHTTPTestResponse,
			options...,
		).Endpoint()
	}
	var HealthZeroEndpoint endpoint.Endpoint
	{
		HealthZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/health/"),
			EncodeHTTPHealthZeroRequest,
			DecodeHTTPHealthResponse,
			options...,
		).Endpoint()
	}

	return svc.Endpoints{
		SyncSwapOrderEndpoint:           SyncSwapOrderZeroEndpoint,
		SwapOrderPageEndpoint:           SwapOrderPageZeroEndpoint,
		FindSwapOrderEndpoint:           FindSwapOrderZeroEndpoint,
		ListSwapTokenEndpoint:           ListSwapTokenZeroEndpoint,
		GetSwapApproveAllowanceEndpoint: GetSwapApproveAllowanceZeroEndpoint,
		ApproveSwapTransactionEndpoint:  ApproveSwapTransactionZeroEndpoint,
		SwapTradeEndpoint:               SwapTradeZeroEndpoint,
		SwapQuoteEndpoint:               SwapQuoteZeroEndpoint,
		TestEndpoint:                    TestZeroEndpoint,
		HealthEndpoint:                  HealthZeroEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

// CtxValuesToSend configures the http client to pull the specified keys out of
// the context and add them to the http request as headers.  Note that keys
// will have net/http.CanonicalHeaderKey called on them before being send over
// the wire and that is the form they will be available in the server context.
func CtxValuesToSend(keys ...string) httptransport.ClientOption {
	return httptransport.ClientBefore(func(ctx context.Context, r *http.Request) context.Context {
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				r.Header.Set(k, v)
			}
		}
		return ctx
	})
}

// HTTP Client Decode

// DecodeHTTPSyncSwapOrderResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded SyncSwapOrderResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPSyncSwapOrderResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.SyncSwapOrderResponse
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPSwapOrderPageResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded SwapOrderPageResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPSwapOrderPageResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.SwapOrderPageResponse
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPFindSwapOrderResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded FindSwapOrderResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPFindSwapOrderResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.FindSwapOrderResponse
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPListSwapTokenResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded ListSwapTokenResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPListSwapTokenResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.ListSwapTokenResponse
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPGetSwapApproveAllowanceResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded GetSwapApproveAllowanceResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPGetSwapApproveAllowanceResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.GetSwapApproveAllowanceResponse
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPApproveSwapTransactionResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded ApproveSwapTransactionResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPApproveSwapTransactionResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.ApproveSwapTransactionResponse
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPSwapTradeResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded SwapTradeResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPSwapTradeResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.SwapTradeResponse
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPSwapQuoteResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded SwapQuoteResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPSwapQuoteResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.SwapQuoteResponse
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPTestResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded TestResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPTestResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.TestResponse
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPHealthResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded HealthResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPHealthResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.HealthResponse
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// HTTP Client Encode

// EncodeHTTPSyncSwapOrderZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a syncswaporder request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPSyncSwapOrderZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.SyncSwapOrderRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap_order",
		"sync",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPSyncSwapOrderOneRequest is a transport/http.EncodeRequestFunc
// that encodes a syncswaporder request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPSyncSwapOrderOneRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.SyncSwapOrderRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap_order",
		"sync",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPSwapOrderPageZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a swaporderpage request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPSwapOrderPageZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.SwapOrderPageRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap_order",
		"page",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("from_address", fmt.Sprint(req.FromAddress))

	tmp, err = json.Marshal(req.Paginator)
	if err != nil {
		return errors.Wrap(err, "failed to marshal req.Paginator")
	}
	strval = string(tmp)
	values.Add("paginator", strval)

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPSwapOrderPageOneRequest is a transport/http.EncodeRequestFunc
// that encodes a swaporderpage request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPSwapOrderPageOneRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.SwapOrderPageRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap_order",
		"page",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("from_address", fmt.Sprint(req.FromAddress))

	tmp, err = json.Marshal(req.Paginator)
	if err != nil {
		return errors.Wrap(err, "failed to marshal req.Paginator")
	}
	strval = string(tmp)
	values.Add("paginator", strval)

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPFindSwapOrderZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a findswaporder request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPFindSwapOrderZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.FindSwapOrderRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap_order",
		"find",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("tx_hash", fmt.Sprint(req.TxHash))

	values.Add("from_address", fmt.Sprint(req.FromAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPFindSwapOrderOneRequest is a transport/http.EncodeRequestFunc
// that encodes a findswaporder request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPFindSwapOrderOneRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.FindSwapOrderRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap_order",
		"find",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("tx_hash", fmt.Sprint(req.TxHash))

	values.Add("from_address", fmt.Sprint(req.FromAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPListSwapTokenZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a listswaptoken request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPListSwapTokenZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.ListSwapTokenRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap_token",
		"list",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("token_address", fmt.Sprint(req.TokenAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPListSwapTokenOneRequest is a transport/http.EncodeRequestFunc
// that encodes a listswaptoken request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPListSwapTokenOneRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.ListSwapTokenRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap_token",
		"list",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("token_address", fmt.Sprint(req.TokenAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPGetSwapApproveAllowanceZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a getswapapproveallowance request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPGetSwapApproveAllowanceZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.GetSwapApproveAllowanceRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap",
		"approve",
		"allowance",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("token_address", fmt.Sprint(req.TokenAddress))

	values.Add("wallet_address", fmt.Sprint(req.WalletAddress))

	values.Add("aggregator_address", fmt.Sprint(req.AggregatorAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPGetSwapApproveAllowanceOneRequest is a transport/http.EncodeRequestFunc
// that encodes a getswapapproveallowance request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPGetSwapApproveAllowanceOneRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.GetSwapApproveAllowanceRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap",
		"approve",
		"allowance",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("token_address", fmt.Sprint(req.TokenAddress))

	values.Add("wallet_address", fmt.Sprint(req.WalletAddress))

	values.Add("aggregator_address", fmt.Sprint(req.AggregatorAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPApproveSwapTransactionZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a approveswaptransaction request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPApproveSwapTransactionZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.ApproveSwapTransactionRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap",
		"approve",
		"transaction",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("token_address", fmt.Sprint(req.TokenAddress))

	values.Add("amount", fmt.Sprint(req.Amount))

	values.Add("aggregator_address", fmt.Sprint(req.AggregatorAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPApproveSwapTransactionOneRequest is a transport/http.EncodeRequestFunc
// that encodes a approveswaptransaction request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPApproveSwapTransactionOneRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.ApproveSwapTransactionRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap",
		"approve",
		"transaction",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("token_address", fmt.Sprint(req.TokenAddress))

	values.Add("amount", fmt.Sprint(req.Amount))

	values.Add("aggregator_address", fmt.Sprint(req.AggregatorAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPSwapTradeZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a swaptrade request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPSwapTradeZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.SwapTradeRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap",
		"trade",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("amount", fmt.Sprint(req.Amount))

	values.Add("from_token_address", fmt.Sprint(req.FromTokenAddress))

	values.Add("to_token_address", fmt.Sprint(req.ToTokenAddress))

	values.Add("slippage", fmt.Sprint(req.Slippage))

	values.Add("from_address", fmt.Sprint(req.FromAddress))

	values.Add("dest_receiver", fmt.Sprint(req.DestReceiver))

	values.Add("aggregator_address", fmt.Sprint(req.AggregatorAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPSwapTradeOneRequest is a transport/http.EncodeRequestFunc
// that encodes a swaptrade request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPSwapTradeOneRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.SwapTradeRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap",
		"trade",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("amount", fmt.Sprint(req.Amount))

	values.Add("from_token_address", fmt.Sprint(req.FromTokenAddress))

	values.Add("to_token_address", fmt.Sprint(req.ToTokenAddress))

	values.Add("slippage", fmt.Sprint(req.Slippage))

	values.Add("from_address", fmt.Sprint(req.FromAddress))

	values.Add("dest_receiver", fmt.Sprint(req.DestReceiver))

	values.Add("aggregator_address", fmt.Sprint(req.AggregatorAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPSwapQuoteZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a swapquote request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPSwapQuoteZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.SwapQuoteRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap",
		"quote",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("amount", fmt.Sprint(req.Amount))

	values.Add("from_token_address", fmt.Sprint(req.FromTokenAddress))

	values.Add("to_token_address", fmt.Sprint(req.ToTokenAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPSwapQuoteOneRequest is a transport/http.EncodeRequestFunc
// that encodes a swapquote request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPSwapQuoteOneRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.SwapQuoteRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"swap",
		"quote",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("chainID", fmt.Sprint(req.ChainID))

	values.Add("amount", fmt.Sprint(req.Amount))

	values.Add("from_token_address", fmt.Sprint(req.FromTokenAddress))

	values.Add("to_token_address", fmt.Sprint(req.ToTokenAddress))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPTestZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a test request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPTestZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.TestRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"test",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("type", fmt.Sprint(req.Type))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPTestOneRequest is a transport/http.EncodeRequestFunc
// that encodes a test request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPTestOneRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.TestRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"test",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("type", fmt.Sprint(req.Type))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPHealthZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a health request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPHealthZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.HealthRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"health",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("type", fmt.Sprint(req.Type))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPHealthOneRequest is a transport/http.EncodeRequestFunc
// that encodes a health request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPHealthOneRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.HealthRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"health",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("type", fmt.Sprint(req.Type))

	r.URL.RawQuery = values.Encode()
	return nil
}

func errorDecoder(buf []byte) error {
	var w errorWrapper
	if err := json.Unmarshal(buf, &w); err != nil {
		const size = 8196
		if len(buf) > size {
			buf = buf[:size]
		}
		return fmt.Errorf("response body '%s': cannot parse non-json request body", buf)
	}

	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}
