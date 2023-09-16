# Saudi CERT Alert Parser & Public API

<div style="text-align:center;">
    <img src="./assets/CERT-logo-white.svg" alt="Saudi CERT Logo" width="400" />
</div>

Given the mission of Saudi CERT in enhancing cybersecurity awareness within the Kingdom of Saudi Arabia, this project aims to develop a parser to extract and structure alerts from Saudi CERT and expose them through a public API.

## Features
* Parse Saudi CERT security warnings.
* Retrieve information about security alerts, including severity, affected products, and recommendations.
* Store parsed data in a MongoDB database.
* Public API for accessing Saudi CERT security warnings.


## Running the Application Using Docker-Compose


1. **Docker and docker-compose**:
   Ensure you have Docker and docker-compose installed on your system.
   - [Get Docker](https://docs.docker.com/get-docker/)
   - [Install Docker Compose](https://docs.docker.com/compose/install/)

2. **Repository**:
   Clone the repository (if it's on a Git repository) or navigate to the directory where your `Dockerfile` and `docker-compose.yml` files reside.


### Deployment Steps:

```yaml
version: '3.9'

services:
  app:
    build:
      context: .
      dockerfile: ./docker/Dockerfile
    ports:
      - 8000:8000
    networks:
      - saudi_cert_network
    depends_on:
      - mongodb
  mongodb:
    image: mongo:latest
    ports:
      - 27017:27017
    environment:
      MONGO_INITDB_ROOT_USERNAME: saudi_cert_user
      MONGO_INITDB_ROOT_PASSWORD: saudi_cert_password
    volumes:
      - mongodb_data_volume:/data/db
    networks:
      - saudi_cert_network
  mongo-express:
    image: mongo-express:latest
    ports:
      - 8081:8081
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: saudi_cert_user
      ME_CONFIG_MONGODB_ADMINPASSWORD: saudi_cert_password
      ME_CONFIG_MONGODB_URL: mongodb://saudi_cert_user:saudi_cert_password@mongodb:27017/
    depends_on:
      - mongodb
    networks:
      - saudi_cert_network
volumes:
  mongodb_data_volume:
networks:
  saudi_cert_network:
```

Navigate to the directory containing your docker-compose.yml file in a terminal. Then run:

```bash
$ docker-compose up --build
```

## 1. Project Structure

#### 1.1 Parser
The parser is responsible for:
* Fetching the latest alerts from Saudi CERT's official website.
* Extracting relevant data such as threat name, severity, description, affected products, best practices, and links to official advisories.
* Structuring this data into a format suitable for database storage.

#### 1.2 Public API
The API will:
* Allow users to fetch the latest alerts.
* Permit detailed searches based on product names, severity, and other criteria.
* Provide endpoints for each type of threat for more specific searches.


This project is a parser + public API built using the following tech stack:

- Go
- MongoDB

## Endpoints

```go
app.Get("/api/v1/alerts", controllers.GetAlerts)
app.Get("/api/v1/alerts/:id", controllers.GetAlertByID)
```

## Example Response

```bash
{
    _id: ObjectId('6505ad858e7bdb36a6733c07'),
    severity_level: 'High',
    name: 'Lenovo Alert',
    image_url: 'https://cert.gov.sa/media/Lenovo_DnBzkUi.png',
    original_link: 'https://cert.gov.sa/en/security-warnings/lenovo-alert187654/',
    details: {
        best_practice: 'The CERT team encourages users to review Lenovo security advisory and update the affected products:https://support.lenovo.com/us/en/product_security/LEN-118374 https://support.lenovo.com/us/en/product_security/LEN-118320 ',
        description: 'Lenovo has released security updates to address multiple vulnerabilities in the following products:',
        targeted_sector: 'All',
        threats: '  ',
        warning_number: '2023-5482',
        warning_date: '1 March, 2023',
        affected_products: [
            {
                name: 'Converged HX'
            },
            {
                name: 'Desktop'
            },
            {
                name: 'Hyperscale'
            },
            {
                name: 'Storage'
            },
            {
                name: 'System x'
            },
            {
                name: 'ThinkAgile'
            },
            {
                name: 'ThinkServer'
            },
            {
                name: 'ThinkStation'
            },
            {
                name: 'ThinkSystem'
            }
        ],
        threat_list: [
            {
                name: 'Information disclosure'
            },
            {
                name: 'Escalation of Privilege'
            },
            {
                name: 'Denial of Service (DoS) Attack'
            }
        ],
        recommendations: [
            {
                link: 'https://support.lenovo.com/us/en/product_security/LEN-118374'
            },
            {
                link: 'https://support.lenovo.com/us/en/product_security/LEN-118320'
            }
        ]
    }
}

```