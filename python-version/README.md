# Nametag Generator (Python Version)

A modern nametag generator built with FastAPI and WeasyPrint. This application allows you to create professional-looking event nametags and generate perfectly sized PDFs for printing.

## Features

- Generate nametags with custom content (name, role, event details)
- Preview nametags in browser
- Download perfectly-sized PDFs (6.3cm × 8.2cm) for printing
- API endpoint for programmatic access

## Setup

### Prerequisites

- Python 3.8+
- Cairo graphics library (required by WeasyPrint)

### Installation

1. Install Cairo graphics library:

   ```
   # For Ubuntu/Debian
   sudo apt-get install libcairo2-dev libpango1.0-dev libgdk-pixbuf2.0-dev libffi-dev shared-mime-info

   # For macOS
   brew install cairo pango gdk-pixbuf libffi
   ```

2. Create a virtual environment and activate it:

   ```
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```

3. Install the dependencies:

   ```
   pip install -r requirements.txt
   ```

4. Run the application:

   ```
   uvicorn app:app --reload
   ```

5. Open your browser and navigate to [http://127.0.0.1:8000](http://127.0.0.1:8000)

## API Usage

FastAPI provides automatic interactive API documentation at `/docs` and `/redoc` endpoints.

You can generate nametag PDFs programmatically via the API:

```bash
curl -X POST http://127.0.0.1:8000/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "template": "simple_horizontal",
    "event_name": "Annual Conference 2023",
    "first_name": "Jane",
    "last_name": "Doe",
    "role": "Speaker",
    "dates": "21 - 22 August 2023",
    "location": "Nairobi, Kenya"
  }' \
  --output nametag.pdf
```

## Creating New Templates

To create a new template:

1. Add a new HTML file in the `templates` directory with a descriptive name, e.g., `vertical_design.html`
2. Ensure the template includes styling for the nametag and responsive handling of all data fields
3. Your new template will automatically appear in the template selection dropdown

## Deployment

To deploy to production:

```
uvicorn app:app --host 0.0.0.0 --port 8000 --workers 4
```

Or use Gunicorn with Uvicorn workers:

```
gunicorn -w 4 -k uvicorn.workers.UvicornWorker app:app
```

## WeasyPrint Note

WeasyPrint is a powerful HTML/CSS to PDF converter, but it has some limitations compared to modern browsers:

1. **Limited flexbox support**: Our templates use traditional CSS layouts instead of flexbox
2. **No CSS Grid support**: Avoid using display: grid in templates
3. **Limited JavaScript support**: PDF generation is static, no JS execution

When creating new templates, use traditional CSS positioning, tables, and block/inline elements for best results.
