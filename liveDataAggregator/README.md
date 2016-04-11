====================
LIVE DATA AGGREGATOR
====================
Dials single address that produces a stream of deltas, aggregates deltas and serves to multiple websocket connections

# Components
Connections - manages TCP / websocket connections
Parser - translates incoming TCP byte stream into readable format for aggregator
Aggregator - takes a stream of incoming deltas and stores an in-memory sum
* Parser and Aggregator are specifically typed and will need to change along with data formats

# Process
## Data
Live data comes via TCP socket
Parser translates into map[string]string
LightAggregator aggregates while broadcast is blocked
LightAggregator blocks catchup, sends to FullAggregator and broadcasts to existing connections
LightAggregator unblocks catchup

## Connections
Incoming connections come via TCP socket
Connections stored in new connections pool until catchup is unblocked
Connections blocks broadcast, catches up all new connections and places into existing connections pool
Connections unblocks broadcast
