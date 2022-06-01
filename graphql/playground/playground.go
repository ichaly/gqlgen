package playground

import (
	"html/template"
	"net/http"
	"net/url"
)

var page = template.Must(template.New("graphiql").Parse(`<!DOCTYPE html>
<html>
<head>
  <title>{{.title}}</title>
  <style>
      html,
      body {
          height: 100%;
          margin: 0;
          overflow: hidden;
          width: 100%;
      }

      #graphiql {
          height: 100vh;
      }
  </style>
  <link
      rel="stylesheet"
      href="https://cdn.jsdelivr.net/npm/graphiql-with-graphiql-explorer@0.15.1/graphiqlWithExtensions.css"
      integrity="sha256-PiHC3pxlImCsNqR79k/SieFXPFOKGNpmqqu83g1gghw="
      crossorigin="anonymous">
  <script
      src="https://cdn.jsdelivr.net/npm/whatwg-fetch@2.0.3/fetch.min.js"
      integrity="sha384-dcF7KoWRaRpjcNbVPUFgatYgAijf8DqW6NWuqLdfB5Sb4Cdbb8iHX7bHsl9YhpKa"
      crossorigin="anonymous"
  ></script>
  <script
      src="https://cdn.jsdelivr.net/npm/react@16.8.6/umd/react.production.min.js"
      integrity="sha384-qn+ML/QkkJxqn4LLs1zjaKxlTg2Bl/6yU/xBTJAgxkmNGc6kMZyeskAG0a7eJBR1"
      crossorigin="anonymous"
  ></script>
  <script
      src="https://cdn.jsdelivr.net/npm/react-dom@16.8.6/umd/react-dom.production.min.js"
      integrity="sha384-85IMG5rvmoDsmMeWK/qUU4kwnYXVpC+o9hoHMLi4bpNR+gMEiPLrvkZCgsr7WWgV"
      crossorigin="anonymous"
  ></script>
  <script
      src="https://cdn.jsdelivr.net/npm/graphiql-with-graphiql-explorer@0.15.1/graphiqlWithExtensions.min.js"
      integrity="sha256-4+5xt0s1fmQi6n064zU/ZDcCvgroBNF/kWNYVdiTNPg="
      crossorigin="anonymous"
  ></script>
</head>
<body>
<div id="graphiql"></div>
<script>
{{- if .endpointIsAbsolute}}
  var fetchURL = {{.endpoint}};
{{- else}}
  var fetchURL = location.protocol + '//' + location.host + {{.endpoint}};
{{- end}}
  function graphQLFetcher (graphQLParams) {
    var headers = {
      Accept: 'application/json',
      'Content-Type': 'application/json'
    }
    return fetch(fetchURL, {
      method: 'post',
      headers: headers,
      body: JSON.stringify(graphQLParams)
    }).then(function (response) {
      return response.text()
    }).then(function (responseBody) {
      try {
        return JSON.parse(responseBody)
      } catch (error) {
        return responseBody
      }
    })
  }

  ReactDOM.render(
    React.createElement(GraphiQLWithExtensions.GraphiQLWithExtensions, {
      fetcher: graphQLFetcher
    }),
    document.getElementById('graphiql')
  )
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
			"version":              "1.8.2",
			"cssSRI":               "sha256-CDHiHbYkDSUc3+DS2TU89I9e2W3sJRUOqSmp7JC+LBw=",
			"jsSRI":                "sha256-X8vqrqZ6Rvvoq4tvRVM3LoMZCQH8jwW92tnX0iPiHPc=",
			"reactSRI":             "sha256-Ipu/TQ50iCCVZBUsZyNJfxrDk0E2yhaEIz0vqI+kFG8=",
			"reactDOMSRI":          "sha256-nbMykgB6tsOFJ7OdVmPpdqMFVk4ZsqWocT6issAPUF0=",
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
