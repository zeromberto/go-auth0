package management

import "encoding/json"

const (
	// LogStreamTypeAmazonEventBridge constant.
	LogStreamTypeAmazonEventBridge = "eventbridge"
	// LogStreamTypeAzureEventGrid constant.
	LogStreamTypeAzureEventGrid = "eventgrid"
	// LogStreamTypeHTTP constant.
	LogStreamTypeHTTP = "http"
	// LogStreamTypeDatadog constant.
	LogStreamTypeDatadog = "datadog"
	// LogStreamTypeSplunk constant.
	LogStreamTypeSplunk = "splunk"
	// LogStreamTypeSumo constant.
	LogStreamTypeSumo = "sumo"
)

// LogStream is used to export tenant log
// events to a log event analysis service.
//
// See: https://auth0.com/docs/customize/log-streams
type LogStream struct {
	// The hook's identifier.
	ID *string `json:"id,omitempty"`

	// The name of the hook. Can only contain alphanumeric characters, spaces
	// and '-'. Can neither start nor end with '-' or spaces.
	Name *string `json:"name,omitempty"`

	// The type of the log-stream. Can be one of "http", "eventbridge",
	// "eventgrid", "datadog" or "splunk".
	Type *string `json:"type,omitempty"`

	// The status of the log-stream. Can be one of "active", "paused", or "suspended".
	Status *string `json:"status,omitempty"`

	// Sink for validation.
	Sink interface{} `json:"-"`
}

// MarshalJSON is a custom serializer for the LogStream type.
func (ls *LogStream) MarshalJSON() ([]byte, error) {
	type logStream LogStream
	type logStreamWrapper struct {
		*logStream
		RawSink json.RawMessage `json:"sink,omitempty"`
	}

	w := &logStreamWrapper{(*logStream)(ls), nil}

	if ls.Sink != nil {
		b, err := json.Marshal(ls.Sink)
		if err != nil {
			return nil, err
		}
		w.RawSink = b
	}

	return json.Marshal(w)
}

// UnmarshalJSON is a custom deserializer for the LogStream type.
func (ls *LogStream) UnmarshalJSON(b []byte) error {
	type logStream LogStream
	type logStreamWrapper struct {
		*logStream
		RawSink json.RawMessage `json:"sink,omitempty"`
	}

	w := &logStreamWrapper{(*logStream)(ls), nil}

	err := json.Unmarshal(b, w)
	if err != nil {
		return err
	}

	if ls.Type != nil {
		var v interface{}

		switch *ls.Type {
		case LogStreamTypeAmazonEventBridge:
			v = &LogStreamSinkAmazonEventBridge{}
		case LogStreamTypeAzureEventGrid:
			v = &LogStreamSinkAzureEventGrid{}
		case LogStreamTypeHTTP:
			v = &LogStreamSinkHTTP{}
		case LogStreamTypeDatadog:
			v = &LogStreamSinkDatadog{}
		case LogStreamTypeSplunk:
			v = &LogStreamSinkSplunk{}
		case LogStreamTypeSumo:
			v = &LogStreamSinkSumo{}
		default:
			v = make(map[string]interface{})
		}

		err = json.Unmarshal(w.RawSink, &v)
		if err != nil {
			return err
		}

		ls.Sink = v
	}

	return nil
}

// LogStreamSinkAmazonEventBridge is used to export logs to Amazon EventBridge.
type LogStreamSinkAmazonEventBridge struct {
	// AWS Account Id
	AccountID *string `json:"awsAccountId,omitempty"`
	// AWS Region
	Region *string `json:"awsRegion,omitempty"`
	// AWS Partner Event Source
	PartnerEventSource *string `json:"awsPartnerEventSource,omitempty"`
}

