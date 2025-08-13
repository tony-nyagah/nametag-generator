package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//go:embed templates/*
var templateFS embed.FS

// Home page template
var indexTmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Event Nametag Generator</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;800&family=Space+Grotesk:wght@400;700&display=swap" rel="stylesheet">
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
    <style>
        :root {
            --primary: #FF5470;
            --secondary: #5FB3FC;
            --accent: #FFBF47;
            --dark: #232323;
            --light: #F7F7F9;
            --success: #36D39A;
            --shadow-offset: 4px;
        }
        
        * {
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Inter', system-ui, sans-serif;
            background-color: var(--light);
            line-height: 1.6;
            color: var(--dark);
            max-width: 1000px;
            margin: 0 auto;
            padding: 40px 20px;
        }
        
        h1, h2, h3, h4 {
            font-family: 'Space Grotesk', sans-serif;
            font-weight: 700;
        }
        
        h1 {
            font-size: 3rem;
            margin-bottom: 1rem;
            position: relative;
            display: inline-block;
        }
        
        h1::after {
            content: "";
            position: absolute;
            width: 100%;
            height: 12px;
            bottom: 8px;
            left: 0;
            z-index: -1;
            background-color: var(--accent);
        }
        
        .container {
            display: grid;
            grid-template-columns: 1.5fr 1fr;
            gap: 40px;
        }
        
        @media (max-width: 900px) {
            .container {
                grid-template-columns: 1fr;
            }
        }
        
        .card {
            background: white;
            border: 3px solid var(--dark);
            border-radius: 8px;
            padding: 30px;
            box-shadow: var(--shadow-offset) var(--shadow-offset) 0 var(--dark);
            margin-bottom: 30px;
        }
        
        .card-title {
            background: var(--secondary);
            color: var(--dark);
            margin: -30px -30px 20px -30px;
            padding: 20px 30px;
            border-bottom: 3px solid var(--dark);
            font-weight: 800;
            font-size: 1.4rem;
        }
        
        form {
            display: grid;
            gap: 20px;
        }
        
        .field {
            display: grid;
            gap: 8px;
        }
        
        label {
            font-weight: 600;
            font-size: 0.95rem;
        }
        
        input, select {
            font-family: 'Inter', system-ui, sans-serif;
            font-size: 1rem;
            padding: 12px 16px;
            border: 3px solid var(--dark);
            border-radius: 6px;
            background: white;
            box-shadow: 2px 2px 0 var(--dark);
            transition: transform 0.1s, box-shadow 0.1s;
        }
        
        input:focus, select:focus {
            outline: none;
            box-shadow: 4px 4px 0 var(--secondary);
            transform: translate(-2px, -2px);
        }
        
        input::placeholder {
            color: #aaa;
        }
        
        input.error {
            border-color: var(--primary);
            background-color: rgba(255, 84, 112, 0.05);
            animation: shake 0.3s;
        }
        
        @keyframes shake {
            0%, 100% { transform: translateX(0); }
            25% { transform: translateX(-8px); }
            75% { transform: translateX(8px); }
        }
        
        button {
            font-family: 'Space Grotesk', sans-serif;
            font-weight: 700;
            font-size: 1.1rem;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            padding: 16px;
            background: var(--primary);
            color: white;
            border: 3px solid var(--dark);
            border-radius: 6px;
            cursor: pointer;
            box-shadow: var(--shadow-offset) var(--shadow-offset) 0 var(--dark);
            transition: transform 0.1s ease, box-shadow 0.1s ease;
        }
        
        button:hover {
            transform: translate(-2px, -2px);
            box-shadow: calc(var(--shadow-offset) + 2px) calc(var(--shadow-offset) + 2px) 0 var(--dark);
        }
        
        button:active {
            transform: translate(2px, 2px);
            box-shadow: 1px 1px 0 var(--dark);
        }
        
        .api-section {
            background: var(--dark);
            color: white;
            padding: 20px;
            border-radius: 8px;
            overflow: auto;
        }
        
        .api-section h3 {
            color: var(--accent);
            margin-top: 0;
        }
        
        pre {
            background: rgba(255,255,255,0.1);
            padding: 15px;
            border-radius: 4px;
            overflow: auto;
            font-family: 'Courier New', monospace;
        }
        
        .template-preview {
            border: 3px solid var(--dark);
            border-radius: 8px;
            overflow: hidden;
            box-shadow: var(--shadow-offset) var(--shadow-offset) 0 var(--dark);
            background: white;
            height: 450px;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
			margin-bottom: 30px;
        }
        
        .preview-placeholder {
            text-align: center;
            color: #aaa;
        }
        
        .preview-placeholder svg {
            width: 80px;
            height: 80px;
            margin-bottom: 20px;
            opacity: 0.5;
        }
        
        .htmx-indicator {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
        }
        
        .loading-spinner {
            width: 38px;
            height: 38px;
        }
        
        .actions {
            display: flex;
            gap: 15px;
        }
    </style>
