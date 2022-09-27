# Overview

A template supports the most basic things for user - user, and user - group communicate in **game client**. The differences of this template from normal social platform chat mechanism are that room state will be reset if all the users connected to it have left and group messages are not archived anywhere.

# Features

- [x] Firebase Authentication
- [ ] User
  - [x] Send private message
  - [x] Notify friends status
  - [ ] Friend relationship
    - [x] Get friend list
    - [ ] Add friend
    - [ ] Remove friend
- [x] Room
  - [x] Create room
  - [x] Join room (Multi-room join mode)
  - [x] Leave room
  - [x] Kick out of room
  - [x] Transfer ownership
  - [x] Invite to room
  - [x] Respond to room invitation
  - [x] Notify room chages
  - [x] Send group message (Don't store into database)

# Tools And Technologies

- [NestJS](https://nestjs.com/)
- [Socket.IO](https://socket.io/)
- [Prisma](https://www.prisma.io/)
- [RESTful API](https://restfulapi.net/)
- [NGINX Web Server](https://en.wikipedia.org/wiki/Nginx)
- [Redis](https://redis.io/)
- [Redis Commander](https://github.com/joeferner/redis-commander)
- [PostgreSQL](https://www.postgresql.org)
- [Adminer](https://www.adminer.org)
- [Firebase Authentication](https://firebase.google.com/docs/auth)
- [Docker](https://www.docker.com)
- [Docker Compose](https://docs.docker.com/compose)

# Details

## Architecture

![Communication Server Structure](https://raw.githubusercontent.com/TP-OG/communication-server/main/docs/img/architecture.jpg)

## Database Design

![Communication Server Database Design](https://raw.githubusercontent.com/TP-OG/communication-server/main/docs/img/database.jpg)

## Documentaions

RESTful API [here](https://tp-og.github.io/communication-server/api-docs).

Event-Driven API [here](https://tp-og.github.io/communication-server/socketio-docs).

Application [here](https://tp-og.github.io/communication-server/app-docs).

# Setup

```bash
$ git clone git@github.com:TP-OG/communication-server.git

$ cd communication-server

$ cp .env.example .env
```

Then fill in the `.env` file.

## Development

```bash
$ docker-compose up

$ docker-compose exec app npx prisma migrate dev
```

## Demo

```bash
$ docker-compose -f docker-compose.demo.yml up

$ docker-compose exec app npx prisma migrate deploy
```

# License

- ##### This project is distributed under the [MIT License](LICENSE).
- ##### Copyright of [@TP-O](https://github.com/TP-O), 2022.
