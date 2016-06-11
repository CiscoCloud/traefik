package middlewares

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/containous/traefik/types"
	"github.com/streamrail/concurrent-map"
)

const (
	loggerReqidHeader = "X-Traefik-Reqid"
)

/*
Logger writes each request and its response to the access log.
It gets some information from the logInfoResponseWriter set up by previous middleware.
*/
type Logger struct {
	file          *os.File
	kafkaProducer *sarama.AsyncProducer
	topic         string
}

// AccessLog holds configuration of the access log
type AccessLog struct {
	Filename string // Name of file on server
	Brokers  string // Kafka broker names/IP addresses (comma-separated)
	Topic    string // Kafka topic
}

type FrontendInfo struct {
	Name string          // Frontend name
	Info *types.Frontend // Frontend information
}

// MonascaSpec describes the structure of the metrics message written to Kafka.
// See https://wiki.openstack.org/wiki/Monasca/Message_Schema#Metrics_Message
type MonascaSpec struct {
	Metric       MetricSpec `json:"metric"`
	Meta         MetaSpec   `json:"meta"`
	CreationTime int64      `json:"creation_time"`
}

// MetricSpec describes the metric portion of the message.
// Dimensions are keys for message access.
type MetricSpec struct {
	Name       string            `json:"name"`
	Dimensions map[string]string `json:"dimensions"`
	Timestamp  int64             `json:"timestamp"`
	Value      uint64            `json:"value"`
}

// MetaSpec describes meta data for the message
type MetaSpec struct {
	TenantID string `json:"tenantId"`
	Region   string `json:"region"`
}

// Logging handler to log frontend name, backend name, and elapsed time
type frontendBackendLoggingHandler struct {
	reqid       string
	logger      *Logger
	handlerFunc http.HandlerFunc
}

var (
	reqidCounter        uint64       // Request ID
	infoRwMap           = cmap.New() // Map of reqid to response writer
	backend2FrontendMap *map[string]*FrontendInfo
	label2DimensionMap  = map[string]string{
		"project_id":   "SHIPPED_PROJECT_ID",
		"project_name": "SHIPPED_PROJECT_NAME",
		"service_id":   "SHIPPED_SERVICE_ID",
		"service_name": "SHIPPED_SERVICE_NAME",
		"env_id":       "SHIPPED_ENVIRONMENT_ID",
		"env_name":     "SHIPPED_ENVIRONMENT_NAME",
	}
)

// logInfoResponseWriter is a wrapper of type http.ResponseWriter
// that tracks frontend and backend names and request status and size
type logInfoResponseWriter struct {
	rw           http.ResponseWriter
	backend      string
	frontendInfo *FrontendInfo
	status       int
	size         int
}

// NewLogger returns a new Logger instance with either or both of an output file and a Kafka producer
func NewLogger(accessLog AccessLog) *Logger {
	logger := Logger{}
	if len(accessLog.Filename) > 0 {
		if fi, err := os.OpenFile(accessLog.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
			log.Fatal("Error opening file", err)
		} else {
			logger.file = fi
		}
	}
	if len(accessLog.Brokers) > 0 {
		config := sarama.NewConfig()
		config.Producer.Retry.Max = 5
		config.Producer.RequiredAcks = sarama.WaitForAll
		brokers := strings.Split(accessLog.Brokers, ",")
		if producer, err := sarama.NewAsyncProducer(brokers, config); err != nil {
			log.Fatal("Unable to create access log Kafka producer", err)
		} else {
			logger.kafkaProducer = &producer
			logger.topic = accessLog.Topic
		}
	}
	return &logger
}

// SetBackend2FrontendMap is called by server.go to set up frontend translation
func SetBackend2FrontendMap(newMap *map[string]*FrontendInfo) {
	backend2FrontendMap = newMap
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if l.file == nil && l.kafkaProducer == nil {
		next(rw, r)
	} else {
		reqid := strconv.FormatUint(atomic.AddUint64(&reqidCounter, 1), 10)
		r.Header[loggerReqidHeader] = []string{reqid}
		defer deleteReqid(r, reqid)
		frontendBackendLoggingHandler{reqid, l, next}.ServeHTTP(rw, r)
	}
}

// Delete a reqid from the map and the request's headers
func deleteReqid(r *http.Request, reqid string) {
	infoRwMap.Remove(reqid)
	delete(r.Header, loggerReqidHeader)
}

// Save the backend name for the Logger
func saveBackendNameForLogger(r *http.Request, backendName string) {
	if reqidHdr := r.Header[loggerReqidHeader]; len(reqidHdr) == 1 {
		reqid := reqidHdr[0]
		if infoRw, ok := infoRwMap.Get(reqid); ok {
			infoRw.(*logInfoResponseWriter).SetBackend(backendName)
			infoRw.(*logInfoResponseWriter).SetFrontend((*backend2FrontendMap)[backendName])
		}
	}
}

// Close closes the Logger (i.e. the file).
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

