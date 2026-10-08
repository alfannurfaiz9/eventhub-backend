# EVENTHUB API
A project to event management

[![License](https://img.shields.io/badge/License-MIT-green)](https://opensource.org/license/mit) 
[![Go](https://img.shields.io/badge/Go-v1.27.1-blue?logo=go)](https://go.dev/) 
[![Gin gonic](https://img.shields.io/badge/Gin%20gonic-v1.12.0-blue?logo=gin&logoColor=white)](https://gin-gonic.com/) 
[![PostgreSQL](https://img.shields.io/badge/Postgre%20SQL-v16.15-blue?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-v8.4.7-red?logo=redis)](https://redis.io/)
[![Swagger](https://img.shields.io/badge/Swagger-v1.6.1-green?logo=swagger)](https://swagger.io/)

## Technologies
- **Go**
- **Gin gonic** For backend framework
- **PostgreSQL** For database
- **Redis** For caching
- **Swagger** For API documentation

## Features
- Authentication (Register, Login, Logout, Forgot password)
- User profile (Profile, Edit profile, Notification, Change password)
- Event list
- Community list
- Organizer (Create event, Dashboard & Analytics)
- Admin (Manage event, Manage user, Dashboard & Analytics)

## How to use
### Database Setup
1. Create your environment in the root directory named ``.env``
```bash
DB_HOST={YOUR_DB_HOST}
DB_PORT={YOUR_DB_PORT}
DB_PASS={YOUR_DB_PASS}
DB_NAME={YOUR_DB_NAME}
DB_USER={YOUR_DB_USER}
```

### Redis Setup
```bash
$ git clone developer
```

### Run
1. Clone this repository
```bash
$ git clone developer
```
2. Install the dependency
```bash
$ go mod download
```

## Routes
| Endpoint                | Method | Description        |
| -----------             | -------| -------------------|
| /auth/login             | POST   | Login              |
| /auth/register          | POST   | Register           |
| /auth/logout            | POST   | Logout             |
| /auth/forgot-password   | PATCH  | Forgot Password    |

### Documentation
For complete documentation, visit to ``/swagger/index.html``

## How to contribute
- Fork this repository
- Create your changes
- Pull request

## License
This project is licensed under the MIT License 
<br>
``LICENSE``

## Contacts
[Email](mailto:alfannurfaizwork@gmail.com)    

## Related Project
[Frontend](https://github.com/kodacampmain/koda-b7-react)