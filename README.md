# google-photo-downloader
Go Application to download Google Photo Album Content

## Concepts Include

- OAuth2 Authentication with Google Photo API
- Creating HTTP Server and get user consent
- Accessing All Albums of the User
- Download Contents(photo, videos etc.) of the album
- Usage of Go Routine to enable parallel downloading of multiple media items at the same time.

## To get started:

- Create a New Project in Google Cloud Platform
- Enable Google Photo API for your new project
- Create Credential using OAuth 2.0 Client ID
  - Select Application Type as "Web Application"
  - Provide Authorised JavaScript origins as http://localhost:8080
  - Provide Authorised redirect URIs as http://localhost:8080/callback
  - Hit save.
- Create .env file at project root and populate it using same key as in env-example.txt file.
- You are good to go!

