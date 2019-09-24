Get customized placeholder images for your mockups and website designs :rocket:


## Features
 * Free
 * Open Source
 * Fast (delivered via a CDN)
 
## Usage
The base URL is: `https://nocontent.xyz/img`

To customize the image you will get, prodive one or more of the following options as query string parameters:

|name|description||value|
|---|---|---|---|
|x|width of image|optional (default: 200)|number (in pixels)|
|y|height of image|optional (default: 200)|number (in pixels)|
|bg|background color|optional (default: ffffff)|hex-value (3 or 6 character)|
|fg|text color|optional (default: 333333)|hex-value (3 or 6 character)|
|label|text on the image|optional (default: 'height' x 'width') | string (max-length: 20)

## Example #1
Without any parameters you will get a small white image with black text and 200x200 px
```
https://nocontent.xyz/img
```
![200x200](https://nocontent.xyz/img)

## Example #2
To get a purple image with the dimensions 600x400, query this address:
```
https://nocontent.xyz/img?x=600&y=400&bg=980080
```
![600x400](https://nocontent.xyz/img?x=600&y=400&bg=980080)

## Example #3
Change the text color to blue and provide a custom text with this URL:
```
https://nocontent.xyz/img?label=Hello%20World&fg=0099FF
```
![200x200](https://nocontent.xyz/img?label=Hello%20World&fg=0099FF)



:computer: with :heart: by [codemonauts](https://codemonauts.com)
