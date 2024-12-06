# goraspio: a go library for raspberry pi gpio operations

<p align="center">
    <img style="width: 20%" src="assets/goraspio_logo.png">
</p>

> [!WARNING]
> this library is still in development and is not yet ready for use

_goraspio_ is a go library that contains all common operations for using a raspberry pi as a controller.

_goraspio_ offers the following modules:

- _gpio_: module for gpio communication. it includes digital output, digital input, pwm, spi, and i2c.
- _num_: numeric operations
- _algos_: useful algorithms
- _utils_: useful operations
- _actuator_: functional code for actuators
- _sensor_: functional code for sensors
- _client_: udp and tcp clients for sending information to a server
- _model_: code for using onnx machine learning models
