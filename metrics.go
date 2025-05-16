package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {	
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}


func (apiCfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "type/html")
	w.WriteHeader(http.StatusOK)

	message := fmt.Sprintf(`
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>`, apiCfg.fileserverHits.Load())
	w.Write([]byte(message))
}
