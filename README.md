# Users Service

## Supported Routes

### User Login
[GIN-debug] POST   /api/v1/func/users/login

### User Sign (creating new account)
```
[GIN-debug] POST   /api/v1/func/users/signup
```
### One Time Password(OTP) login is generic route. It support :
 1. Confirm Email using OTP
 2. Login using Two Factor Authentication
 3. Login using Activating and Using Authenticator(like Google Authenticator)
```
[GIN-debug] POST   /api/v1/func/users/otp_login
```
### Extending JWT Token Life
```
[GIN-debug] POST   /api/v1/func/users/re_auth
```
### Enabling Two Factor Authentication
```
[GIN-debug] GET    /api/v1/func/users/enable_two_factor
```
### Enable Authenticator (QR code / Secret Code)
```
[GIN-debug] GET    /api/v1/func/users/enable_auth
```
### Email Confirm Route
```
[GIN-debug] GET    /api/v1/func/users/confirm_email
```
### Fetching User Detail Using User Code
```
[GIN-debug] GET    /api/v1/users/code/:userCode
```
### Getting Login User Detail
```
[GIN-debug] GET    /api/v1/users/profile
```

