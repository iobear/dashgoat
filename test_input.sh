#!/bin/bash

# An array of Scandinavian city names
cities=("Copenhagen" "Stockholm" "Oslo" "Helsinki" "Gothenburg" "Aarhus" "Bergen" "Uppsala" "Tampere" "Odense" "Stavanger" "Malmö" "Linköping" "Trondheim" "Västerås" "Lund" "Turku" "Oulu" "Reykjavik" "Jyväskylä")

# An array of possible status values
statuses=("ok" "warning" "error")

for i in {1..20}; do
    # Get the city name from the cities array
    city="${cities[$((i - 1))]}"
    echo $city
    
    # Vary the status in a cyclic fashion using modulo
    status_value="${statuses[$((i % 3))]}"
    
    curl --request POST \
      --url http://127.0.0.1:2000/update \
      --header 'content-type: application/json' \
      --data "{
    \"host\": \"$city\",
    \"service\": \"HTTP\",
    \"status\": \"$status_value\",
    \"message\": \"Hello World\",
    \"dependon\": \"tr-ch1\",
    \"updatekey\": \"my-precious!\"
    }"
    # Optional: Add sleep to avoid hitting the server too fast.
    sleep .1
done

