package filters

import (
	"fmt"
	"net/http"

	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/topology"
)

func WithPopulatedTopologyInfo(handler http.Handler, populator topology.Populator) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet || request.Method == http.MethodConnect {
			topologyInfo, err := populator.PopulateTopologyInfo(request)
			if err != nil {
				responsewriters.InternalError(writer, request, fmt.Errorf("failed to populate TopologyInfo: %v", err))
				return
			}

			request = request.WithContext(topology.WithTopologyInfo(request.Context(), topologyInfo))
		}
		handler.ServeHTTP(writer, request)
	})
}
