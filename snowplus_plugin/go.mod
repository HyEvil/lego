module yym/snowpluslib

go 1.14

require (
	github.com/emicklei/proto v1.9.0
	github.com/go-playground/validator/v10 v10.2.0
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/kubernetes/gengo v0.0.0-20200205140755-e0e292d8aa12
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sahilm/fuzzy v0.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cast v1.3.0
	github.com/thoas/go-funk v0.6.0
	golang.org/x/tools v0.0.0-20200318031718-dba9bee06b6c
	k8s.io/gengo v0.0.0-20200205140755-e0e292d8aa12 // indirect
	k8s.io/klog v1.0.0 // indirect
)

replace golang.org/x/tools v0.0.0-20200318031718-dba9bee06b6c => ./lib/tools
