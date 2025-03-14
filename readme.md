# Setup Instructions

### 1. Resolve Dependencies
Before running the application, install all required dependencies. If any issues arise, resolve them before proceeding.

### 2. Database Setup
- Create an empty database.
- Configure your environment variables based on your local setup.

### 3. Run Migrations
Apply database migrations by running:  
```sh
make migrate.migration
```

### 4. API Testing
For easier testing, use the provided **Postman collection** located in the `postman/` folder.