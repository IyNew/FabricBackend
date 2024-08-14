## Check if the server is alive
Use this command to check if the server is alive:
```
curl tjws-06.cse.nd.edu:6999/api/hello
```

If not, log in to the workstation, and use the following command to restart the server:
```
chmod +x background_deploy.sh
nohup ./background_deploy.sh &
```

## Function test

Query existing records
```bash
curl tjws-06.cse.nd.edu:6999/api/record/all
```

Create a new record with current time stamp
```json
{
    "drone": "grey",
    "story": "1",
    "zip": "96701",
    "datetime": "2024-01-02T00:00:00Z",
    "temperature": "24.8",
    "wind": "3.01",
    "gust": "9.3",
    "timesincelastmaintenance": "16.9192835242178",
    "flighthours": "29.3651318764408",
    "pitch": "1.46646019661782",
    "roll": "1.68349152692755",
    "yaw": "1.13800845548443",
    "vibex": "0.257336027891003",
    "vibey": "1.69496992420949",
    "vibez": "0.811888101171306",
    "nsat": "4",
    "noise": "35",
    "currentslope": "0.005",
    "brownout": "FALSE",
    "batterylevel": "0.8",
    "crash": "FALSE"
}
```

```bash
curl -X POST -H "Content-Type: application/json" -d '{"drone": "grey","story": "1","zip": "96701","datetime": "2024-01-02T00:00:00Z","temperature": "24.8","wind": "3.01","gust": "9.3","timesincelastmaintenance": "16.9192835242178","flighthours": "29.3651318764408","pitch": "1.46646019661782","roll": "1.68349152692755","yaw": "1.13800845548443","vibex": "0.257336027891003","vibey": "1.69496992420949","vibez": "0.811888101171306","nsat": "4","noise": "35","currentslope": "0.005","brownout": "FALSE","batterylevel": "0.8","crash": "FALSE"}' http://tjws-06.cse.nd.edu:6999/api/record/create
```

Query the newly added records
```bash
curl http://tjws-06.cse.nd.edu:6999/api/record/grey 
```

