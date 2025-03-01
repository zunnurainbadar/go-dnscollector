package loggers

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/dmachard/go-dnscollector/dnsutils"
	"github.com/dmachard/go-logger"
)

func TestPrometheus_BadAuth(t *testing.T) {
	// init the logger
	config := dnsutils.GetFakeConfig()
	g := NewPrometheus(config, logger.New(false), "dev", "test")

	tt := []struct {
		name       string
		uri        string
		handler    func(w http.ResponseWriter, r *http.Request)
		method     string
		statusCode int
	}{
		{
			name:       "total clients",
			uri:        "/metrics",
			handler:    g.httpServer.Handler.ServeHTTP,
			method:     http.MethodGet,
			statusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// init httptest
			request := httptest.NewRequest(tc.method, tc.uri, strings.NewReader(""))
			request.SetBasicAuth(config.Loggers.Prometheus.BasicAuthLogin, "badpassword")
			responseRecorder := httptest.NewRecorder()

			// call handler
			tc.handler(responseRecorder, request)

			// checking status code
			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}
		})
	}
}

func TestPrometheus_GetMetrics(t *testing.T) {
	// init the logger
	config := dnsutils.GetFakeConfig()
	g := NewPrometheus(config, logger.New(false), "dev", "test")

	// record one dns message to simulate some incoming data
	noerror_record := dnsutils.GetFakeDnsMessage()
	noerror_record.DNS.Type = dnsutils.DnsQuery
	g.Record(noerror_record)

	nx_record := dnsutils.GetFakeDnsMessage()
	nx_record.DNS.Type = dnsutils.DnsReply
	nx_record.DNS.Rcode = dnsutils.DNS_RCODE_NXDOMAIN
	g.Record(nx_record)

	sf_record := dnsutils.GetFakeDnsMessage()
	sf_record.DNS.Type = dnsutils.DnsReply
	sf_record.DNS.Rcode = dnsutils.DNS_RCODE_SERVFAIL
	g.Record(sf_record)

	// compute metrics, this function is called every second
	g.ComputeMetrics()

	tt := []struct {
		name       string
		method     string
		handler    func(w http.ResponseWriter, r *http.Request)
		want       string
		statusCode int
	}{
		{
			name:       "total domains",
			method:     http.MethodGet,
			handler:    g.httpServer.Handler.ServeHTTP,
			want:       config.Loggers.Prometheus.PromPrefix + `_domains_total{stream_id="collector"} 1`,
			statusCode: http.StatusOK,
		},
		{
			name:       "total requesters",
			method:     http.MethodGet,
			handler:    g.httpServer.Handler.ServeHTTP,
			want:       config.Loggers.Prometheus.PromPrefix + `_requesters_total{stream_id="collector"} 1`,
			statusCode: http.StatusOK,
		},
		{
			name:       "total nxdomain",
			method:     http.MethodGet,
			handler:    g.httpServer.Handler.ServeHTTP,
			want:       config.Loggers.Prometheus.PromPrefix + `_nxdomains_total{stream_id="collector"} 1`,
			statusCode: http.StatusOK,
		},
		{
			name:       "total sfdomain",
			method:     http.MethodGet,
			handler:    g.httpServer.Handler.ServeHTTP,
			want:       config.Loggers.Prometheus.PromPrefix + `_sfdomains_total{stream_id="collector"} 1`,
			statusCode: http.StatusOK,
		},
		{
			name:       "total dns messages",
			method:     http.MethodGet,
			handler:    g.httpServer.Handler.ServeHTTP,
			want:       config.Loggers.Prometheus.PromPrefix + `_dnsmessages_total{stream_id="collector"} 3`,
			statusCode: http.StatusOK,
		},
		{
			name:       "total queries",
			method:     http.MethodGet,
			handler:    g.httpServer.Handler.ServeHTTP,
			want:       config.Loggers.Prometheus.PromPrefix + `_queries_total{stream_id="collector"} 1`,
			statusCode: http.StatusOK,
		},
		{
			name:       "total replies",
			method:     http.MethodGet,
			handler:    g.httpServer.Handler.ServeHTTP,
			want:       config.Loggers.Prometheus.PromPrefix + `_replies_total{stream_id="collector"} 2`,
			statusCode: http.StatusOK,
		},
		{
			name:       "total qtypes",
			method:     http.MethodGet,
			handler:    g.httpServer.Handler.ServeHTTP,
			want:       config.Loggers.Prometheus.PromPrefix + `_qtypes_total{query_type="A",stream_id="collector"} 3`,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// init httptest
			request := httptest.NewRequest(tc.method, "/metrics", strings.NewReader(""))
			request.SetBasicAuth(config.Loggers.Prometheus.BasicAuthLogin, config.Loggers.Prometheus.BasicAuthPwd)
			responseRecorder := httptest.NewRecorder()

			// call handler
			tc.handler(responseRecorder, request)

			// checking status code
			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			// checking content
			metrics := strings.TrimSpace(responseRecorder.Body.String())
			if regexp.MustCompile(tc.want).MatchString(metrics) != true {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}
