package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"net"
	"strings"

	"github.com/SGNL-ai/TraTs-Demo-Svcs/gateway/pkg/common"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"go.uber.org/zap"
)

const (
	OIDC_ID_TOKEN_TYPE = "urn:ietf:params:oauth:token-type:id_token"
	TXN_TOKEN_TYPE     = "urn:ietf:params:oauth:token-type:txn_token"
	GRANT_TYPE         = "urn:ietf:params:oauth:grant-type:token-exchange"
)

const (
	AUDIENCE = "https://alphastocks.com/"
)

type RequestDetails struct {
	Path            string            `json:"path"`
	Method          common.HttpMethod `json:"method"`
	QueryParameters json.RawMessage   `json:"queryParameters"`
	Headers         json.RawMessage   `json:"headers"`
	Body            json.RawMessage   `json:"body"`
}

type txnToken struct {
	IssuedTokenType string `json:"issued_token_type"`
	AccessToken     string `json:"access_token"`
}

func getRequesterIP(r *http.Request) (string, error) {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		parts := strings.Split(xForwardedFor, ",")
		for _, part := range parts {
			ip := strings.TrimSpace(part)
			if validIP := net.ParseIP(ip); validIP != nil {
				return ip, nil
			}
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		if validIP := net.ParseIP(r.RemoteAddr); validIP != nil {
			return r.RemoteAddr, nil
		}

		return "", fmt.Errorf("failed to parse RemoteAddr: %v", err)
	}

	if validIP := net.ParseIP(ip); validIP != nil {
		return ip, nil
	}

	return "", fmt.Errorf("no valid IP address found")
}

func getRequestContext(r *http.Request) (map[string]interface{}, error) {
	ip, err := getRequesterIP(r)
	if err != nil {
		return nil, err
	}

	request_context := make(map[string]interface{})

	request_context["req_ip"] = ip

	return request_context, nil
}

func readAndReplaceBody(r *http.Request) (json.RawMessage, error) {
	if r.Body == nil {
		return []byte("{}"), nil
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body.Close()

	r.Body = io.NopCloser(bytes.NewBuffer(data))

	if len(data) == 0 {
		return []byte("{}"), nil
	}

	return data, nil
}

func convertMapToJson(data map[string]string) (json.RawMessage, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		bytes = []byte("{}")
	}

	return json.RawMessage(bytes), nil
}

// TODO: handle keys with multiple values
func convertHeaderToJson(headers http.Header) (json.RawMessage, error) {
	headerMap := make(map[string]string)
	for key, values := range headers {
		headerMap[key] = values[0]
	}

	return convertMapToJson(headerMap)
}

func getRequestDetails(r *http.Request) (*RequestDetails, error) {
	body, err := readAndReplaceBody(r)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %w", err)
	}

	headersJson, err := convertHeaderToJson(r.Header)
	if err != nil {
		return nil, fmt.Errorf("error reading request header: %w", err)
	}

	//TODO: handle keys with multiple values
	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		queryParams[key] = values[0]
	}

	queryParamsJson, err := convertMapToJson(queryParams)
	if err != nil {
		return nil, fmt.Errorf("error reading query parameters: %w", err)
	}

	details := &RequestDetails{
		Path:            r.URL.Path,
		Method:          common.HttpMethod(r.Method),
		QueryParameters: queryParamsJson,
		Headers:         headersJson,
		Body:            body,
	}

	return details, nil
}

func GetTxnTokenMiddleware(txnTokenServiceURL *url.URL, x509Source *workloadapi.X509Source, tokenetesSpiffeID spiffeid.ID, logger *zap.Logger) func(http.Handler) http.Handler {
	tokenetesMtlsClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsconfig.MTLSClientConfig(x509Source, x509Source, tlsconfig.AuthorizeID(tokenetesSpiffeID)),
		},
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestDetails, err := getRequestDetails(r)
			if err != nil {
				logger.Error("Error creating request details.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			oidcIdToken, ok := r.Context().Value(common.OIDC_ID_TOKEN_CONTEXT_KEY).(string)
			if !ok {
				logger.Error("Failed to retrieve OIDC id_token.")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			requestDetailsJSON, err := json.Marshal(requestDetails)
			if err != nil {
				logger.Error("Failed to marshal request details.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			encodedRequestDetails := base64.RawURLEncoding.EncodeToString(requestDetailsJSON)

			requestContext, err := getRequestContext(r)
			if err != nil {
				logger.Error("Error generating request context.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			requestContextJSON, err := json.Marshal(requestContext)
			if err != nil {
				logger.Error("Failed to marshal request context.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			encodedRequestContext := base64.RawURLEncoding.EncodeToString(requestContextJSON)

			requestData := url.Values{}
			requestData.Set("grant_type", GRANT_TYPE)
			requestData.Set("requested_token_type", TXN_TOKEN_TYPE)
			requestData.Set("audience", AUDIENCE)
			requestData.Set("subject_token", oidcIdToken)
			requestData.Set("subject_token_type", OIDC_ID_TOKEN_TYPE)
			requestData.Set("request_details", encodedRequestDetails)
			requestData.Set("request_context", encodedRequestContext)

			tokenEndpointURL := common.AppendPathToURL(txnTokenServiceURL, "token_endpoint").String()

			req, err := http.NewRequest(http.MethodPost, tokenEndpointURL, bytes.NewBufferString(requestData.Encode()))
			if err != nil {
				logger.Error("Failed to create the http request for txn-token endpoint.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, err := tokenetesMtlsClient.Do(req)
			if err != nil {
				logger.Error("Failed to request txn token from tokenetes.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logger.Error("Received non-ok http status from tokenetes service.", zap.Int("status", resp.StatusCode))

				if resp.StatusCode == http.StatusForbidden {
					http.Error(w, "Access Forbidden", http.StatusForbidden)
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}

				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("Failed to read the response from tokenetes", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			var token txnToken
			if err := json.Unmarshal(body, &token); err != nil {
				logger.Error("Failed to parse transaction token", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			if token.IssuedTokenType != TXN_TOKEN_TYPE {
				logger.Error("Issued invalid token type in txn-token response.", zap.String("token-type", string(token.IssuedTokenType)))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			if token.AccessToken == "" {
				logger.Error("Received empty access token from tokenetes.")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			r.Header.Set("Txn-Token", token.AccessToken)

			// ⚠️ Setting the "Txn-Token" header in the response as well. This is done only for this example application to demonstrate
			// TraTs generation in the application's UI interactively. There is no reason to do this in a real application, and it should not be
			// done as it can be a security threat. ⚠️
			w.Header().Set("Txn-Token", token.AccessToken)

			next.ServeHTTP(w, r)
		})
	}
}
