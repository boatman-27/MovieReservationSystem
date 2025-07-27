# Movie Reservation System

This project is a complete backend system for a movie booking application, built with Go (Gin) and PostgreSQL. It allows users to sign up, log in, browse movies, view showtimes, and book or cancel reservations. Admin users can manage movies, showtimes, and view user activity.

> This implementation is inspired by the [roadmap.sh Movie Reservation System Project](https://roadmap.sh/projects/movie-reservation-system). For a deeper dive into the design and architecture, please refer to the [Project Wiki](https://github.com/labasubagia/movie-reservation-system/wiki).

---

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **ORM**: [sqlx](https://github.com/jmoiron/sqlx)
- **Authentication**: JWT (Access and Refresh Tokens)

---

## Authentication

- Users sign up and receive JWT access and refresh tokens.
- Protected routes require a valid access token.
- Admin routes are further secured with role-based authorization.

---

## Features

### User

- Signup / Login
- JWT Token generation
- Role-based promotion to admin

### Movies

- Get all movies
- Get movie by ID
- Admin can add, update, delete movies

### Showtimes

- Add, update, delete showtimes (admin)
- Get showtimes by movie ID
- Check seat availability

### Reservations

- Book seats
- Cancel reservations
- View upcoming reservations
- Admin can view all user reservations

---

## Getting Started

¯### 1. Clone the repository

```bash
git clone https://github.com/boatman-27/MovieReservationSystem
cd MovieReservationSystem
```

### 2. Create `.env` file

```env
POSTGRES_DSN="user=postgres password=password dbname=name sslmode=disable"
```

### 3. Run the server

```bash
go run main.go
```

### 4. Import the SQL Tables

- #### Users table

```sQL
CREATE TABLE users (
	userid VARCHAR(36) PRIMARY KEY,
	name TEXT NOT NULL,
	email VARCHAR(50) NOT NULL UNIQUE,
	password TEXT NOT NULL,
	role VARCHAR(5) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
	createdat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

- #### Movies table

```sQL
CREATE TABLE movies (
	movieid VARCHAR(10) PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT,
	genre TEXT NOT NULL,
	duration INT NOT NULL, -- in minutes
	director TEXT NOT NULL,
	posterimage Text NOT NULL,
	releasedate TIMESTAMP NOT NULL,
	createdat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

- #### Showtimes table

```sQL
CREATE TABLE showtimes (
	showtimeid VARCHAR(10) PRIMARY KEY,
	movieid VARCHAR(10) NOT NULL REFERENCES movies(movieid) ON DELETE CASCADE,
	starttime TIMESTAMP NOT NULL,
	endtime TIMESTAMP NOT NULL,
	venue TEXT NOT NULL,
	priceperseat NUMERIC(6,2) NOT NULL,
	availableseats INT NOT NULL CHECK (availableseats >= 0),
	createdat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

- #### Reservations table

```sQL
CREATE TABLE reservations (
    reservationid    VARCHAR(10) PRIMARY KEY,
    userid           TEXT NOT NULL REFERENCES users(userid),
    showtimeid       TEXT NOT NULL REFERENCES showtimes(showtimeid),
    numberofseats    INT NOT NULL CHECK (numberofseats > 0),
    totalprice       NUMERIC(10, 2) NOT NULL CHECK (totalprice >= 0),
    reservationdate  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

## API Endpoints

### /account

- `POST /signup` - Register new user
- `POST /login` - Authenticate user

### /protected _(Requires JWT)_

- `GET /movies` - Get all movies
- `POST /get-movie-byid`
- `POST /get-showtime-and-movie`
- `POST /get-seatsinfo`
- `POST /book-seats`
- `POST /upcoming-reservations`
- `POST /cancel-reservation`

### /admin _(Requires Admin Role)_

- `POST /promote`
- `POST /add-movie`
- `PATCH /update-movie`
- `POST /delete-movie`
- `POST /add-showtime`
- `PATCH /update-showtime`
- `POST /delete-showtime`
- `GET /all-reservations`
- `POST /user-reservations`
