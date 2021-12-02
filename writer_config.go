package transports

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"io"
	"math"
	"net"
)

func NewWriter(options ...writerOption) io.WriteCloser {
	var config writerConfiguration
	WriterOptions.apply(options...)(&config)

	var this io.WriteCloser
	this = newDialWriter(config) // lowest level: establish connection

	if config.CustomCompression != nil {
		this = config.CustomCompression(this) // TODO: does this close the underlying stream?
	}

	// TODO: framing writer

	return this
}

func (writerSingleton) DialContext(value context.Context) writerOption {
	return func(this *writerConfiguration) { this.DialContext = value }
}
func (writerSingleton) DialNetwork(value string) writerOption {
	return func(this *writerConfiguration) { this.DialNetwork = value }
}
func (writerSingleton) DialAddress(value string) writerOption {
	return func(this *writerConfiguration) { this.DialAddress = value }
}
func (writerSingleton) Dialer(value Dialer) writerOption {
	return func(this *writerConfiguration) { this.Dialer = value }
}

func (writerSingleton) TLSConfig(value *tls.Config) writerOption {
	return func(this *writerConfiguration) { this.TLSConfig = value }
}
func (writerSingleton) BufferSize(value uint32) writerOption {
	return func(this *writerConfiguration) { this.BufferSize = value }
}
func (writerSingleton) GZipCompression(level int) writerOption {
	return WriterOptions.CustomCompression(func(writer io.Writer) io.WriteCloser {
		compressor, _ := gzip.NewWriterLevel(writer, level)
		return compressor
	})
}
func (writerSingleton) CustomCompression(value func(io.Writer) io.WriteCloser) writerOption {
	return func(this *writerConfiguration) { this.CustomCompression = value }
}
func (writerSingleton) MaxMessageSize(value uint16) writerOption {
	return func(this *writerConfiguration) { this.MaxMessageSize = value }
}

func (writerSingleton) Logger(value Logger) writerOption {
	return func(this *writerConfiguration) { this.Logger = value }
}
func (writerSingleton) Monitor(value Monitor) writerOption {
	return func(this *writerConfiguration) { this.Monitor = value }
}

func (writerSingleton) apply(options ...writerOption) writerOption {
	return func(this *writerConfiguration) {
		for _, item := range WriterOptions.defaults(options...) {
			item(this)
		}
	}
}
func (writerSingleton) defaults(options ...writerOption) []writerOption {
	return append([]writerOption{
		WriterOptions.DialContext(context.Background()),
		WriterOptions.DialNetwork("tcp"),
		WriterOptions.DialAddress(""),
		WriterOptions.Dialer(defaultDialer),

		WriterOptions.BufferSize(0),
		WriterOptions.TLSConfig(nil),
		WriterOptions.GZipCompression(gzip.BestCompression),
		WriterOptions.CustomCompression(nil),
		WriterOptions.MaxMessageSize(math.MaxUint16),

		WriterOptions.Logger(writerEmpty),
		WriterOptions.Monitor(writerEmpty),
	}, options...)
}

var (
	defaultDialer Dialer = &net.Dialer{}
	writerEmpty          = &writerNop{}
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type writerConfiguration struct {
	DialContext       context.Context
	DialNetwork       string
	DialAddress       string
	Dialer            Dialer
	BufferSize        uint32
	MaxMessageSize    uint16
	CustomCompression func(io.Writer) io.WriteCloser
	TLSConfig         *tls.Config
	Logger            Logger
	Monitor           Monitor
}
type writerOption func(*writerConfiguration)
type writerSingleton struct{}

var WriterOptions writerSingleton

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type writerNop struct{}

func (*writerNop) Write([]byte) (int, error)     { return 0, nil }
func (*writerNop) Close() error                  { return nil }
func (*writerNop) Printf(string, ...interface{}) {}