</head>
<body>
    <h1>Nametag Generator</h1>
    
    <div class="container">
        <div>
            <div class="card">
                <div class="card-title">Create Your Nametag</div>
                <form id="nametag-form" hx-post="/generate" hx-target="#preview-content" hx-trigger="submit, change from:#template" hx-indicator=".htmx-indicator" hx-swap="innerHTML">
                    <div class="field">
                        <label for="template">Choose Template:</label>
                        <select name="template" id="template" required>
                            {{range .Templates}}
                            <option value="{{.}}">{{.}}</option>
                            {{end}}
                        </select>
                    </div>
                    
                    <div class="field">
                        <label for="logoUrl">Logo URL:</label>
                        <input type="url" name="logoUrl" id="logoUrl" placeholder="https://example.com/logo.png">
                    </div>
                    
                    <div class="field">
                        <label for="eventName">Event Name:</label>
                        <input type="text" name="eventName" id="eventName" required>
                    </div>
                    
                    <div class="field">
                        <label for="firstName">First Name:</label>
                        <input type="text" name="firstName" id="firstName" required>
                    </div>
                    
                    <div class="field">
                        <label for="lastName">Last Name:</label>
                        <input type="text" name="lastName" id="lastName" required>
                    </div>
                    
                    <div class="field">
                        <label for="role">Role:</label>
                        <input type="text" name="role" id="role" placeholder="Attendee, Speaker, Sponsor, etc." required>
                    </div>
                    
                    <div class="field">
                        <label for="dates">Dates:</label>
                        <input type="text" name="dates" id="dates" placeholder="21 - 22 August 2025">
                    </div>
                    
                    <div class="field">
                        <label for="location">Location:</label>
                        <input type="text" name="location" id="location" placeholder="City, Country">
                    </div>
                    
                    <div class="actions">
                        <button type="submit" style="background: var(--secondary);">Preview</button>
                        <a href="/generate" target="_blank" class="printable-btn" style="display:none;" id="print-link">
                            <button type="button">Generate Printable</button>
                        </a>
                    </div>
                </form>
            </div>
        </div>
        
        <div>
            <div class="template-preview">
                <div id="preview-content">
                    <div class="preview-placeholder">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                            <circle cx="8.5" cy="8.5" r="1.5"></circle>
                            <polyline points="21 15 16 10 5 21"></polyline>
                        </svg>
                        <p>Fill in the form and click "Preview"<br>to see your nametag</p>
                    </div>
                </div>
                <div class="htmx-indicator" style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%);">
                    <div class="loading-spinner">
                        <svg width="38" height="38" viewBox="0 0 38 38" xmlns="http://www.w3.org/2000/svg" stroke="var(--secondary)">
                            <g fill="none" fill-rule="evenodd">
                                <g transform="translate(1 1)" stroke-width="2">
                                    <circle stroke-opacity=".5" cx="18" cy="18" r="18"/>
                                    <path d="M36 18c0-9.94-8.06-18-18-18">
                                        <animateTransform attributeName="transform" type="rotate" from="0 18 18" to="360 18 18" dur="1s" repeatCount="indefinite"/>
                                    </path>
                                </g>
                            </g>
                        </svg>
                    </div>
                </div>
            </div>
            
            <div class="card api-section">
                <h3>API Usage</h3>
                <pre>POST /api/generate
Content-Type: application/json

{
  "template": "Simple Horizontal Template",
  "logoUrl": "https://example.com/logo.png", 
  "eventName": "Annual Conference 2025",
  "firstName": "Jane",
  "lastName": "Doe",
  "role": "Speaker",
  "dates": "21 - 22 August 2025",
  "location": "Nairobi, Kenya"
}</pre>
            </div>
        </div>
    </div>
    
    <!-- HTMX handles the interaction -->
