# go-berry: a go library for all raspberry pi communication protocols

<p align="center">
    <img style="width: 20%" src="assets/goraspio_logo.png">
</p>

> [!WARNING]
> this library is still in development and is not yet ready for use

_go-berry_ is a go library that contains all common operations for using a raspberry pi with the go programming language. it includes the following communication protocols:

- gpio (general purpose input/output)
- pwm (pulse width modulation)
- spi (serial peripheral interface)
- i2c (inter-integrated circuit)
- udp (user data protocol)

## installation

to install the library, run the following command:

```bash
go get github.com/vistormu/go-berry
```

and don't forget to enable all the communication protocols that you want to use with `raspi-config`

## extra modules

_go-berry_ also includes some extra modules that can be used for different purposes. these modules are not directly related to the raspberry pi communication protocols, but they can be useful for different applications. my background is in robotics, so these modules are more focused on robotics applications.

- _utils_: useful operations
- _utils/num_: numeric operations
- _utils/signal_: useful algorithms for signal processing
- _peripherals_: functional code for sensors and actuators

## future work

the library is still in development, and there are many features that need to be added. 

some of the features that will be added in the future are:

- tcp client
- tcp and udp servers
- serial (serial communication)
- bluetooth 
- camera?

also, the library needs to be tested on different raspberry pi models 

lastly, the library needs to be documented properly

## other libraries

the code for this library is copied and modified from the following libraries:

- [go-rpio](https://github.com/stianeikeland/go-rpio) (MIT License)
- [go-i2c](https://github.com/d2r2/go-i2c) (MIT License)
