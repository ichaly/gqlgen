package playground

import (
	"html/template"
	"net/http"
	"net/url"
)

var page = template.Must(template.New("graphiql").Parse(`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{.title}}</title>
  <base href="//cdn.jsdelivr.net/npm/altair-static@{{.version}}/build/dist/">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <link rel="icon" type="image/x-icon" href="favicon.ico">
  <link href="styles.css" rel="stylesheet" />
</head>
<body>
  <app-root>
    <style>
      .loading-screen {
        display: none;
      }
    </style>
    <div class="loading-screen styled">
      <div class="loading-screen-inner">
        <div class="loading-screen-logo-container">
          <img src="assets/img/logo_350.svg" alt="Altair">
        </div>
        <div class="loading-screen-loading-indicator">
          <span class="loading-indicator-dot"></span>
          <span class="loading-indicator-dot"></span>
          <span class="loading-indicator-dot"></span>
        </div>
      </div>
    </div>
  </app-root>
  <script rel="preload" as="script" type="text/javascript" src="runtime.js"></script>
  <script rel="preload" as="script" type="text/javascript" src="polyfills.js"></script>
  <script rel="preload" as="script" type="text/javascript" src="main.js"></script>
  <script>
{{- if .endpointIsAbsolute}}
    const url = {{.endpoint}};
    const subscriptionUrl = {{.subscriptionEndpoint}};
{{- else}}
    const url = location.protocol + '//' + location.host + {{.endpoint}};
    const wsProto = location.protocol == 'https:' ? 'wss:' : 'ws:';
    const subscriptionUrl = wsProto + '//' + location.host + {{.endpoint}};
{{- end}}
    var altairOptions = {
      endpointURL: url,
      subscriptionsEndpoint: subscriptionUrl,
      initialHeaders:{},
      initialVariables:'{}'
    };
    window.addEventListener("load", function() {
      AltairGraphQL.init(altairOptions);
    });
  </script>
</body>
</html>`))

// Handler responsible for setting up the playground
func Handler(title string, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		err := page.Execute(w, map[string]interface{}{
			"title":                title,
			"endpoint":             endpoint,
			"endpointIsAbsolute":   endpointHasScheme(endpoint),
			"subscriptionEndpoint": getSubscriptionEndpoint(endpoint),
			"version":              "4.4.2",
		})
		if err != nil {
			panic(err)
		}
	}
}

// endpointHasScheme checks if the endpoint has a scheme.
func endpointHasScheme(endpoint string) bool {
	u, err := url.Parse(endpoint)
	return err == nil && u.Scheme != ""
}

// getSubscriptionEndpoint returns the subscription endpoint for the given
// endpoint if it is parsable as a URL, or an empty string.
func getSubscriptionEndpoint(endpoint string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}

	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	default:
		u.Scheme = "ws"
	}

	return u.String()
}