// Logging handler to log frontend name, backend name, and elapsed time
func (fblh frontendBackendLoggingHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	infoRw := &logInfoResponseWriter{rw: rw}
	infoRwMap.Set(fblh.reqid, infoRw)
	fblh.handlerFunc(infoRw, req)

	username := "-"
	url := *req.URL
	if url.User != nil {
		if name := url.User.Username(); name != "" {
			username = name
		}
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}

	ts := startTime.Format("02/Jan/2006:15:04:05 -0700")
	method := req.Method
	uri := url.RequestURI()
	if qmIndex := strings.Index(uri, "?"); qmIndex > 0 {
		uri = uri[0:qmIndex]
	}
	proto := req.Proto
	referer := req.Referer()
	agent := req.UserAgent()

	frontendName, backendName := "", infoRw.GetBackend()
	frontendInfo := infoRw.GetFrontend()
	if frontendInfo != nil {
		frontendName = strings.TrimPrefix(frontendInfo.Name, "frontend-")
	}
	status := infoRw.GetStatus()
	size := infoRw.GetSize()

	elapsed := time.Now().UTC().Sub(startTime.UTC())
	if fblh.logger.file != nil {
		fmt.Fprintf(fblh.logger.file, `%s - %s [%s] "%s %s %s" %d %d "%s" "%s" %s "%s" "%s" %s%s`,
			host, username, ts, method, uri, proto, status, size, referer, agent, fblh.reqid, frontendName, backendName, elapsed, "\n")
	}
	if fblh.logger.kafkaProducer != nil && frontendInfo != nil && frontendInfo.Info != nil {
		if labels := frontendInfo.Info.Labels; labels != nil && len(labels["project_id"]) > 0 {
			timestamp := time.Now().Unix() * 1E3
			dimensions := map[string]string{
				"frontend": frontendName,
				"backend":  backendName,
				"url":      uri,
			}
			for label, dimension := range label2DimensionMap {
				if len(labels[label]) > 0 {
					dimensions[dimension] = labels[label]
				}
			}
			err := fblh.writeMetric(&MetricSpec{
				Name:       "elapsedMs-traefik",
				Dimensions: dimensions,
				Timestamp:  timestamp,
				Value:      uint64(elapsed / 1000000),
			})
			if err == nil && status >= 200 && status < 600 {
				metricName := ""
				switch {
				case status < 300:
					metricName = "2xx"
				case status < 400:
					metricName = "3xx"
				case status < 500:
					metricName = "4xx"
				default:
					metricName = "5xx"
				}
				fblh.writeMetric(&MetricSpec{
					Name:       metricName,
					Dimensions: dimensions,
					Timestamp:  timestamp,
					Value:      1,
				})
			}
		}
	}
}

/**
 * Write a metric to Kafka
 */
func (fblh frontendBackendLoggingHandler) writeMetric(metric *MetricSpec) (e error) {
	monascaSpec := MonascaSpec{
		Metric: *metric,
		Meta: MetaSpec{
			TenantID: metric.Dimensions["SHIPPED_PROJECT_ID"],
		},
		CreationTime: metric.Timestamp,
	}
	if b, err := json.Marshal(monascaSpec); err != nil {
		e = fmt.Errorf("Can't marshal metric %#v: %s", monascaSpec, err.Error())
		log.Errorf(e.Error())
	} else {
		msg := &sarama.ProducerMessage{
			Topic: fblh.logger.topic,
			Value: sarama.StringEncoder(b),
		}
		select {
		case (*fblh.logger.kafkaProducer).Input() <- msg:
			log.Debugf("Wrote metric %s=%v for %s", metric.Name, metric.Value, metric.Dimensions["frontend"])
		case err := <-(*fblh.logger.kafkaProducer).Errors():
			e = fmt.Errorf("Kafka write for %s failed: %s", metric.Dimensions["frontend"], err.Error())
			log.Errorf(e.Error())
		}
	}
	return
}

/**
 * Extract environment, project, and service from a Shipped appID
 */
func extractNamesFromAppID(appID string) (envName, projectName, serviceName, configID string) {
	if tokens := strings.SplitN(appID, "--", 4); len(tokens) == 4 {
		envName, projectName, serviceName, configID = tokens[0], tokens[1], tokens[2], tokens[3]
	}
	return
}

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (acc *AccessLog) String() string {
	return fmt.Sprintf("%#v", acc)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Currently unused.
func (acc *AccessLog) Set(value string) error {
	return nil
}

// Type is type of the struct
func (acc *AccessLog) Type() string {
	return fmt.Sprint("accesslog")
}

func (lirw *logInfoResponseWriter) Header() http.Header {
	return lirw.rw.Header()
}

func (lirw *logInfoResponseWriter) Write(b []byte) (int, error) {
	if lirw.status == 0 {
		lirw.status = http.StatusOK
	}
	size, err := lirw.rw.Write(b)
	lirw.size += size
	return size, err
}

func (lirw *logInfoResponseWriter) WriteHeader(s int) {
	lirw.rw.WriteHeader(s)
	lirw.status = s
}

func (lirw *logInfoResponseWriter) Flush() {
	f, ok := lirw.rw.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func (lirw *logInfoResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return lirw.rw.(http.Hijacker).Hijack()
}

func (lirw *logInfoResponseWriter) GetStatus() int {
	return lirw.status
}

func (lirw *logInfoResponseWriter) GetSize() int {
	return lirw.size
}

func (lirw *logInfoResponseWriter) GetBackend() string {
	return lirw.backend
}

func (lirw *logInfoResponseWriter) GetFrontend() *FrontendInfo {
	return lirw.frontendInfo
}

func (lirw *logInfoResponseWriter) SetBackend(backend string) {
	lirw.backend = backend
}

func (lirw *logInfoResponseWriter) SetFrontend(frontendInfo *FrontendInfo) {
	lirw.frontendInfo = frontendInfo
}
