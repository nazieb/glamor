package glamor

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo"
)

type lambdaFn func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func WrapServer(e *echo.Echo) lambdaFn {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		body := strings.NewReader(request.Body)
		req := httptest.NewRequest(request.HTTPMethod, request.Path, body)
		for k, v := range request.Headers {
			req.Header.Add(k, v)
		}

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := rec.Result()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
				Headers:    headersToMap(res.Header),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: res.StatusCode,
			Body:       string(b),
			Headers:    headersToMap(res.Header),
		}, nil
	}
}

func headersToMap(headers http.Header) map[string]string {
	result := make(map[string]string)

	for key, values := range headers {
		var resultValue string

		if len(values) == 0 {
			resultValue = ""
		} else {
			resultValue = values[0]
		}

		result[key] = resultValue
	}

	return result
}
