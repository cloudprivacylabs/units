# Units

This is a RESTful front-end to the Unified Code for Units of Measure
(UCUM) library that implements the measure service functionality for
Layered Schema Architecture. It extends the UCUM interface with a
regular-expression based unit matching functionality to extend UCUM to
handle common nonstandard measure representations, such as 5'6" for
length.

The homepage for UCUM is: https://ucum.nlm.nih.gov/

## Building/Running

This service contains a NodeJS application (the UCUM implementation)
and a Go API front-end that deals with normalizations not supported by
UCUM. Easiest way to run it is using a Docker image. To build, run:

```
docker build -t units:latest .
```

The docker image runs two servers. The UCUM server is running on port
8080, and the Go front-end server is running on port 8081. Normally
you would want to expose the 8081 service:

To run:

```
docker run -p 9090:8081 units:latest
```

This will make the `units` service available on `localhost:9090`.

## The API

There are three endpoints provided by this service:

### Unit normalization

```
curl http://localhost:9090/unit?value=5'6"
```

This will return:

```
{
  "unit": "[in_i]",
  "valid": true,
  "value": "66"
}
```

This output means that the unit input is valid, its unit is `[in_i]`,
which is the UCUM code for inch, and its value is `66`, meaning the
value is `66 inch`.

The input:

```
curl http://localhost:9090/unit?value=12meter
```

The output is:

```
{
  "msg": "The UCUM code for meter is m.\nDid you mean m?",
  "unit": "m",
  "valid": false,
  "value": "12"
}
```

Note that the unit conversion is done, but it is marked as invalid.

### Unit Validation

This API directly calls the UCUM implementation to validate a unit.

```
curl http://localhost:9090/validate?unit=meter

{
  "status": "invalid",
  "ucumCode": "m",
  "unit": {
    "code": "m",
    "name": "meter",
    "guidance": "unit of length = 1.09361 yards"
  },
  "msg": [
    "The UCUM code for meter is m.\nDid you mean m?"
  ]
}
```

### Unit Conversion

This API directly calls the UCUM conversion API.

```
curl http://localhost:9090/convert?value=1&unit=cm&output=m

{
  "status": "succeeded",
  "toVal": 0.01,
  "msg": [
    
  ],
  "fromUnit": {
    "isBase_": true,
    "name_": "centimeter",
    "csCode_": "cm",
    "ciCode_": "CM",
    "property_": "length",
    "magnitude_": 0.01,
    "dim_": {
      "dimVec_": [
        1,
        0,
        0,
        0,
        0,
        0,
        0
      ]
    },
    "printSymbol_": "cm",
    "class_": null,
    "isMetric_": false,
    "variable_": "L",
    "cnv_": null,
    "cnvPfx_": 1,
    "isSpecial_": false,
    "isArbitrary_": false,
    "moleExp_": 0,
    "synonyms_": "centimeters; centimetres",
    "source_": "LOINC",
    "loincProperty_": "Len",
    "category_": "Clinical",
    "guidance_": null,
    "csUnitString_": null,
    "ciUnitString_": null,
    "baseFactorStr_": null,
    "baseFactor_": null,
    "defError_": false
  },
  "toUnit": {
    "isBase_": true,
    "name_": "meter",
    "csCode_": "m",
    "ciCode_": "M",
    "property_": "length",
    "magnitude_": 1,
    "dim_": {
      "dimVec_": [
        1,
        0,
        0,
        0,
        0,
        0,
        0
      ]
    },
    "printSymbol_": "m",
    "class_": null,
    "isMetric_": false,
    "variable_": "L",
    "cnv_": null,
    "cnvPfx_": 1,
    "isSpecial_": false,
    "isArbitrary_": false,
    "moleExp_": 0,
    "synonyms_": "meters; metres; distance",
    "source_": "UCUM",
    "loincProperty_": "Len",
    "category_": "Clinical",
    "guidance_": "unit of length = 1.09361 yards",
    "csUnitString_": null,
    "ciUnitString_": null,
    "baseFactorStr_": null,
    "baseFactor_": null,
    "defError_": false
  }
}
```
