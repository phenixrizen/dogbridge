package distro

// Components describes the collector components expected in the dogbridge distribution.
var Components = []string{
	"datadogreceiver",
	"statsdreceiver",
	"otlpreceiver",
	"filelogreceiver",
	"prometheusreceiver",
	"k8sattributesprocessor",
	"transformprocessor",
	"batchprocessor",
	"memorylimiterprocessor",
	"otlpexporter",
	"prometheusremotewriteexporter",
	"lokiexporter",
	"elasticsearchexporter",
}
