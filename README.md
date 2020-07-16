# Aeroqual AQY 1 Connector
This activity pulls data from Aeroqual AQY1 instrument.

## Installation
### Flogo CLI
```bash
flogo install github.com/matt-doug-davidson/aeroqualaqy1
```
## Settings
| Setting     | Type   | Required  | Description |
|:------------|:-------|:----------|:------------|
| host  | string      | True | The host running the MQTT broker|
| port | string | True | The MQTT port (typically 1883)|
| username | string | True | Username for log-in to the AQY1 instrument |
| password | string | True | Password for log-in to the AQY1 instrument |
| instrument | string | True | Instrument identification.Note: instrument ID can have spaces. |
| entity | string | True | Entity to which measurements are related. Exported identifier to other subsystems |
| mappings | string | True | JSON format that contains mapping between instrument key and output field. See example below. When entered in  Flogo Web UI don't need quotes. |

mappings example:
```json
{
   "NO2": {
     "field": "NO2"
   },
   "O3": {
     "field": "O3"
   },
   "TEMP": {
     "field": "Temperature"
   },
   "PM2.5": {
     "field": "PM2_5"
   },
   "PM10": {
     "field": "PM10"
   },
   "RH": {
     "field": "RelativeHumidity"
   },
   "DP": {
     "field": "DewPoint"
   }
 }
 ```

## Input/Output
```json
    "input": [
    ],
    "output": [
      {
        "name": "connectorMsg",
        "type": "object",
        "description": "The message for connectorMsg object"
      }
    ]
```
To build

go build github.com/matt-doug-davidson/aeroqualaqy1

This is it.
