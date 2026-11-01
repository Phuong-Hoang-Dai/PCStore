#!/bin/sh
# Creates the "DDStore logs" data view (the default one) and saved search in Kibana, for the
# logs Filebeat ships to ddstore-logs-*. Each call overwrites what is there, so it can rerun.
set -eu

KIBANA=http://kibana:5601

# Kibana answers 503 until it is ready, which can take a minute or two.
until curl -fs "$KIBANA/api/status" >/dev/null; do sleep 5; done

kibana() {
    curl -fsS -X POST "$KIBANA$1" -H 'kbn-xsrf: true' -H 'Content-Type: application/json' -d "$2" >/dev/null
}

kibana /api/data_views/data_view '{
    "override": true,
    "data_view": {"id": "ddstore-logs", "title": "ddstore-logs-*", "name": "DDStore logs", "timeFieldName": "@timestamp"}
}'
kibana /api/data_views/default '{"data_view_id": "ddstore-logs", "force": true}'

kibana '/api/saved_objects/search/ddstore-logs?overwrite=true' '{
    "attributes": {
        "title": "DDStore logs",
        "columns": ["service", "level", "kind", "transaction_id", "message", "duration_ms"],
        "sort": [["@timestamp", "desc"]],
        "kibanaSavedObjectMeta": {
            "searchSourceJSON": "{\"query\":{\"query\":\"\",\"language\":\"kuery\"},\"filter\":[],\"indexRefName\":\"kibanaSavedObjectMeta.searchSourceJSON.index\"}"
        }
    },
    "references": [{"name": "kibanaSavedObjectMeta.searchSourceJSON.index", "type": "index-pattern", "id": "ddstore-logs"}]
}'

echo "Kibana is set up: Discover > DDStore logs"
