#!/usr/bin/env python3
import os
import json
from typing import Optional, List
from urllib.parse import urlencode, urlparse
from fastapi import FastAPI, Request, Form, File, UploadFile, Query
from fastapi.responses import HTMLResponse, FileResponse, StreamingResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from weasyprint import HTML, CSS
from io import BytesIO
import requests
import base64

# Create FastAPI app
app = FastAPI(
    title="Nametag Generator",
    description="Generate professionally designed event nametags with perfect PDF output",
    version="1.0.0",
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Setup templates
templates = Jinja2Templates(directory="templates")

# Configuration
DEFAULT_TEMPLATE = "simple_horizontal"

BASE_CSS = """
@page {
  size: 6.3cm 8.2cm;
  margin: 0;
}

:root {
  --brand: #7e287b;
  --brand-2: #231f54;
  --brand-3: #663c96;
  --text: #111;
  --muted: #555;
  --border: #d9d9d9;
}

body {
  margin: 0;
  padding: 0;
  font-family: Inter, system-ui, -apple-system, "Segoe UI", Roboto,
    "Helvetica Neue", Arial, "Noto Sans", "Liberation Sans", sans-serif;
  color: var(--text);
}

.nametag {
  width: 6.3cm;
  height: 8.2cm;
  box-sizing: border-box;
  margin: 0;
  padding: 0.5cm;
  background: #fff;
  position: relative;
}

.logo-wrap {
  text-align: center;
  margin-bottom: 0.3cm;
}

.logo-img {
  height: 1.1cm;
  max-width: 100%;
  display: inline-block;
}

.event-name {
  color: var(--brand);
  font-weight: bold;
}

.role-badge {
  color: var(--brand);
  border-color: var(--brand);
}

.divider {
  border-top-color: var(--brand);
}

.date-location .location {
  color: var(--brand);
}

.receiver-name {
  font-size: 26pt;
  font-weight: 900;
  text-transform: uppercase;
  line-height: 1.02;
  margin: 0.5cm 0;
}
"""


# Define API models
class NametagRequest(BaseModel):
    template: str = DEFAULT_TEMPLATE
    logo_url: Optional[str] = None
    event_name: str
    first_name: str
    last_name: str
    role: str
    dates: Optional[str] = None
    location: Optional[str] = None

    class Config:
        schema_extra = {
            "example": {
                "template": "simple_horizontal",
                "logo_url": "https://example.com/logo.png",
                "event_name": "Annual Conference 2023",
                "first_name": "Jane",
                "last_name": "Doe",
                "role": "Speaker",
                "dates": "21 - 22 August 2023",
                "location": "Nairobi, Kenya",
            }
        }


def get_available_templates():
    """Get a list of available templates in the templates directory."""
    templates_dir = "templates"
    templates = []

    for file in os.listdir(templates_dir):
        if (
            file.endswith(".html")
            and not file.startswith("_")
            and file not in ["index.html", "preview.html"]
        ):
            template_name = file[:-5]  # Remove .html extension
            templates.append(template_name)

    return templates or [DEFAULT_TEMPLATE]


def get_image_as_base64(url):
    """Convert an image URL to a base64 data URI for embedding in HTML."""
    try:
        if not url:
            return None

        # Check if it's already a data URI
        if url.startswith("data:"):
            return url

        # Check if it's a valid URL
        parsed = urlparse(url)
        if not parsed.scheme or not parsed.netloc:
            return None

        response = requests.get(url, timeout=5)
        response.raise_for_status()

        # Get content type from response
        content_type = response.headers.get("Content-Type", "image/png")

        # Convert to base64
        img_data = base64.b64encode(response.content).decode("utf-8")
        return f"data:{content_type};base64,{img_data}"
    except Exception as e:
        print(f"Error fetching image from {url}: {e}")
        return None


@app.get("/", response_class=HTMLResponse)
async def home(request: Request):
    """Render home page with form."""
    available_templates = get_available_templates()
    return templates.TemplateResponse(
        "index.html", {"request": request, "templates": available_templates}
    )


@app.post("/preview", response_class=HTMLResponse)
async def preview(
    request: Request,
    template: str = Form(DEFAULT_TEMPLATE),
    logo_url: Optional[str] = Form(None),
    event_name: str = Form(...),
    first_name: str = Form(...),
    last_name: str = Form(...),
    role: str = Form(...),
    dates: Optional[str] = Form(None),
    location: Optional[str] = Form(None),
):
    """Generate a preview of the nametag."""
    print(f"Original logo URL: {logo_url}")

    # Try to get the logo as base64 if a URL is provided
    logo_data_uri = get_image_as_base64(logo_url) if logo_url else None

    print(f"Converted to data URI: {'Yes' if logo_data_uri else 'No'}")

    data = {
        "logo_url": logo_data_uri or logo_url or "",
        "event_name": event_name,
        "first_name": first_name,
        "last_name": last_name,
        "role": role,
        "dates": dates or "",
        "location": location or "",
    }

    # Create query string for the PDF URL
    pdf_params = urlencode({k: v for k, v in data.items() if v})
    pdf_params = f"template={template}&{pdf_params}"

    # Render the nametag template
    html_content = templates.get_template(f"{template}.html").render(**data)

    # Return the preview template with the nametag content
    return templates.TemplateResponse(
        "preview.html",
        {
            "request": request,
            "html_content": html_content,
            "pdf_url": f"/generate-pdf?{pdf_params}",
        },
    )


@app.get("/generate-pdf")
async def generate_pdf(
    template: str = Query(DEFAULT_TEMPLATE),
    logo_url: Optional[str] = Query(None),
    event_name: str = Query(...),
    first_name: str = Query(...),
    last_name: str = Query(...),
    role: str = Query(...),
    dates: Optional[str] = Query(None),
    location: Optional[str] = Query(None),
):
    """Generate a printable PDF version of the nametag."""
    # Try to get the logo as base64 if a URL is provided
    logo_data_uri = get_image_as_base64(logo_url) if logo_url else None

    data = {
        "logo_url": logo_data_uri or logo_url or "",
        "event_name": event_name,
        "first_name": first_name,
        "last_name": last_name,
        "role": role,
        "dates": dates or "",
        "location": location or "",
    }

    # Render nametag HTML content
    html_content = templates.get_template(f"{template}.html").render(**data)

    # Generate PDF with WeasyPrint
    pdf_buffer = BytesIO()
    HTML(string=html_content).write_pdf(pdf_buffer, stylesheets=[CSS(string=BASE_CSS)])
    pdf_buffer.seek(0)

    # Generate a filename based on the person's name
    filename = f"{first_name}-{last_name}-nametag.pdf".replace(" ", "-").lower()

    # Return the PDF as a downloadable file
    return StreamingResponse(
        pdf_buffer,
        media_type="application/pdf",
        headers={"Content-Disposition": f"attachment; filename={filename}"},
    )


@app.post("/api/generate")
async def api_generate(nametag: NametagRequest):
    """API endpoint for generating PDFs programmatically."""
    # Try to get the logo as base64 if a URL is provided
    logo_data_uri = get_image_as_base64(nametag.logo_url) if nametag.logo_url else None

    data = {
        "logo_url": logo_data_uri or nametag.logo_url or "",
        "event_name": nametag.event_name,
        "first_name": nametag.first_name,
        "last_name": nametag.last_name,
        "role": nametag.role,
        "dates": nametag.dates or "",
        "location": nametag.location or "",
    }

    # Render nametag HTML content
    html_content = templates.get_template(f"{nametag.template}.html").render(**data)

    # Generate PDF with WeasyPrint
    pdf_buffer = BytesIO()
    HTML(string=html_content).write_pdf(pdf_buffer, stylesheets=[CSS(string=BASE_CSS)])
    pdf_buffer.seek(0)

    # Generate a filename based on the person's name
    filename = f"{nametag.first_name}-{nametag.last_name}-nametag.pdf".replace(
        " ", "-"
    ).lower()

    # Return the PDF as a downloadable file
    return StreamingResponse(
        pdf_buffer,
        media_type="application/pdf",
        headers={"Content-Disposition": f"attachment; filename={filename}"},
    )


if __name__ == "__main__":
    import uvicorn

    port = int(os.environ.get("PORT", 8000))
    uvicorn.run("app:app", host="0.0.0.0", port=port, reload=True)
