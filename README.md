# 1Devices

## Overview

This API is implemented in  **Golang**  and uses  **PostgreSQL**  as its persistence layer. Both the API and the database are containerized using  **Docker**, and orchestrated via  **Docker Compose**. The application runs on  `127.0.0.1`  at port  `8080`.

The API is designed to manage telecommunication devices, adhering to specific business rules and domain validations. Below are the key features and constraints of the  **Device Domain**:

### Device Domain

-   **Id**: Unique identifier for the device.
    
-   **Name**: Name of the device.
    
-   **Brand**: Brand of the device.
    
-   **State**: Current state of the device, which can be one of the following:
    
    -   `available`
        
    -   `in-use`
        
    -   `inactive`
        
-   **Creation Time**: Timestamp indicating when the device was created.
    

### Supported Functionalities

-   **Create a new device**: Add a new device to the system.
    
-   **Update an existing device**: Fully or partially update the details of an existing device.
    
-   **Fetch a single device**: Retrieve details of a specific device by its ID.
    
-   **Fetch all devices**: Retrieve a list of all devices in the system.
    
-   **Fetch devices by brand**: Retrieve devices filtered by their brand.
    
-   **Fetch devices by state**: Retrieve devices filtered by their state.
    
-   **Delete a single device**: Remove a device from the system.
    

### Domain Validations

-   **Creation Time**: The creation time of a device cannot be updated.
    
-   **Name and Brand**: The  `name`  and  `brand`  properties cannot be updated if the device is in the  `in-use`  state.
    
-   **In-use Devices**: Devices that are in the  `in-use`  state cannot be deleted.


## Usage

### 1. Clone the Repository

    git clone https://github.com/paulocesarvasco/1Devices_API
    cd 1Devices_API

### 2. Start the Application

    docker-compose up

### 3. Access the service

    http://127.0.0.1:8080/api/v1/

An HTML file has been added to enhance the usability of the services. By accessing the service through a browser, the API will route and serve this file. This interface allows users to view the list of devices, as well as edit or delete them directly from the browser.

## Documentation

The API documentation was implemented using the OpenAPI pattern and can be found in the `doc` folder. And it is possible to render the file in:

    https://editor.swagger.io/

## Development workflow

For project development management, a GitHub board was created where tasks were tracked from the documentation phase through TDD and feature development.

A **branch model** was established, where the `main` branch serves as the productive branch, and all new features must follow the naming convention `feature/*` before being merged.

A **CI pipeline** was implemented to enforce the branch model, build the project, and run tests. The pipeline also evaluates test coverage, which must meet a minimum threshold of **70%**, and ensures that the source code follows Golang's standard formatting.

The project can be accessed here:

    https://github.com/users/paulocesarvasco/projects/1/views/1
