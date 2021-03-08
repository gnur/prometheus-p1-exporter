# prometheus p1 exporter 

The prometheus p1 exporter is a simple (by design and purpose) binary that can read data from a smart meter through a serial port and expose these metrics to be scraped by prometheus.

## usage/configuration

```
usage: prometheus-p1-exporter --serial.port=SERIAL.PORT [<flags>]

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9602"  
                                 Address on which to expose metrics and web interface.
      --web.telemetry-path="/metrics"  
                                 Path under which to expose metrics.
      --serial.port=SERIAL.PORT  Serial port for the connection to the P1 interface.
```


## limitations

Lines are processed as they come in, no checksums are handled.

For now, only a few fields are exported: 

- power used in tariff 1 in Wh
- power used in tariff 2 in Wh
- current power draw in W

These values come from the following oids:

- 1-0:1.8.1 (meter reading, tariff 1 in kWh)
- 1-0:2.8.2 (meter reading, tariff 2 in kWh)
- 1-0:1.7.0 (Actual power delivered in kW)

## sources and acknowledgements

This repo for providing me with a direction: https://github.com/marceldegraaf/smartmeter  
This document for providing me with a overview of how p1 works: http://files.domoticaforum.eu/uploads/Smartmetering/DSMR%20v4.0%20final%20P1.pdf  
This page (Dutch) for a quick overview of how to use cu to get an example reading: https://infi.nl/nieuws/hobbyproject-slimme-meterkast-met-raspberry-pi/
