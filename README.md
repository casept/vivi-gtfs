# vivi-gtfs

Tool for receiving realtime data from the Websocket powering https://trainmap.vivi.lv/, and converting it to a GTFS-RT feed.

The feed has been tested to work fine in MOTIS.

## Usage

Download the non-realtime GTFS feed from https://www.pv.lv/xml/atdtransit/GTFS.zip. This is needed to match certain identifiers.

Then, run the program like so:

```sh
./vivi-gtfs -gtfs-path GTFS.zip
```

By default, the GTFS-RT feed is available at `http://localhost:1337/gtfs-rt.pb`.
