# Classify

Classify is an efficient tool for extraction of structured field sequences from HTML/XML data sources. It finds repetitive patterns and returns sequence of fields or XPaths to extract those fields.

## When is this useful

Classify allows to create scraping/parsing solution for specific data source and do it quickly.

## Requirements

Go 1.13+

## Installation

```
git clone https://github.com/olesho/classify
cd bin/fields
go install
```

## Usage

Input data could be any URL or file.

### Extracting values in JSON:
```curl -s YOUR_URL_HERE | fields```
or
```curl -s YOUR_URL_HERE | fields -json```

### Extracting XPaths:
```curl -s YOUR_URL_HERE | fields -xpath```

Keep in mind that some websites use dynamic content extensively. So CURLed version might differ significantly from the one you see in browser. You might want to use: ```chromium --dump-dom 'YOUR_URL_HERE' | fields``` or ```google-chrome --dump-dom 'YOUR_URL_HERE' | fields``` instead of CURL 