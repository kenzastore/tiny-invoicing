# Technology Stack: Tiny Invoicing

## Overview
This document outlines the core technologies utilized in the Tiny Invoicing application.

## Programming Language
- **Go**: The primary programming language for the application's backend logic and API services.

## Database
- **MySQL**: A relational database management system used for persistent data storage, including invoices, client information, and other application data.

## Testing Libraries
- **go-sqlmock**: Used for mocking SQL database interactions in unit tests.
- **Interface Dependency Injection**: Used to decouple handlers from the database for cleaner unit testing with custom mocks.
