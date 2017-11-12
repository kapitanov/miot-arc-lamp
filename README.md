miot-arc-lamp
=============

Controller app for Arc Reactor Indicator.
Handles requests via MQTT.

Input messages
--------------

Send anything to `/arc/request` to get current indicator state. Response message will be pushed to `/arc/state`.

Output messages
---------------

Subscribe to `/arc/state` to receive indicator state updates. Message format:

```json
{
    "ring" : "blink",
    "core_r" : "on",
    "core_g" : "off",
    "core_b" : "off"
}
```

Possible states:

* `off`
* `on`
* `blink`

Configuration
-------------

App is configured via environment variables:

* `IMAP_ADDR` - IMAP server host and port (e.g. `imap.gmail.com:993`)
* `IMAP_USERNAME` - IMAP server login
* `IMAP_PASSWORD` - IMAP server password
* `IMAP_LABEL` - IMAP server label to check
* `MQTT_HOSTNAME` - MQTT broker hostname
* `MQTT_USERNAME` - MQTT broker login
* `MQTT_PASSWORD` - MQTT broker password

How to run
----------

1. Create `.env` file with all required environment variables (see above)
2. Run command:
   ```shell
   docker-compose up -d
   ```