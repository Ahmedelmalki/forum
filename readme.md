# Forum Project

## Objectives

This project consists of creating a web forum that allows:

- Communication between users.
- Associating categories to posts.
- Liking and disliking posts and comments.
- Filtering posts.

---

## SQLite

To store the data in our forum (like users, posts, comments, etc.), we will use the database library **SQLite**.

SQLite is a popular choice as an embedded database software for local/client storage in application software such as web browsers. It enables us to create a database as well as control it using queries.

### Requirements:
- We must use at least one **SELECT**, one **CREATE**, and one **INSERT** query.

### Recommendation:
- Build an **entity-relationship diagram** (ERD) based on our database structure to achieve better performance.

For more information about SQLite, check the [SQLite documentation](https://www.sqlite.org/).

---

## Authentication

### Registration:
Users must be able to register on the forum by providing:
- **Email**: If the email is already taken, we must return an error response.
- **Username**
- **Password**: (Bonus) Passwords must be encrypted when stored.

### Login:
- The forum must check if the provided email exists in the database and validate the credentials.
- If the password is incorrect, we must return an error response.

### Sessions:
- We will use cookies to manage user sessions.
- Each session must include an **expiration date**.
- (Bonus) Use **UUIDs** to manage sessions.

---

## Communication

### Posts and Comments:
- Only **registered users** can create posts and comments.
- Posts can be associated with one or more **categories** (our choice).
- Posts and comments must be **visible to all users** (registered or not).

---

## Likes and Dislikes

- Only **registered users** can like or dislike posts and comments.
- The number of likes and dislikes must be **visible to all users** (registered or not).

---

## Filter Mechanism

We will implement a filter mechanism to allow users to filter displayed posts by:
- **Categories** (acts as subforums).
- **Created posts** (logged-in user's posts).
- **Liked posts** (logged-in user's liked posts).

*Note*: The last two filters are available only to registered users.

---

## Docker

We will use Docker to containerize our forum project. For more details, refer to the [Docker basics guide](https://docs.docker.com/get-started/).

---

## Instructions

- Use **SQLite**.
- Handle website errors and HTTP statuses.
- Handle all types of technical errors.
- Follow coding **best practices**.
- Include test files for **unit testing**.

---

## Allowed Packages

- All standard Go packages.
- `sqlite3`
- `bcrypt`
- `UUID`

*Note*: We must not use any frontend libraries or frameworks like React, Angular, or Vue.

---

## What We Will Learn

- Basics of web development:
  - **HTML**
  - **HTTP**
  - **Sessions and cookies**
- Using and setting up **Docker**:
  - Containerizing an application.
  - Creating images.
  - Compatibility/dependency management.
- SQL language:
  - Database manipulation.
- Basics of **encryption**.
