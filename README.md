# Github Gist Clone Backend

[![Go Build And Test](https://github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/actions/workflows/go.yml/badge.svg)](https://github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/actions/workflows/go.yml)

Backend REST API for GitHub Gist Clone built using Golang, Gin, GORM, Docker, and PostgreSQL and [other awesome Go libraries](https://github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/blob/master/go.mod).

## Running

### Pre-requisites

Apart from regular tools like Docker, Docker Compose, Go, etc. you will need following files to run the project:

1. `./app.env` - Environment variables for the application (Example)

```properties
POSTGRES_HOST=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password123
POSTGRES_DB=gist-backend
POSTGRES_PORT=5432

PORT=8000
CLIENT_ORIGIN=http://localhost:3000

ACCESS_TOKEN_PRIVATE_KEY=LS0tLS1CRUdJTiBSU0EgUFJJVkFHRSBLRVktLS0tLQpNSUlCUEFJQkFBSkJBTzVIKytVM0xrWC91SlRvRHhWN01CUURXSTdGU0l0VXNjbGFFKzlaUUg5Q2VpOGIxcUVmCnJxR0hSVDVWUis4c3UxVWtCUVpZTER3MnN3RTVWbjg5c0ZVQ0F3RUFBUUpCQUw4ZjRBMUlDSWEvQ2ZmdWR3TGMKNzRCdCtwOXg0TEZaZXMwdHdtV3Vha3hub3NaV0w4eVpSTUJpRmI4a25VL0hwb3piTnNxMmN1ZU9wKzVWdGRXNApiTlVDSVFENm9JdWxqcHdrZTFGY1VPaldnaXRQSjNnbFBma3NHVFBhdFYwYnJJVVI5d0loQVBOanJ1enB4ckhsCkUxRmJxeGtUNFZ5bWhCOU1HazU0Wk1jWnVjSmZOcjBUQWlFQWhML3UxOVZPdlVBWVd6Wjc3Y3JxMTdWSFBTcXoKUlhsZjd2TnJpdEg1ZGdjQ0lRRHR5QmFPdUxuNDlIOFIvZ2ZEZ1V1cjg3YWl5UHZ1YStxeEpXMzQrb0tFNXdJZwpQbG1KYXZsbW9jUG4rTkVRdGhLcTZuZFVYRGpXTTlTbktQQTVlUDZSUEs0PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==
ACCESS_TOKEN_PUBLIC_KEY=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d2RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTzVIKytVM0xrWC91SlRvRHhWN01CUURXSTdGU0l0VQpzY2xhRSs5WlFIOUNlaThiMXFFZnJxR0hSVDVWUis4c3UxVWtCUVpZTER3MnN3RTVWbjg5c0ZVQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==
ACCESS_TOKEN_EXPIRED_IN=15m
ACCESS_TOKEN_MAXAGE=15

REFRESH_TOKEN_PRIVATE_KEY=LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlCT1FJQkFBSkJBSWFJcXZXeldQSndnYjR1SEhFQ01RdHFZMTI5b2F5RzVZMGlGcG51a0J1VHpRZVlQWkE4Cmx4OC9lTUh3Rys1MlJGR3VxMmE2N084d2s3TDR5dnY5dVY4Q0F3RUFBUUpBRUZ6aEJqOUk3LzAxR285N01CZUgKSlk5TUJLUEMzVHdQQVdwcSswL3p3UmE2ZkZtbXQ5NXNrN21qT3czRzNEZ3M5T2RTeWdsbTlVdndNWXh6SXFERAplUUloQVA5UStrMTBQbGxNd2ZJbDZtdjdTMFRYOGJDUlRaZVI1ZFZZb3FTeW40YmpBaUVBaHVUa2JtZ1NobFlZCnRyclNWZjN0QWZJcWNVUjZ3aDdMOXR5MVlvalZVRlVDSUhzOENlVHkwOWxrbkVTV0dvV09ZUEZVemhyc3Q2Z08KU3dKa2F2VFdKdndEQWlBdWhnVU8yeEFBaXZNdEdwUHVtb3hDam8zNjBMNXg4d012bWdGcEFYNW9uUUlnQzEvSwpNWG1heWtsaFRDeWtXRnpHMHBMWVdkNGRGdTI5M1M2ZUxJUlNIS009Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0t
REFRESH_TOKEN_PUBLIC_KEY=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUbLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBSWFJcXZXeldCSndnYjR1SEhFQ01RdHFZMTI5b2F5Rwo1WTBpRnBudWtCdVR6UWVZUFpBOGx4OC9lTUh3Rys1MlJGR3VxMmE2N084d2s3TDR5dnY5dVY4Q0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==
REFRESH_TOKEN_EXPIRED_IN=60m
REFRESH_TOKEN_MAXAGE=60

EMAIL_FROM=admin@admin.com
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_USER=cb939a070d2d82
SMTP_PASS=719e9aba66e0f5
SMTP_PORT=2525
```
`ACCESS_TOKEN_MAXAGE` - Time in minutes
`REFRESH_TOKEN_MAXAGE` - Time in minutes

2. `./pgadmin.env` - PostgreSQL admin credentials (Example)

```properties
PGADMIN_DEFAULT_EMAIL=user@domain.com
PGADMIN_DEFAULT_PASSWORD=SuperSecret
```

**Note:** Modifying any of the above files might require you to change port numbers, etc. in the [`docker-compose.yml`](https://github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/blob/master/docker-compose.yaml) file.

### Starting the application

```bash
$ docker-compose up -d
```

The REST API is documented using swagger and can be accessed at http://localhost:8000/swagger/index.html (`8000` Port Number)

The docs are also checked into the repository and can be accessed at `./docs/`

### Stopping the application

```bash
$ docker-compose down
```

## Modifying the application

After modifying the Go source code, you can rebuild the application using the following command:

```bash
$ docker build -t yvyom/github-gist-backend:v1.0-alpha .
```

You can either use a different tag or remove the older image and use `pull_policy: if_not_present` in the [`docker-compose.yml`](https://github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/blob/master/docker-compose.yaml) file to pull the latest image.

## Directory Structure

```bash
.
├── app.env             # Application environment variables
├── controllers         # Application controllers
├── docker-compose.yaml # Docker compose file
├── Dockerfile          # Dockerfile for building the application
├── docs                # Swagger docs
├── go.mod              # Go module file
├── go.sum              # Go sum file
├── initializers        # Application initializers (Connecting to DB, Loading Env Files, etc.)
├── main.go             # Main file
├── middleware          # Application middleware (Checking if user is authenticated (DeserializeUser), etc.)
├── models              # Application models
├── pgadmin.env         # PostgreSQL admin credentials
├── README.md           # README file
├── routes              # Application routes
├── scripts             # Scripts for running the application (enabling uuid-ossp on database for uuid generation)
├── templates           # Application templates (for sending emails)
└── utils               # Application utilities (for sending emails, generating tokens, etc.)
```

Following actions are currently implemented:

- Check if username is available
- Check the basic health of api
- Create a comment on a gist
- Create a gist
- Follow a user
- Get the comments of a gist
- Get the current logged in user details.
- Get the gist by gist id, does not load gist comments
- Get the publicly visible details of a user, does not load gists
- Get the publicly visible gist Ids of a user
- Get the publicly visible gists of a user, does not load the gist comments
- Log out a user
- Refresh access token with refresh token
- Register a new user
- Resend verification email
- Reset password
- Send reset code for password reset
- Sign in a user
- Star a gist
- Unfollow a user
- Un-star a gist
- Update gist data
- Update user metadata
- Verify users email address

More endpoints would be implemented soon!

## Contributing to this project

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

Be Respectful and Kind to each other. We are all here to learn and grow together.

# Give a ⭐️ if you like this project!