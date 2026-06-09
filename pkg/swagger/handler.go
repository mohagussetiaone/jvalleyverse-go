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
  <title>JValleyVerse API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
  <style>
    body { margin: 0; background: #fafafa; font-family: sans-serif; }
    #login-bar {
      display: flex; align-items: center; gap: 10px;
      padding: 10px 20px; background: #1b1b1b; color: #fff;
      flex-wrap: wrap;
    }
    #login-bar span { font-size: 14px; font-weight: bold; color: #61affe; }
    #login-bar input {
      padding: 6px 10px; border-radius: 4px; border: 1px solid #555;
      background: #333; color: #fff; font-size: 13px; width: 200px;
    }
    #login-bar button {
      padding: 6px 16px; border-radius: 4px; border: none;
      background: #61affe; color: #fff; font-size: 13px;
      cursor: pointer; font-weight: bold;
    }
    #login-bar button:hover { background: #4a90d9; }
    #login-status { font-size: 13px; }
    #login-status.ok { color: #49cc90; }
    #login-status.err { color: #f93e3e; }
  </style>
</head>
<body>

  <div id="login-bar">
    <span>🔐 Quick Login:</span>
    <input id="login-email" type="email" placeholder="Email" value="admin@jvalleyverse.com" />
    <input id="login-password" type="password" placeholder="Password" value="Admin@123" />
    <button onclick="doLogin()">Login &amp; Set Token</button>
    <span id="login-status"></span>
  </div>

  <div id="swagger-ui"></div>

  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-standalone-preset.js"></script>

  <script>
    var currentToken = localStorage.getItem('swagger_token') || '';

    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/api/docs/openapi.json",
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [SwaggerUIBundle.plugins.DownloadUrl],
        layout: "StandaloneLayout",
        persistAuthorization: true,
        requestInterceptor: function(req) {
          var token = localStorage.getItem('swagger_token');
          if (token) {
            req.headers['Authorization'] = 'Bearer ' + token;
          }
          return req;
        }
      });

      // Restore token status on page load
      if (currentToken) {
        setStatus('✔ Token loaded from storage', 'ok');
      }
    };

    function doLogin() {
      var email    = document.getElementById('login-email').value;
      var password = document.getElementById('login-password').value;
      var status   = document.getElementById('login-status');

      setStatus('Logging in...', '');

      fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: email, password: password })
      })
      .then(function(res) { return res.json(); })
      .then(function(data) {
        var token = data.token || data.access_token;
        if (!token) {
          setStatus('✘ Login failed: ' + (data.error || JSON.stringify(data)), 'err');
          return;
        }
        localStorage.setItem('swagger_token', token);
        currentToken = token;
        setStatus('✔ Logged in as ' + email, 'ok');

        // Also authorize Swagger UI natively so the lock icons show as locked
        window.ui.preauthorizeApiKey('BearerAuth', token);
      })
      .catch(function(err) {
        setStatus('✘ Error: ' + err.message, 'err');
      });
    }

    function setStatus(msg, cls) {
      var el = document.getElementById('login-status');
      el.textContent = msg;
      el.className = cls;
    }
  </script>
</body>
</html>`
}

// GetOpenAPISpec returns OpenAPI 3.0.0 specification
func GetOpenAPISpec() string {
	return openAPISpec
}
