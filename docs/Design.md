## High Level Design

```mermaid
graph LR
  A[User] --> B[CLI]
  A --> C[UI]
  B --> D{Handler}
  C --> D
  D -->|local calls| X[CSV Saved<br> Locally]
  D -->|reomote calls| Z[Deployed RDB]
```

## Handler Components
- Gonna use SQLite to import the data in form a csv.
- Storing the SQL data on machine at the location where the binary is.
- This can be interacting using CLI & UI.

# TODO
- [X] Add export and import functionality.
- [] Add querying functionality.
- [] Add filtering based on tags.
    - [] Explore DB approach creating separate table to store references.
    - [] Explore redis to create cached references on initialisation.

