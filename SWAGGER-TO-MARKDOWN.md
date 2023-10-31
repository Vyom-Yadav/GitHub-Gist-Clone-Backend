# GitHub Gist Backend REST API

The REST API for GitHub Gist Backend

## Base URL

Base URL: localhost:8000/api/ 
Change the port no. if you are running the server on a different port.

## Version: 1.0-alpha

### /auth/forgotpassword

#### POST

##### Summary:

Send reset code for password reset

##### Parameters

| Name                | Located in | Description                               | Required | Schema                                                    |
|---------------------|------------|-------------------------------------------|----------|-----------------------------------------------------------|
| ForgotPasswordInput | body       | The Input for sending password reset code | Yes      | [models.ForgotPasswordInput](#models.ForgotPasswordInput) |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 400  | Bad Request           | object |
| 401  | Unauthorized          | object |
| 500  | Internal Server Error | object |

### /auth/github/callback

#### GET

##### Summary:

Register a new user through GitHub

##### Parameters

| Name        | Located in | Description                                                     | Required | Schema |
|-------------|------------|-----------------------------------------------------------------|----------|--------|
| code        | query      | code received from GitHub API after user authorizes application | Yes      | string |
| newUsername | query      | new username to be used for the user                            | No       | string |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 201  | Created               | object |
| 400  | Bad Request           | object |
| 409  | Conflict              | object |
| 500  | Internal Server Error | object |
| 502  | Bad Gateway           | object |

### /auth/github/clientid

#### GET

##### Summary:

Get GitHub Client ID

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 500  | Internal Server Error | object |

### /auth/login

#### POST

##### Summary:

Sign in a user

##### Parameters

| Name        | Located in | Description    | Required | Schema                                    |
|-------------|------------|----------------|----------|-------------------------------------------|
| SignInInput | body       | Sign in a user | Yes      | [models.SignInInput](#models.SignInInput) |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 400  | Bad Request           | object |
| 409  | Conflict              | object |
| 500  | Internal Server Error | object |

### /auth/logout

#### GET

##### Summary:

Log out a user

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |

### /auth/refresh

#### GET

##### Summary:

Refresh access token with refresh token

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 403  | Forbidden             | object |
| 500  | Internal Server Error | object |

### /auth/register

#### POST

##### Summary:

Register a new user

##### Parameters

| Name        | Located in | Description                                      | Required | Schema                                    |
|-------------|------------|--------------------------------------------------|----------|-------------------------------------------|
| SignUpInput | body       | User object that needs to be added to the system | Yes      | [models.SignUpInput](#models.SignUpInput) |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 201  | Created               | object |
| 400  | Bad Request           | object |
| 409  | Conflict              | object |
| 500  | Internal Server Error | object |
| 502  | Bad Gateway           | object |

### /auth/resendverificationemail

#### POST

##### Summary:

Resend verification email

##### Parameters

| Name                         | Located in | Description                                                | Required | Schema                                                                      |
|------------------------------|------------|------------------------------------------------------------|----------|-----------------------------------------------------------------------------|
| ResendVerificationEmailInput | body       | Resend verification email to the user with the given email | Yes      | [models.ResendVerificationEmailInput](#models.ResendVerificationEmailInput) |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 400  | Bad Request           | object |
| 409  | Conflict              | object |
| 500  | Internal Server Error | object |

### /auth/resetpassword

#### PATCH

##### Summary:

Reset password

##### Parameters

| Name               | Located in | Description                              | Required | Schema                                                  |
|--------------------|------------|------------------------------------------|----------|---------------------------------------------------------|
| ResetPasswordInput | body       | The input required to reset the password | Yes      | [models.ResetPasswordInput](#models.ResetPasswordInput) |
| username           | query      | The username of the user                 | Yes      | string                                                  |
| resetToken         | query      | The token required to reset the password | Yes      | string                                                  |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 400  | Bad Request | object |

### /auth/usernameavailable/{username}

#### GET

##### Summary:

Check if username is available

##### Parameters

| Name     | Located in | Description                        | Required | Schema |
|----------|------------|------------------------------------|----------|--------|
| username | path       | The username to check availability | Yes      | string |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |

### /auth/verifyemail

#### GET

##### Summary:

Verify users email address

##### Parameters

| Name             | Located in | Description                                  | Required | Schema |
|------------------|------------|----------------------------------------------|----------|--------|
| username         | query      | Username of the user to be verified          | Yes      | string |
| verificationCode | query      | Verify the added user object to the database | Yes      | string |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 400  | Bad Request | object |
| 409  | Conflict    | object |

### /gists/{gistId}

#### GET

##### Summary:

Get the gist by gist id, does not load gist comments

##### Parameters

| Name   | Located in | Description        | Required | Schema |
|--------|------------|--------------------|----------|--------|
| gistId | path       | The ID of the gist | Yes      | string |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 404  | Not Found   | object |

### /gists/{gistId}/comments

#### GET

##### Summary:

Get the comments of a gist

##### Parameters

| Name   | Located in | Description        | Required | Schema |
|--------|------------|--------------------|----------|--------|
| gistId | path       | The ID of the gist | Yes      | string |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 404  | Not Found   | object |

### /gists/{gistId}/stargazers

#### GET

##### Summary:

Get the stargazers of a gist

##### Parameters

| Name   | Located in | Description        | Required | Schema |
|--------|------------|--------------------|----------|--------|
| gistId | path       | The ID of the gist | Yes      | string |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 400  | Bad Request | object |
| 404  | Not Found   | object |

### /health

#### GET

##### Summary:

Check the basic health of api

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |

### /users/{username}

#### GET

##### Summary:

Get the publicly visible details of a user, does not load gists

##### Parameters

| Name     | Located in | Description         | Required | Schema |
|----------|------------|---------------------|----------|--------|
| username | path       | The username to get | Yes      | string |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 404  | Not Found   | object |

### /users/{username}/followers

#### GET

##### Summary:

Get the followers of a user

##### Parameters

| Name     | Located in | Description                                      | Required | Schema |
|----------|------------|--------------------------------------------------|----------|--------|
| username | path       | The username of the user to get the followers of | Yes      | string |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 404  | Not Found             | object |
| 500  | Internal Server Error | object |

### /users/{username}/following

#### GET

##### Summary:

Get the list of users a user is following

##### Parameters

| Name     | Located in | Description                                      | Required | Schema |
|----------|------------|--------------------------------------------------|----------|--------|
| username | path       | The username of the user to get the following of | Yes      | string |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 404  | Not Found             | object |
| 500  | Internal Server Error | object |

### /users/{username}/follows/{otherUser}

#### GET

##### Summary:

Whether a username follows another username

##### Parameters

| Name      | Located in | Description                             | Required | Schema |
|-----------|------------|-----------------------------------------|----------|--------|
| username  | path       | The username of the follower            | Yes      | string |
| otherUser | path       | The username of the user being followed | Yes      | string |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 404  | Not Found   | object |

### /users/{username}/gistIds

#### GET

##### Summary:

Get the publicly visible gist Ids of a user

##### Parameters

| Name     | Located in | Description                   | Required | Schema |
|----------|------------|-------------------------------|----------|--------|
| username | path       | The username to get gists for | Yes      | string |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 404  | Not Found   | object |

### /users/{username}/gists

#### GET

##### Summary:

Get the publicly visible gists of a user, does not load the gist comments

##### Parameters

| Name     | Located in | Description                   | Required | Schema |
|----------|------------|-------------------------------|----------|--------|
| username | path       | The username to get gists for | Yes      | string |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 404  | Not Found   | object |

### /users/{username}/starredGist/{gistId}

#### GET

##### Summary:

Whether a user has starred a gist

##### Parameters

| Name     | Located in | Description                                           | Required | Schema |
|----------|------------|-------------------------------------------------------|----------|--------|
| username | path       | The username of the user to check the starred gist of | Yes      | string |
| gistId   | path       | The ID of the gist to check if it is starred          | Yes      | string |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 400  | Bad Request | object |
| 404  | Not Found   | object |

### /users/{username}/starredGists

#### GET

##### Summary:

Get the list of starred gists of a user

##### Parameters

| Name     | Located in | Description                                          | Required | Schema |
|----------|------------|------------------------------------------------------|----------|--------|
| username | path       | The username of the user to get the starred gists of | Yes      | string |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 404  | Not Found             | object |
| 500  | Internal Server Error | object |

### /users/comments

#### POST

##### Summary:

Create a comment on a gist

##### Parameters

| Name               | Located in | Description                    | Required | Schema                                                      |
|--------------------|------------|--------------------------------|----------|-------------------------------------------------------------|
| CreateCommentInput | body       | The Input for creating comment | Yes      | [models.CommentOnGistRequest](#models.CommentOnGistRequest) |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 201  | Created     | object |
| 400  | Bad Request | object |

### /users/details

#### PATCH

##### Summary:

Update user metadata

##### Parameters

| Name                   | Located in | Description                          | Required | Schema                                                              |
|------------------------|------------|--------------------------------------|----------|---------------------------------------------------------------------|
| UpdateUserDetailsInput | body       | The Input for updating user metadata | Yes      | [models.UpdateUserDetailsRequest](#models.UpdateUserDetailsRequest) |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |
| 400  | Bad Request | object |

### /users/follow/{userToFollow}

#### PATCH

##### Summary:

Follow a user

##### Parameters

| Name         | Located in | Description                        | Required | Schema |
|--------------|------------|------------------------------------|----------|--------|
| userToFollow | path       | The username of the user to follow | Yes      | string |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 400  | Bad Request           | object |
| 404  | Not Found             | object |
| 500  | Internal Server Error | object |

### /users/gists

#### PATCH

##### Summary:

Update gist data

##### Parameters

| Name            | Located in | Description                      | Required | Schema                                                |
|-----------------|------------|----------------------------------|----------|-------------------------------------------------------|
| UpdateGistInput | body       | The Input for updating user gist | Yes      | [models.UpdateGistRequest](#models.UpdateGistRequest) |

##### Responses

| Code | Description  | Schema |
|------|--------------|--------|
| 200  | OK           | object |
| 400  | Bad Request  | object |
| 401  | Unauthorized | object |
| 404  | Not Found    | object |

#### POST

##### Summary:

Create a gist

##### Parameters

| Name            | Located in | Description                 | Required | Schema                                                |
|-----------------|------------|-----------------------------|----------|-------------------------------------------------------|
| CreateGistInput | body       | The Input for creating gist | Yes      | [models.CreateGistRequest](#models.CreateGistRequest) |

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 201  | Created     | object |
| 400  | Bad Request | object |

### /users/gists/{gistId}/star

#### PATCH

##### Summary:

Star a gist

##### Parameters

| Name   | Located in | Description                | Required | Schema |
|--------|------------|----------------------------|----------|--------|
| gistId | path       | The ID of the gist to star | Yes      | string |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 404  | Not Found             | object |
| 500  | Internal Server Error | object |

### /users/gists/{gistId}/unstar

#### PATCH

##### Summary:

Un-star a gist

##### Parameters

| Name   | Located in | Description                   | Required | Schema |
|--------|------------|-------------------------------|----------|--------|
| gistId | path       | The ID of the gist to un-star | Yes      | string |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 404  | Not Found             | object |
| 500  | Internal Server Error | object |

### /users/me

#### GET

##### Summary:

Get the current logged in user details.

##### Responses

| Code | Description | Schema |
|------|-------------|--------|
| 200  | OK          | object |

### /users/unfollow/{userToUnfollow}

#### PATCH

##### Summary:

Unfollow a user

##### Parameters

| Name           | Located in | Description                          | Required | Schema |
|----------------|------------|--------------------------------------|----------|--------|
| userToUnfollow | path       | The username of the user to unfollow | Yes      | string |

##### Responses

| Code | Description           | Schema |
|------|-----------------------|--------|
| 200  | OK                    | object |
| 400  | Bad Request           | object |
| 404  | Not Found             | object |
| 500  | Internal Server Error | object |

### Models

#### models.CommentOnGistRequest

| Name    | Type   | Description | Required |
|---------|--------|-------------|----------|
| content | string |             | Yes      |
| gistId  | string |             | Yes      |

#### models.CreateGistRequest

| Name    | Type    | Description | Required |
|---------|---------|-------------|----------|
| content | string  |             | Yes      |
| name    | string  |             | Yes      |
| private | boolean |             | No       |
| title   | string  |             | Yes      |

#### models.ForgotPasswordInput

| Name  | Type   | Description | Required |
|-------|--------|-------------|----------|
| email | string |             | Yes      |

#### models.ResendVerificationEmailInput

| Name  | Type   | Description | Required |
|-------|--------|-------------|----------|
| email | string |             | Yes      |

#### models.ResetPasswordInput

| Name            | Type   | Description | Required |
|-----------------|--------|-------------|----------|
| password        | string |             | Yes      |
| passwordConfirm | string |             | Yes      |

#### models.SignInInput

| Name     | Type   | Description | Required |
|----------|--------|-------------|----------|
| email    | string |             | Yes      |
| password | string |             | Yes      |

#### models.SignUpInput

| Name            | Type   | Description | Required |
|-----------------|--------|-------------|----------|
| email           | string |             | Yes      |
| firstName       | string |             | Yes      |
| lastName        | string |             | Yes      |
| password        | string |             | Yes      |
| passwordConfirm | string |             | Yes      |
| username        | string |             | Yes      |

#### models.UpdateGistRequest

| Name    | Type    | Description | Required |
|---------|---------|-------------|----------|
| content | string  |             | No       |
| gistId  | string  |             | Yes      |
| name    | string  |             | No       |
| private | boolean |             | No       |
| title   | string  |             | No       |

#### models.UpdateUserDetailsRequest

| Name           | Type   | Description | Required |
|----------------|--------|-------------|----------|
| location       | string |             | No       |
| profilePicture | string |             | No       |
| statusIcon     | string |             | No       |
| tagline        | string |             | No       |
| twitter        | string |             | No       |
| website        | string |             | No       |