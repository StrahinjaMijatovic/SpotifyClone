package handlers

import (
	"net/http"

	"example.com/api-gateway/docs"
	"github.com/gin-gonic/gin"
)

// SwaggerJSON vraća swagger.json specifikaciju
func SwaggerJSON(c *gin.Context) {
	data, err := docs.SwaggerJSON.ReadFile("swagger.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read swagger.json"})
		return
	}
	c.Data(http.StatusOK, "application/json", data)
}

// SwaggerUI servira Swagger UI HTML stranicu
func SwaggerUI(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Spotify Clone API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
        html { box-sizing: border-box; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
        .swagger-ui .topbar { display: none; }
        .swagger-ui .info .title { color: #1db954; }
        .swagger-ui .btn.execute { background-color: #1db954; border-color: #1db954; }
        .swagger-ui .btn.execute:hover { background-color: #1ed760; border-color: #1ed760; }
        .swagger-ui .opblock.opblock-post { border-color: #1db954; background: rgba(29, 185, 84, 0.1); }
        .swagger-ui .opblock.opblock-post .opblock-summary-method { background: #1db954; }
        .swagger-ui .opblock.opblock-get { border-color: #61affe; background: rgba(97, 175, 254, 0.1); }
        .swagger-ui .opblock.opblock-delete { border-color: #f93e3e; background: rgba(249, 62, 62, 0.1); }
        .swagger-ui .opblock.opblock-put { border-color: #fca130; background: rgba(252, 161, 48, 0.1); }
        .custom-header {
            background: linear-gradient(135deg, #1db954 0%, #191414 100%);
            padding: 20px 40px;
            color: white;
        }
        .custom-header h1 { margin: 0 0 5px 0; font-size: 28px; }
        .custom-header p { margin: 0; opacity: 0.9; }
    </style>
</head>
<body>
    <div class="custom-header">
        <h1>🎵 Spotify Clone API</h1>
        <p>Interaktivna API dokumentacija - testirajte endpointe direktno iz browsera</p>
    </div>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                persistAuthorization: true,
                tagsSorter: "alpha",
                operationsSorter: "alpha"
            });
            window.ui = ui;
        };
    </script>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
