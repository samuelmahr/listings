Write an API endpoint that returns a filtered set of listings from the data provided in `listing-details.csv`

API:
```
GET /listings?min_price=100000&max_price=200000&min_bed=2&max_bed=2&min_bath=2&max_bath=2
```


| Parameter    | Description                           |
|--------------|---------------------------------------|
| `min_price`  | The minimum listing price in dollars. |
| `max_price`  | The maximum listing price in dollars. |
| `min_bed`    | The minimum number of bedrooms.       |
| `max_bed`    | The maximum number of bedrooms.       |
| `min_bath`   | The minimum number of bathrooms.      |
| `max_bath`   | The maximum number of bathrooms.      |


The expected response is a GeoJSON FeatureCollection of listings:

```json
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [
          -112.1,
          33.4
        ]
      },
      "properties": {
        "id": "123ABC",
        "price": 200000,
        "street": "123 Walnut St",
        "bedrooms": 3,
        "bathrooms": 2,
        "sq_ft": 1500
      }
    }
  ]
}
```

All query parameters are optional, all minimum and maximum fields should be inclusive (e.g. min_bed=2&max_bed=4 should return listings with 2, 3, or 4 bedrooms).

At a minimum:
- Your API endpoint URL is `/listings`
- Your API responds with valid GeoJSON (you can check the output using http://geojson.io)
- Your API should correctly filter any combination of API parameters
- Use a datastore