// LogStreamSinkAzureEventGrid is used to export logs to Azure Event Grid.
type LogStreamSinkAzureEventGrid struct {
	// Azure Subscription Id
	SubscriptionID *string `json:"azureSubscriptionId,omitempty"`
	// Azure Resource Group
	ResourceGroup *string `json:"azureResourceGroup,omitempty"`
	// Azure Region
	Region *string `json:"azureRegion,omitempty"`
	// Azure Partner Topic
	PartnerTopic *string `json:"azurePartnerTopic,omitempty"`
}

// LogStreamSinkHTTP is used to export logs to Custom Webhooks.
type LogStreamSinkHTTP struct {
	// HTTP ContentFormat
	ContentFormat *string `json:"httpContentFormat,omitempty"`
	// HTTP ContentType
	ContentType *string `json:"httpContentType,omitempty"`
	// HTTP Endpoint
	Endpoint *string `json:"httpEndpoint,omitempty"`
	// HTTP Authorization
	Authorization *string `json:"httpAuthorization,omitempty"`
	// Custom HTTP headers
	CustomHeaders []*LogStreamSinkHTTPCustomHeaders `json:"httpCustomHeaders,omitempty"`
}

type LogStreamSinkHTTPCustomHeaders struct {
	// The custom header key
	Header *string `json:"header,omitempty"`
	// The custom header value
	Value *string `json:"value,omitempty"`
}

// LogStreamSinkDatadog is used to export logs to Datadog.
type LogStreamSinkDatadog struct {
	// Datadog Region
	Region *string `json:"datadogRegion,omitempty"`
	// Datadog Api Key
	APIKey *string `json:"datadogApiKey,omitempty"`
}

// LogStreamSinkSplunk is used to export logs to Splunk.
type LogStreamSinkSplunk struct {
	// Splunk Domain
	Domain *string `json:"splunkDomain,omitempty"`
	// Splunk Token
	Token *string `json:"splunkToken,omitempty"`
	// Splunk Port
	Port *string `json:"splunkPort,omitempty"`
	// Splunk Secure
	Secure *bool `json:"splunkSecure,omitempty"`
}

// LogStreamSinkSumo is used to export logs to Sumo Logic.
type LogStreamSinkSumo struct {
	// Sumo Source Address
	SourceAddress *string `json:"sumoSourceAddress,omitempty"`
}

// LogStreamManager manages Auth0 LogStream resources.
type LogStreamManager struct {
	*Management
}

func newLogStreamManager(m *Management) *LogStreamManager {
	return &LogStreamManager{m}
}

// Create a log stream.
//
// See: https://auth0.com/docs/api/management/v2#!/log-streams
func (m *LogStreamManager) Create(l *LogStream, opts ...RequestOption) error {
	return m.Request("POST", m.URI("log-streams"), l, opts...)
}

// Read a log stream.
//
// See: https://auth0.com/docs/api/management/v2#!/Log_Streams/get_log_streams_by_id
func (m *LogStreamManager) Read(id string, opts ...RequestOption) (l *LogStream, err error) {
	err = m.Request("GET", m.URI("log-streams", id), &l, opts...)
	return
}

// List all log streams.
//
// See: https://auth0.com/docs/api/management/v2#!/log-streams/get_log_streams
func (m *LogStreamManager) List(opts ...RequestOption) (ls []*LogStream, err error) {
	err = m.Request("GET", m.URI("log-streams"), &ls, opts...)
	return
}

// Update a log stream.
//
// The following fields may be updated in a PATCH operation: Name, Status, Sink.
//
// Note: For log streams of type eventbridge and eventgrid, updating the sink is
// not permitted.
//
// See: https://auth0.com/docs/api/management/v2#!/log-streams
func (m *LogStreamManager) Update(id string, l *LogStream, opts ...RequestOption) (err error) {
	return m.Request("PATCH", m.URI("log-streams", id), l, opts...)
}

// Delete a log stream.
//
// See: https://auth0.com/docs/api/management/v2#!/log-streams
func (m *LogStreamManager) Delete(id string, opts ...RequestOption) (err error) {
	return m.Request("DELETE", m.URI("log-streams", id), nil, opts...)
}
