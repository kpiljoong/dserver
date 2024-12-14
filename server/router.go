package server

import (
	"dserver/config"
	"dserver/utils"
	"sync"

	"log"
	"net/http"
)

var (
	mux     *http.ServeMux
	muxLock sync.RWMutex
)

// GetRouter returns the current router.
func GetRouter() *http.ServeMux {
	muxLock.RLock()
	defer muxLock.RUnlock()
	return mux
}

// SetRouter updates the router dynamically.
func SetRouter(newMux *http.ServeMux) {
	muxLock.Lock()
	defer muxLock.Unlock()
	mux = newMux
}

// DynamicHandler is a wrapper that dynamically serves requests using the latest router.
type DynamicHandler struct{}

func (h *DynamicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := GetRouter()
	if router != nil {
		router.ServeHTTP(w, r)
	} else {
		http.Error(w, "No router configured", http.StatusInternalServerError)
	}
}

func RegisterRoutes(cfg *config.Config, verbose bool) {
	newMux := http.NewServeMux()
	routeMap := make(map[string][]config.RouteConfig)

	for _, route := range cfg.Routes {
		routeMap[route.Path] = append(routeMap[route.Path], route)
	}

	for path, routes := range routeMap {
		if verbose {
			log.Printf("Registering route: %s\n", path)
		}
		newMux.HandleFunc(path, CreateHandler(routes, verbose))
	}

	SetRouter(newMux)
}

func CreateHandler(routeConfigs []config.RouteConfig, verbose bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routeConfigs {
			if route.Method != r.Method {
				continue
			}

			if len(route.Responses) == 0 {
				utils.ApplyDelay(route.DelayMs)
				utils.WriteResponse(w, route.StatusCode, route.ContentType, route.ResponseBody)
				return
			}

			requestParams := r.URL.Query()
			for _, response := range route.Responses {
				if !utils.MatchQueryParams(requestParams, response.Query) {
					continue
				}

				delayMs := route.DelayMs
				if response.DelayMs != nil {
					delayMs = *response.DelayMs
				}
				utils.ApplyDelay(delayMs)

				statusCode := route.StatusCode
				if response.StatusCode != nil {
					statusCode = *response.StatusCode
				}

				contentType := route.ContentType
				if response.ContentType != nil {
					contentType = *response.ContentType
				}

				utils.WriteResponse(w, statusCode, contentType, response.ResponseBody)
				return
			}
			http.Error(w, "No matching response found", http.StatusNotFound)
			return
		}
		http.Error(w, "No matching route found", http.StatusNotFound)
	}
}
