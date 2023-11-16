package filters

import (
	"fmt"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/topology"
	"net/http"
)

func WithTopologyInfo(handler http.Handler, resolver topology.TopologyInfoResolver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		info, err := resolver.NewTopologyInfo(req)
		if err != nil {
			responsewriters.InternalError(w, req, fmt.Errorf("failed to create TopologyInfo: %v", err))
			return
		}

		req = req.WithContext(topology.WithTopologyInfo(ctx, info))
		handler.ServeHTTP(w, req)
	})
}
