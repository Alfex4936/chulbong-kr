# chulbong-:kr:

![dotgraph](https://github.com/Alfex4936/chulbong-kr/assets/2356749/7b1e06ec-4478-4514-aa8f-9831481fd4d8)


```mermaid
erDiagram
    Users ||--o{ Markers : "can create"
    Users ||--o{ MarkerDislikes : "can dislike"
    Users ||--o{ Comments : "can comment"
    Users ||--o{ OpaqueTokens : "has"
    Users ||--o{ PasswordResetTokens : "has"
    Users ||--o{ Favorites : "can favorite"
    Markers ||--|{ Photos : "can have"
    Markers ||--|{ Comments : "can have"
    Markers ||--|{ MarkerDislikes : "can be disliked"
    Markers ||--|{ MarkerFacilities : "can have"
    Markers ||--|{ MarkerAddressFailures : "can fail"
    Markers ||--|{ Reports : "can have"
    Reports ||--|{ ReportPhotos : "can have"
    Notifications ||--|{ UserNotifications : "can have"

```

## Technologies and Frameworks

Below is the list of technologies and frameworks used in the Chulbong-KR backend:

- **Programming Language**: Go v1.22.3
- **Web Framework**: Fiber v2
- **Database**: MySQL 8
- **Caching**: Dragonfly (Redis-compatible)
- **Search Engine**: ZincSearch (Elasticsearch-compatible)
- **Chatting**: gorilla Websocket

## Dependency Management

- **Dependency Injection**: The project utilizes [Uber Fx](https://github.com/uber-go/fx) for managing dependencies effectively, following best practices in modularity and maintainability.

## Coding Conventions

- **Standards Followed**:
  - Google's Go programming language conventions
  - Uber's Go style guide

## Code Quality

- **Linter Configuration**:
  - See `revive.toml` for the linter settings used to ensure code quality and consistency across the project.
