# WolGate

This add-on runs on port **3881** ! 

## Installation

Follow these steps to get the add-on on your system:

1. Click the Home Assistant My button below to open the add-on on your Home
   Assistant instance.

   [![Add repository on my Home Assistant][repository-badge]][repository-url]

2. Find the "WolGate" add-on and click it.
3. Click on the "INSTALL" button
4. Enjoy the add-on!

## Configuration

You have only the configuration to set:
```yaml
# example
- url: https://sub1.domain.com
  macAddress: xx:xx:xx:xx:xx:xx
  broadcastAddress: 255.255.255.255:9
  ip: 192.168.0.200
  port: 1234
- url: https://sub2.domain.com
  macAddress: xx:xx:xx:xx:xx:xx
  broadcastAddress: 255.255.255.255:9
  ip: 192.168.0.200
  port: 5678
```

Now click on save and you are good to go :)