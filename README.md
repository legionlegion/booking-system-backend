# Book4U Backend

This repository contains the server-side code and API endpoints that power the Book4U facility booking system. The backend is hosted at [https://bookingsyal-cbd544b30b67.herokuapp.com/](https://bookingsyal-cbd544b30b67.herokuapp.com/).

## Features

Here are the key features of the Book4U backend:

1. **Root Endpoint (/)**
   - Fetches and displays approved bookings from the `approvedbookings` table within the current and the upcoming week.

2. **User Management**
   - **Register Endpoint (/register)**: Registers users with their username and password into the `users` table.
   - **Authenticate Endpoint (/authenticate)**: On successful user login, it issues JWT tokens. These include an access token for authentication and authorization, and a refresh token for obtaining a new access token.
   - **Refresh Endpoint (/refresh)**: Obtains a new JWT token using refresh tokens securely stored in HTTP-only cookies. 
   - **Logout Endpoint (/logout)**: Invalidates the user's refresh token, logging them out.

3. **Protected Routes (/admin)**
   - Only accessible to users with a valid JWT token; unauthenticated requests receive a 401 Unauthorized status.

4. **Booking Management Endpoints**
   - **/add-booking**: Inserts a new booking into the `requestedbookings` table.
   - **/all-bookings**: Retrieves all bookings from the `approvedbookings` table.
   - **/approve-booking**: Transfers a booking from `requestedbookings` to `approvedbookings`.
   - **/booking-management**: Displays bookings based on user or admin roles.
   - **/delete-pending**: Deletes a booking from `requestedbookings`.
   - **/delete-approved**: Deletes a booking from `approvedbookings`.
   - **/user-bookings**: Fetches the bookings associated with the logged-in user from the `approvedbookings` table.

## Token Management

- Access tokens have a validity of 15 minutes.
- Refresh tokens remain valid for 24 hours, facilitating the generation of new access tokens without repeated user logins.

## Documentation

For a comprehensive overview and details of the project, visit our [detailed project documentation](https://docs.google.com/document/d/1lBbl30woSB4tnFogro37Me5kaN_AsGDz815Tp-ZH0GI/edit?usp=sharing).

## License

Book4U Backend is distributed under the [MIT License](LICENSE).

---