</body>
</html>
`))

type NametagData struct {
	LogoURL   string
	EventName string
	FirstName string
	LastName  string
	Role      string
	Dates     string
	Location  string
}

func getAvailableTemplates() []string {
	var templates []string

	entries, err := templateFS.ReadDir("templates")
	if err != nil {
		log.Printf("Error reading embedded templates: %v", err)
		return []string{"classic"}
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".html") {
			name := strings.TrimSuffix(entry.Name(), ".html")
			templates = append(templates, name)
		}
	}

	if len(templates) == 0 {
		return []string{"classic"}
	}

	return templates
}

func getTemplate(name string) *template.Template {
	templatePath := "templates/" + name + ".html"

	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		log.Printf("Template %s not found, falling back to classic: %v", name, err)
		// Fall back to classic template if file doesn't exist
		content, err = templateFS.ReadFile("templates/Simple Horizontal Template.html")
		if err != nil {
			log.Fatalf("Critical error: default template not found: %v", err)
		}
	}

	return template.Must(template.New(name).Parse(string(content)))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Templates": getAvailableTemplates(),
	}

	if err := indexTmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering index template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	// Support both GET (preview) and POST (generate)
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	templateName := r.FormValue("template")
	if templateName == "" {
		templateName = "Simple Horizontal Template"
	}

	// Prepare nametag data
	data := NametagData{
		LogoURL:   r.FormValue("logoUrl"),
		EventName: r.FormValue("eventName"),
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		Role:      r.FormValue("role"),
		Dates:     r.FormValue("dates"),
		Location:  r.FormValue("location"),
	}

	// Get and execute template
	tmpl := getTemplate(templateName)
	w.Header().Set("Content-Type", "text/html")

	// Check if this is an HTMX request or a direct page load
	isHtmx := r.Header.Get("HX-Request") == "true"
	isPrintable := r.URL.Query().Get("printable") == "true"

	if isHtmx {
		// HTMX request - render template with print button
		// First, generate the nametag
		var nametagHTML bytes.Buffer
		if err := tmpl.Execute(&nametagHTML, data); err != nil {
			log.Printf("Error rendering nametag template %s: %v", templateName, err)
			http.Error(w, "Error generating nametag", http.StatusInternalServerError)
			return
		}

		// Create a wrapper with the nametag and print controls
		printURL := fmt.Sprintf("/generate?template=%s&logoUrl=%s&eventName=%s&firstName=%s&lastName=%s&role=%s&dates=%s&location=%s&printable=true",
			url.QueryEscape(templateName),
			url.QueryEscape(data.LogoURL),
			url.QueryEscape(data.EventName),
			url.QueryEscape(data.FirstName),
			url.QueryEscape(data.LastName),
			url.QueryEscape(data.Role),
			url.QueryEscape(data.Dates),
			url.QueryEscape(data.Location),
		)

		// Render the preview with print controls
		fmt.Fprintf(w, `<div class="preview-nametag">%s</div>
			<div class="preview-controls" style="position:absolute; bottom:15px; left:0; right:0; text-align:center;">
				<a href="%s" target="_blank">
					<button style="background: var(--success); width: auto; padding: 10px 20px; display:inline-flex; align-items:center; gap:8px;">
						<svg style="width: 16px; height: 16px;" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<polyline points="6 9 6 2 18 2 18 9"></polyline>
							<path d="M6 18H4a2 2 0 01-2-2v-5a2 2 0 012-2h16a2 2 0 012 2v5a2 2 0 01-2 2h-2"></path>
							<rect x="6" y="14" width="12" height="8"></rect>
						</svg>
						Print Version
					</button>
				</a>
			</div>`, nametagHTML.String(), printURL)

		// Also update the form's print button URL
		fmt.Fprintf(w, `<script>
			document.getElementById("print-link").href = "%s";
			document.getElementById("print-link").style.display = "block";
		</script>`, printURL)
	} else if isPrintable {
		// Print-ready version with minimal page setup
		fmt.Fprintf(w, `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="utf-8">
			<title>Print: %s %s - %s</title>
			<style>
				@page { size: 6.3cm 8.2cm; margin: 0; }
				body { margin: 0; display: flex; justify-content: center; align-items: center; }
				@media print {
					* { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
				}
			</style>
		</head>
		<body>`, data.FirstName, data.LastName, data.Role)

		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Error rendering printable template %s: %v", templateName, err)
			http.Error(w, "Error generating printable nametag", http.StatusInternalServerError)
			return
		}

		// Add auto-print script
		fmt.Fprintf(w, `
		<script>
			window.onload = function() {
				window.print();
			}
		</script>
		</body>
		</html>`)
	} else {
		// Direct GET request, just render the template
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Error rendering nametag template %s: %v", templateName, err)
			http.Error(w, "Error generating nametag", http.StatusInternalServerError)
		}
	}
}

func apiGenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var request struct {
		Template  string `json:"template"`
		LogoURL   string `json:"logoUrl"`
		EventName string `json:"eventName"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
		Dates     string `json:"dates"`
		Location  string `json:"location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Set default template if not specified
	templateName := request.Template
	if templateName == "" {
		templateName = "Simple Horizontal Template"
	}

	// Prepare nametag data
	data := NametagData{
		LogoURL:   request.LogoURL,
		EventName: request.EventName,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Role:      request.Role,
		Dates:     request.Dates,
		Location:  request.Location,
	}

	// Get and execute template
	tmpl := getTemplate(templateName)
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("API: Error rendering nametag template %s: %v", templateName, err)
		http.Error(w, "Error generating nametag", http.StatusInternalServerError)
	}
}

func main() {
	// Register HTTP handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/api/generate", apiGenerateHandler)

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	serverAddr := ":" + port
	log.Printf("Server starting on http://localhost%s", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
