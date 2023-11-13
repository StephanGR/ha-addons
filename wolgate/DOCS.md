# WolGate Documenation

## Installation

Follow these steps to get the add-on on your system:

1. Navigate in yout Home Assistant frontend to **Supervisor -> Add-on Store**
2. Add this new repository by URL (`https://github.com/StephanGR/ha-addons`)
3. Find the "WolGate" add-on and click it.
4. Click on the "INSTALL" button

## Configuration

You have 4 fields to fill to start up the application:

- wol_macAddress: MAC address of the thing you want to wake up
- vol_broadcastAddress: Broadcast Address
- domains: a list of url
```yaml
# example
  - url: https://sub1.domain.com
    address: 192.168.0.200
    port: 1234
  - url: https://sub2.domain.com
    address: 192.168.0.200
    port: 5678
```
- Network port: Port on which application will be running

Now click on save and you are good to go :)