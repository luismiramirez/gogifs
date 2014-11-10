gogifs
======

## Small Go webapp with Martini

Easy as a pie, two endpoints.

### To upload reactions

```
POST /reactions

Params: title, image
```

### To get a random reaction

```
GET /randomreaction
```

Always using the header ```API-KEY``` that has to be set via environment variable in the host

## TODO

* Tests
* Apply better practices as this is my first Go application
