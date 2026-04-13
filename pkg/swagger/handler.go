package swagger

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
)

//go:embed openapi_spec.json
var openAPISpec string

// SwaggerHandler serves Swagger UI
func SwaggerHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(GetSwaggerHTML())
}

// OpenAPIHandler serves OpenAPI specification
func OpenAPIHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	return c.SendString(GetOpenAPISpec())
}

// GetSwaggerHTML returns Swagger UI HTML
func GetSwaggerHTML() string {
	return `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Swagger UI</title>

  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
  
  <style>
    body {
      margin: 0;
      background: #fafafa;
    }
  </style>
</head>

<body>
  <div id="swagger-ui"></div>

  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-standalone-preset.js"></script>

  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/api/docs/openapi.json",
        dom_id: "#swagger-ui",

        deepLinking: true,

        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],

        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],

        layout: "StandaloneLayout",

        // 🔥 WAJIB (biar token tidak hilang)
        persistAuthorization: true
      })
    }
  </script>
</body>
</html>`
}

// GetOpenAPISpec returns OpenAPI 3.0.0 specification
func GetOpenAPISpec() string {
	return openAPISpec
}
