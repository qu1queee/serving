module knative.dev/serving

go 1.15

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/cli v0.0.0-20200210162036-a4bedce16568 // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible
	github.com/gogo/protobuf v1.3.1
	github.com/google/go-cmp v0.5.4
	github.com/google/go-containerregistry v0.2.1
	github.com/google/gofuzz v1.1.0
	github.com/google/mako v0.0.0-20190821191249-122f8dcef9e3
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/golang-lru v0.5.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/prometheus/client_golang v1.8.0
	github.com/prometheus/client_model v0.2.0
	github.com/tsenart/vegeta/v12 v12.8.4
	go.opencensus.io v0.22.5
	go.uber.org/atomic v1.7.0
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/goleak v1.1.10
	go.uber.org/zap v1.16.0
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.34.0
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	k8s.io/code-generator v0.19.7
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	knative.dev/caching v0.0.0-20210125050654-45e8de7ff96e
	knative.dev/hack v0.0.0-20210120165453-8d623a0af457
	knative.dev/networking v0.0.0-20210125050654-94433ab7f620
	knative.dev/pkg v0.0.0-20210124203454-7101e9d4f6c6
	sigs.k8s.io/yaml v1.2.0
)
