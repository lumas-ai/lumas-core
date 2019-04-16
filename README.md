Lumas enables person detection and HomeKit support to off the shelf IP camera.

![HomeKit notification](images/notification.jpg)

Currently it only support Amcrest IP cameras, but more camera support is coming.
Presently it also only supports one camera at a time.

## Setup

### Requirements

**Software Requirements:** 
* Docker (CE or EE)

**Hardware Requirements:**
* Architecture: x86_64

### Quick Start

#### Run the Lumas server
1) Create a directory called `lumas`
2) Clone this repo to the `lumas` directory
3) Clone the [ONVIF provider extension](https://github.com/lumas-ai/lumas-provider-onvif) into the `lumas` directory
3) Copy the docker-compose.yml file from this repo to the `lumas` directory
4) Run with `docker-compose up`

#### Create a configuration

1) On your workstation, download the [lumasctl](https://github.com/lumas-ai/lumas-core/releases/tag/v0.1.0-alpha.1) binary for your platform
2) Create a `config.yml` file. There's an example config.yml in this repository
3) Apply the configuration with `./lumasctl apply --controller <ADDRESS OF THE LUMAS SERVER>:5389 config.yml`
4) Ensure sure your cameras were applied with `./lumasctl --controller <ADDRESS OF THE LUMAS SERVER>:5389> camera list`

#### Client config file

In order to not have to specify the address of the controller each time, you
can create a client configuration file at `~/.lumasctl.cfg`. It is a yaml
format and accepts the `controller` parameter

For example:
```
---
controller: "192.168.2.207:5389"
```


The basic configuration file structure looks like this:
```
---
cameras:
  - name: "Camera Name"
    provider:
      name: <camera provider extension>
      config:
        <provider parameters>
```

### Cameras parameters

The `cameras` section of the configuration is a list of cameras. Currently
Lumas ony supports one camera at a time. Multi-camera support will be added
soon.

Also, currently only ONVIF cameras are support. Other camera support will be added soon.

Example:
```
---
cameras:
  - name: "<Camera Name>"
    plugin:
      name: onvif
      config:
        rtspAddress: "rtsp://<username>:<password>@<camera address>"
```

### Running

Run the application with `docker-compose up -d`

## Roadmap

* More camera vendor support: Amcrest, Foscam, etc.
* Event system with easily configured responses
* Face recognition - Learn familiar faces over time with custom alerts
* Time lapse - See what's been happening through the day
