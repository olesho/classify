# Classify

Classify is an efficient tool for extraction of structured field sequences from HTML/XML data sources. It finds repetitive patterns and returns sequence of fields or XPaths to extract those fields.

## When is this useful

Classify allows to create scraping/parsing pattern for specific data source and do it quickly. See usage section below.

## Requirements

Go 1.13+

## Installation

```
go get github.com/olesho/classify/sequence
cd classify/bin/fields
go install
```

## Usage

Sample HTML input:
```
<html>
    <body>
        <section> Some Ad </section>
        <section> 
            <h1> Data </h1> 
            <div>
                <h3> Title 1 </h3>
                <p> Some text 1 </p>
                <img src="/src1"> image1 </img>
            </div>
            <div>
                <h3> Title 2 </h3>
                <p> Some text 2 </p>
                <img src="/src2"> image2 </img>
            </div>
            <div>
                <h3> Title 3 </h3>
                <p> Some text 3 </p>
                <img src="/src3"> image3 </img>
            </div>
        </section>
        <section> 
            <h2> Some Menu </h2>
            <ul>
                <li>Menu Item 1</li>
                <li>Menu Item 2</li>
                <li>Menu Item 3</li>
            </ul>
        </section>
    </body>
</html>
```

### As text:
Run:
```fields```
```
Title 1
Some text 1
/src1
image1
--------------------------
Title 2
Some text 2
/src2
image2
--------------------------
Title 3
Some text 3
/src3
image3
--------------------------
```

### As JSON:
```fields -json```
```
{
  "fields": [
    [
      "Title 1",
      "Some text 1",
      "/src1",
      "image1"
    ],
    [
      "Title 2",
      "Some text 2",
      "/src2",
      "image2"
    ],
    [
      "Title 3",
      "Some text 3",
      "/src3",
      "image3"
    ]
  ],
  "stats": {
    "groups_count": 3,
    "group_fields_count": 4
  }
}
```

### As CSV:
```fields -csv```
```
Title 1,Some text 1,/src1,image1
Title 2,Some text 2,/src2,image2
Title 3,Some text 3,/src3,image3
```

### As XPath pattern
XPath pattern generation is still being developed:
```fields -xpath ```
```
/html/body/section/div/h3/@text=
/html/body/section/div/p/@text=
/html/body/section/div/img
/html/body/section/div/@text=
```

### Multiple groups of data:
Data is being extracted and ranked descending by volume with indexes 0, 1, 2 ... In case you need to access other groups put the index N like: ```fields N```

From example above ```fields -csv 1``` will produce following output:
```
Menu Item 1
Menu Item 2
Menu Item 3
``` 
Menu items take less volume and are ranked after main data block.

### Extracting from other sources:
Input data could be any URL or file.
```curl -s YOUR_URL_HERE | fields ```
Keep in mind that some websites use dynamic content extensively. So CURLed version might differ significantly from the one you see in the browser. You might want to use: ```chromium --dump-dom 'YOUR_URL_HERE' | fields``` or ```google-chrome --dump-dom 'YOUR_URL_HERE' | fields``` instead of CURL.

### Shell
For debugging purposes there is a command line tool:
```
cd bin/cshell
go build .
./cshell
```