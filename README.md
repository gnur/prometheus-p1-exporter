# prometheus p1 exporter 

The prometheus p1 exporter is a simple (by design and purpose) binary that can read data from a smart meter through a serial port and expose these metrics to be scraped by prometheus.

## configuration

The only things that can be configured with an environment var is:
 - **SERIAL_DEVICE**: the device that needs to be read (usually something like /dev/ttyUSB0)


## limitations

Lines are processed as they come in, no checksums are handled.

For now, only a few fields are exported: 

- power used in tariff 1 in Wh
- power used in tariff 2 in Wh
- current power draw in W
- gas used in cm2, updated hourly

These values come from the following oids:

- 1-0:1.8.1 (meter reading, tariff 1 in kWh)
- 1-0:2.8.2 (meter reading, tariff 2 in kWh)
- 1-0:1.7.0 (Actual power delivered in kW)
- 0-1:24.2.1 (gas meter reading, includes timestamp of reading)

## sources and acknowledgements

This repo for providing me with a direction: https://github.com/marceldegraaf/smartmeter  
This document for providing me with a overview of how p1 works: http://files.domoticaforum.eu/uploads/Smartmetering/DSMR%20v4.0%20final%20P1.pdf  
This page (Dutch) for a quick overview of how to use cu to get an example reading: https://infi.nl/nieuws/hobbyproject-slimme-meterkast-met-raspberry-pi/
