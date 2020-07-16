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
| username | string | True | Username for log-in to the instrument |
| password | string | True | Password for log-in to the instrument |
| instrument | string | True | Instrument identification |
| entity | string | True | Entity to which measurements are related. Exported identifier to other subsystems |
| mappings | string | True | JSON format that contains mapping between instrument key and output field |

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
