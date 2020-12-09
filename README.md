# Classify

Classify is an efficient tool for extraction of structured field sequences from HTML/XML data sources. It finds repetitive patterns and returns sequence of fields or XPaths to extract those fields.

## When is this useful

Classify allows to create scraping/parsing pattern for specific data source and do it quickly.

Sample HTML:
```
<html>
    <body>
        <section> Some Ad </section>
        <section> 
            <h1> Data </h1> 
            <div>
                <h3> Title 1 </h3>
                <p> Some text 1 </p>
                <img src="/src1"> </img>
            </div>
            <div>
                <h3> Title 2 </h3>
                <p> Some text 2 </p>
                <img src="/src1"> </img>
            </div>
            <div>
                <h3> Title 3 </h3>
                <p> Some text 3 </p>
                <img src="/src1"> </img>
            </div>
        </section>
        <section> 
            <h2> Some Menu </h2>
            <ul>
                <li>Item 1</li>
                <li>Item 2</li>
                <li>Item 3</li>
            </ul>
        </section>
    </body>
</html>
```

Sample text output after running:

```fields```
```
Title 1
Some text 1
image1
--------------------------
Title 2
Some text 2
image2
--------------------------
Title 3
Some text 3
image3
--------------------------
```

```fields -json```
```
{
  "fields": [
    [
      "Title 1",
      "Some text 1",
      "image1"
    ],
    [
      "Title 2",
      "Some text 2",
      "image2"
    ],
    [
      "Title 3",
      "Some text 3",
      "image3"
    ]
  ],
  "stats": {
    "groups_count": 3,
    "group_fields_count": 3
  }
}
```

```fields -csv```
```
Title 1,Some text 1,image1
Title 2,Some text 2,image2
Title 3,Some text 3,image3
```

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