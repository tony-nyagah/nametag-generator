# Event Nametag Generator

A small, print‑ready system for generating event nametags from HTML/CSS templates. It injects event details (logo, names, role, dates, location) into a chosen template and outputs a page at the exact card size.

Note: brand‑agnostic. Templates define their own colors; you can supply any logo/brand styling.

## Features

- HTML/CSS templates sized for printing (6.3cm × 8.2cm)
- Print‑optimized (@page sizing, print color adjust, border‑based dividers)
- Simple variable injection via element IDs and attributes
- Role badge with role‑based color variants (template‑specific)
- Frontend UI (planned) for single badge entry and CSV/Excel batch generation

## Templates

- Location: templates/
- Example: templates/Simple Horizontal Template.html
- Pick a template, inject variables, print at 100% scale.

## Template variables (convention)

Provide these values when rendering a template:

- logoUrl: URL or path to the event/organization logo
- eventName: e.g., "24th Annual Scientific Conference"
- firstName: attendee’s first name
- lastName: attendee’s last name
- role: Attendee | Speaker | Sponsor | Staff (string)
- date: e.g., "21 – 22 August 2025"
- location: e.g., "Nairobi, Kenya"

DOM IDs and attributes (used by current templates):

- #first-name, #last-name, #role, #date, #location
- .event-name (text)
- .logo-img (set src to logoUrl)

Some templates support role color classes (e.g., role-attendee, role-speaker, role-sponsor). Toggle or set the class on the role element to match the role.

## Frontend (planned)

- Single badge form: fields for logo URL/upload, event name, first/last name, role, date, location, template selector, live preview
- Batch mode: upload CSV/Excel, map columns to fields, preview first N, generate all
- Output: per‑badge HTML/PDF and ZIP download; optional multi‑up imposition for A4/Letter
- File handling: drag‑and‑drop, basic validation, sample file download

## Batch input formats

- CSV: headers like firstName,lastName,role,date,location,eventName,logoUrl
- Excel (.xlsx): first row as headers; multiple sheets optional (select on upload)
- Column mapper: UI to map arbitrary headers to required fields

## API sketch

- GET /templates → list available templates (name, path, preview)
- POST /render → body: { template, data } → returns rendered HTML/PDF
- POST /batch → body: { template, rows[] } or CSV/XLSX multipart → returns ZIP of outputs

## Suggested stack

- Backend: Go (net/http), html/template or dom manipulation, encoding/csv, github.com/xuri/excelize/v2 for .xlsx, archiver/zip for zips
- Frontend: lightweight (vanilla/HTMX) or a small framework (Svelte/React) with a preview iframe
- Rendering: browser print dialog at 100%, or headless Chromium for server PDFs if needed

## Printing tips

- Set scale to 100% (no fit‑to‑page)
- Disable headers/footers in the print dialog
- Paper size: use the template’s @page size (6.3cm × 8.2cm)
- Colors: exact printing enabled via print‑color‑adjust; dividers use borders for reliability

## Roadmap

- Generate organization/user IDs on badges (text, QR/Barcode)
- Bulk generation from CSV/JSON/XLSX with ZIP exports
- Multi‑up imposition on A4/Letter sheets
- Saved presets/branding profiles per event
- CLI and simple web UI

## Contributing

- Add new templates under templates/
- Follow the variable ID conventions above
- Keep print size and color optimization in mind
