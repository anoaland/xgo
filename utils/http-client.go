package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	xgoErrors "github.com/anoaland/xgo/errors"
	"github.com/gofiber/fiber/v2"
)

type HttpClient struct {
	Url             string
	Method          string
	Headers         []HttpClientHeaders
	Payload         []byte
	Args            *fiber.Args
	ResponseSuccess interface{}
	ResponseError   interface{}
	ErrorPrefix     string
	LogRequest      bool
	LogResponse     bool
	RespHttpCode    int
}

type HttpClientHeaders struct {
	Key   string
	Value string
}

const JSON_CONTENT_TYPE = "application/json"

func (hc *HttpClient) Send() (interface{}, error) {
	// setup user agent
	client := fiber.AcquireAgent()

	// setup request
	clientReq := client.Request()
	clientReq.SetRequestURI(hc.Url)
	clientReq.Header.SetMethod(hc.Method)

	// setup custom headers
	if hc.Headers != nil {
		for _, header := range hc.Headers {
			clientReq.Header.Set(header.Key, header.Value)
		}
	} else {
		clientReq.Header.SetContentType(JSON_CONTENT_TYPE)
	}

	if hc.Payload != nil {
		clientReq.SetBody(hc.Payload)
	}

	if hc.Args != nil {
		client.Form(hc.Args)
		fiber.ReleaseArgs(hc.Args)
	}

	// log request
	if hc.LogRequest {
		fmt.Println(clientReq.String())
	}

	if err := client.Parse(); err != nil {
		return nil, err
	}

	// get response raw
	respCode, respBody, respErrs := client.Bytes()
	hc.RespHttpCode = respCode
	if respErrs != nil {
		err := errors.New(extractResponseErrors(respErrs))
		return nil, err
	}

	if respCode >= 300 {
		resError, err := resolveResponse(hc.ResponseError, respBody)

		if err != nil {
			return nil, xgoErrors.NewHttpError("❌ FAILED_TO_PARSE_RESPONSE_ERROR", err, 500, 2)
		}

		if hc.Payload != nil {
			fmt.Printf("❌ HTTP ERROR REQUEST PAYLOAD  %s", string(hc.Payload))
		}
		fmt.Printf("❌ HTTP ERROR RESPONSE [%d] %s", respCode, string(respBody))
		err = xgoErrors.NewHttpError("HTTP_CLIENT", errors.New(string(respBody)), respCode, 2)
		return resError, err
	}

	// log response
	if hc.LogResponse {
		fmt.Printf("response : [%d] %s", respCode, respBody)
	}

	return resolveResponse(hc.ResponseSuccess, respBody)

}

func resolveResponse(responseType interface{}, respBody []byte) (interface{}, error) {
	if responseType != nil {

		err := json.Unmarshal(respBody, &responseType)
		if err != nil {
			return nil, err
		}

		return responseType, nil

	}
	return respBody, nil
}

func extractResponseErrors(errors []error) string {
	var message string
	for i, err := range errors {
		if i > 0 {
			message += " | "
		}

		message += err.Error()
	}

	return message
}

func ContentTypeFormHeader() HttpClientHeaders {
	return HttpClientHeaders{
		Key:   "Content-Type",
		Value: "application/x-www-form-urlencoded",
	}
}

func AuthorizationHeader(token string) HttpClientHeaders {
	return HttpClientHeaders{
		Key:   "Authorization",
		Value: "Bearer " + token,
	}
}

func JsonContentTypeHeader() HttpClientHeaders {
	return HttpClientHeaders{
		Key:   "Content-Type",
		Value: JSON_CONTENT_TYPE,
	}
}
