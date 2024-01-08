# User Service

Service that is responsible for user authentication and authorization.
It also provides user profile information.

## Postgres Database

### Users

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

| Column     | Type         | Description                                 |
| ---------- | ------------ | ------------------------------------------- |
| id         | UUID         | Unique identifier for the user              |
| email      | VARCHAR(255) | Unique email address for the user           |
| password   | VARCHAR(255) | Hashed password for the user                |
| created_at | TIMESTAMP    | Timestamp of when the user was created      |
| updated_at | TIMESTAMP    | Timestamp of when the user was last updated |

### Profiles

```sql
CREATE TABLE profiles (
    user_id UUID NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    age INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

| Column     | Type      | Description                                    |
| ---------- | --------- | ---------------------------------------------- |
| user_id    | UUID      | Unique identifier for the user                 |
| created_at | TIMESTAMP | Timestamp of when the profile was created      |
| updated_at | TIMESTAMP | Timestamp of when the profile was last updated |
| full_name  | VARCHAR   | Full name of the user                          |
| age        | INT       | Age of the user                                |
