import {
  BadRequestException,
  ForbiddenException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import Redis from 'ioredis';
import { Room } from './room.type';
import { v4 as uuidv4 } from 'uuid';
import { RedisService } from '../common';
import { PlayerId } from '../player';
import { RedisNamespace } from 'src/common/enum';
import { CreatePersistentRoomsDto, CreateTemporaryRoomsDto } from './dto';

@Injectable()
export class RoomService {
  private readonly _redis: Redis;

  constructor(redisService: RedisService) {
    this._redis = redisService.client;
  }

  /**
   * Check if player is in any room.
   *
   * @param memberId
   * @returns
   */
  private async isMemberOfAny(memberId: PlayerId) {
    const roomIds = await this._redis.llen(
      `${RedisNamespace.UId2RIds}${memberId}`,
    );

    return roomIds > 0;
  }

  /**
   * Store rooms temporarily.
   *
   * @param rooms
   * @returns created rooms and socket ids of members.
   */
  private async storeRooms(rooms: Room[]) {
    const socketIdsList: string[][] = [];
    const joinerIdsList: PlayerId[][] = [];
    const redisPipe = this._redis.pipeline();

    rooms.map((room) => {
      const memberSIdKeys = room.memberIds.map((mId) => {
        // By the way, bind room id with member id
        redisPipe.lpush(`${RedisNamespace.UId2RIds}${mId}`, room.id);

        return `${RedisNamespace.UID2SId}${mId}`;
      });

      joinerIdsList.push(room.memberIds);
      redisPipe
        .mget(...memberSIdKeys)
        .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room));

      return room;
    });

    (await redisPipe.exec()).forEach((res) => {
      if (Array.isArray(res[1])) {
        socketIdsList.push(res[1]);
      }
    });

    return {
      rooms,
      socketIdsList,
      joinerIdsList,
    };
  }

  /**
   * Create temporary rooms with given settings.
   *
   * @param dto
   * @returns created rooms and socket ids of members.
   */
  async createTemporarily(dto: CreateTemporaryRoomsDto) {
    return this.storeRooms(
      dto.rooms.map((setting) => ({
        ...setting,
        id: uuidv4(),
      })),
    );
  }

  /**
   * Create persistent rooms with given settings.
   *
   * @param dto
   * @returns created rooms and socket ids of members.
   */
  async createPersistently(dto: CreatePersistentRoomsDto) {
    return this.storeRooms(
      dto.rooms.map((setting) => ({
        ...setting,
        id: `${dto.gameId}:${setting.id}`,
        gameId: dto.gameId,
      })),
    );
  }

  /**
   * Remove many rooms by given room ids.
   *
   * @param roomIds
   * @returns removed rooms and socket ids of members in removed rooms.
   */
  async remove(roomIds: string[]) {
    const rooms: Room[] = (
      await this._redis.mget(
        ...roomIds.map((rId) => `${RedisNamespace.Room}${rId}`),
      )
    )
      .filter((roomJSON) => roomJSON != null)
      .map((roomJSON) => JSON.parse(roomJSON));
    const socketIdsList: string[][] = [];
    const leaverIdsList: PlayerId[][] = [];
    const redisPipe = this._redis.pipeline();

    rooms.forEach((room) => {
      const memberSIdKeys = room.memberIds.map((mId) => {
        // By the way, remove room id from member id
        redisPipe.lrem(`${RedisNamespace.UId2RIds}${mId}`, 1, room.id);

        return `${RedisNamespace.UID2SId}${mId}`;
      });

      leaverIdsList.push(room.memberIds);
      redisPipe.mget(...memberSIdKeys).del(`${RedisNamespace.Room}${room.id}`);
    });

    (await redisPipe.exec()).forEach((res) => {
      if (Array.isArray(res[1])) {
        socketIdsList.push(res[1]);
      }
    });

    return { rooms, socketIdsList, leaverIdsList };
  }

  /**
   * Add members to room.
   *
   * @param roomId
   * @param memberIds
   * @returns updated room and socket id of members.
   */
  async addMembers(roomId: string, memberIds: PlayerId[]) {
    const room = await this.get(roomId);
    const redisPipe = this._redis.pipeline();
    const memberSIdKeys: string[] = [];

    memberIds.forEach((mId) => {
      if (!room.memberIds.includes(mId)) {
        room.memberIds.push(mId);
        memberSIdKeys.push(`${RedisNamespace.UID2SId}${mId}`);
        redisPipe.lpush(`${RedisNamespace.UId2RIds}${mId}`, room.id);
      }
    });

    const redisRes = await redisPipe
      .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room))
      .mget(...memberSIdKeys)
      .exec();

    return {
      room,
      socketIds: (redisRes[redisPipe.length - 1][1] as string[]) || [],
    };
  }

  /**
   * Remove members from room.
   *
   * @param roomId
   * @param memberIds
   * @returns updated room and socket id of members.
   */
  async removeMembers(roomId: string, memberIds: PlayerId[]) {
    const room = await this.get(roomId);
    const redisPipe = this._redis.pipeline();
    const memberSIdKeys: string[] = [];

    memberIds.forEach((mId) => {
      const removedMemberIndex = room.memberIds.indexOf(mId);

      if (removedMemberIndex !== -1) {
        room.memberIds.splice(removedMemberIndex, 1);
        memberSIdKeys.push(`${RedisNamespace.UID2SId}${mId}`);
        redisPipe.lrem(`${RedisNamespace.UId2RIds}${mId}`, 1, room.id);
      }
    });

    const redisRes = await redisPipe
      .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room))
      .mget(...memberSIdKeys)
      .exec();

    return {
      room,
      socketIds: (redisRes[redisPipe.length - 1][1] as string[]) || [],
    };
  }

  /**
   * Allow chatting in room or not.
   *
   * @param roomId
   * @param mute
   * @returns updated room.
   */
  async allowChat(roomId: string, mute = true) {
    const room = await this.get(roomId);

    if (room.isMuted !== mute) {
      room.isMuted = mute;
    }

    await this._redis.set(
      `${RedisNamespace.Room}${roomId}`,
      JSON.stringify(room),
    );

    return room;
  }

  /**
   * Create a room and add the booker to its member list.
   * If multi-room join is disabled, the booker must not
   * enter any room before creating the room.
   *
   * @param bookerId
   * @param isPublic if true, anyone can join without invitation.
   * @returns updated room.
   */
  async book(bookerId: PlayerId, isPublic = false) {
    if (await this.isMemberOfAny(bookerId)) {
      throw new BadRequestException(
        'Please leave current room before creating a new one',
      );
    }

    const room: Room = {
      id: uuidv4(),
      isPublic,
      isPersistent: false,
      isMuted: false,
      gameId: 0,
      ownerId: bookerId,
      memberIds: [bookerId],
      waitingIds: [],
      refusedIds: [],
    };

    await this._redis
      .pipeline()
      .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room))
      .lpush(`${RedisNamespace.UId2RIds}${bookerId}`, room.id)
      .exec();

    return room;
  }

  /**
   * Get room by id.
   *
   * @param roomId
   * @returns
   */
  async get(roomId: string) {
    const roomJSON = await this._redis.get(`${RedisNamespace.Room}${roomId}`);

    if (roomJSON === null) {
      throw new BadRequestException('Room does not exist!');
    }

    const room: Room = JSON.parse(roomJSON);

    return room;
  }

  /**
   * Get rooms by ids.
   *
   * @param roomIds
   * @returns
   */
  async getMany(roomIds: string[]) {
    const rooms: Room[] = [];
    const redisPipe = this._redis.pipeline();
    roomIds.forEach((rId) => redisPipe.get(`${RedisNamespace.Room}${rId}`));

    (await redisPipe.exec()).forEach((val) => {
      if (typeof val[1] === 'string') {
        rooms.push(JSON.parse(val[1] as string));
      }
    });

    return rooms;
  }

  /**
   * Add player to room. If multi-room join is disabled,
   * the booker must not enter any room before creating
   * the room.
   *
   * @param joinerId
   * @param roomId
   * @returns updated room.
   */
  async join(joinerId: PlayerId, roomId: string) {
    if (await this.isMemberOfAny(joinerId)) {
      throw new BadRequestException(
        'Please leave current room before joining another one!',
      );
    }

    const room = await this.get(roomId);

    if (!room.isPublic || room.isPersistent) {
      throw new ForbiddenException('This room is private!');
    }

    if (room.memberIds.includes(joinerId)) {
      throw new BadRequestException('You have already joined this room!');
    }

    room.memberIds.push(joinerId);
    room.waitingIds.splice(room.waitingIds.indexOf(joinerId), 1);
    room.refusedIds.splice(room.waitingIds.indexOf(joinerId), 1);

    await this._redis
      .pipeline()
      .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room))
      .lpush(`${RedisNamespace.UId2RIds}${joinerId}`, room.id)
      .exec();

    return room;
  }

  /**
   * Remove player from room. Transfer ownership for to a member
   * in room if leaver is owner. Empty room will be deleted.
   *
   * @param leaverId
   * @param roomId
   * @returns updated room.
   */
  async leave(leaverId: PlayerId, roomId: string) {
    const room = await this.get(roomId);
    const deletedMemberIndex = room.memberIds.indexOf(leaverId);

    if (deletedMemberIndex === -1) {
      throw new ForbiddenException('You are not in this room!');
    } else {
      room.memberIds.splice(deletedMemberIndex, 1);
      room.refusedIds.push(leaverId);
    }

    const redisPipe = this._redis.pipeline();

    // Delete room if all members have left
    if (room.memberIds.length === 0 && !room.isPersistent) {
      redisPipe.del(`${RedisNamespace.Room}${room.id}`);
    } else {
      // Assign owner to the first member
      if (leaverId === room.ownerId) {
        room.ownerId = room.memberIds[0];
      }

      redisPipe.set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room));
    }

    await redisPipe
      .lrem(`${RedisNamespace.UId2RIds}${leaverId}`, 1, room.id)
      .exec();

    return room;
  }

  /**
   * Remove player from many rooms. Logic is the same
   * as `leave` method.
   *
   * @param leaverId
   * @param roomIds empty if leaving all rooms.
   * @returns updated rooms.
   */
  async leaveMany(leaverId: PlayerId, ...roomIds: string[]) {
    if (roomIds.length === 0) {
      roomIds = await this._redis.lrange(
        `${RedisNamespace.UId2RIds}${leaverId}`,
        0,
        -1,
      );

      if (roomIds.length === 0) {
        return [];
      }
    }

    let rooms = await this.getMany(roomIds);
    const redisPipe = this._redis.pipeline();

    rooms = rooms.map((room) => {
      room.memberIds.splice(room.memberIds.indexOf(leaverId), 1);
      room.refusedIds.push(leaverId);

      // Delete room if all members have left
      if (room.memberIds.length === 0 && !room.isPersistent) {
        redisPipe.del(`${RedisNamespace.Room}${room.id}`);
      } else {
        // Assign owner to the first member
        if (leaverId === room.ownerId) {
          room.ownerId = room.memberIds[0];
        }

        redisPipe.set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room));
      }

      redisPipe
        .lrem(`${RedisNamespace.UId2RIds}${leaverId}`, 1, room.id)
        .exec();

      return room;
    });

    await redisPipe.exec();

    return rooms;
  }

  /**
   * Kick member out of room. Kicker must be the owner
   * and member must be in the room, otherwise the action
   * is declined.
   *
   * @param kickerId
   * @param memberId
   * @param roomId
   * @returns updated room and socket id of kicked member.
   */
  async kick(kickerId: PlayerId, memberId: PlayerId, roomId: string) {
    if (kickerId === memberId) {
      const room = await this.leave(kickerId, roomId);

      return { room, kickedMemberSId: kickerId };
    }

    const room = await this.get(roomId);

    if (room.ownerId !== kickerId) {
      throw new ForbiddenException('You are not owner of this room!');
    }

    const deletedMemberIndex = room.memberIds.indexOf(memberId);

    if (deletedMemberIndex === -1) {
      throw new BadRequestException('Member is not in this room!');
    }

    room.memberIds.splice(deletedMemberIndex, 1);
    room.refusedIds.push(memberId);

    const [[, kickedMemberSId]] = await this._redis
      .pipeline()
      .get(`${RedisNamespace.UID2SId}${memberId}`)
      .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room))
      .lrem(`${RedisNamespace.UId2RIds}${memberId}`, 1, room.id)
      .exec();

    return {
      room,
      kickedMemberSocketId: kickedMemberSId as string,
    };
  }

  /**
   * Transfer ownership to another member. Decline action
   * if room is empty, actor is not owner, or choosed member
   * does not exist in the room.
   *
   * @param ownerId
   * @param candidateId
   * @param roomId
   * @returns update room.
   */
  async transferOwnership(
    ownerId: PlayerId,
    candidateId: PlayerId,
    roomId: string,
  ) {
    const room = await this.get(roomId);

    if (room.ownerId !== ownerId) {
      throw new ForbiddenException('You are not owner of this room!');
    }

    if (ownerId === candidateId || !room.memberIds.includes(candidateId)) {
      throw new BadRequestException('New owner must be a member in this room!');
    }

    room.ownerId = candidateId;

    await this._redis.set(
      `${RedisNamespace.Room}${room.id}`,
      JSON.stringify(room),
    );

    return room;
  }

  /**
   * Invite a guest into room. Only invite online player and
   * non-exist in room player.
   *
   * @param inviter
   * @param guestId
   * @param roomId
   * @returns updated room and guest socket ids.
   */
  async invite(inviter: PlayerId, guestId: PlayerId, roomId: string) {
    const [[, roomJSON], [, guestSId]] = (await this._redis
      .pipeline()
      .get(`${RedisNamespace.Room}${roomId}`)
      .get(`${RedisNamespace.UID2SId}${guestId}`)
      .exec()) as [error: any, result: string | string[]][];

    if (roomJSON == null) {
      throw new NotFoundException('Room does not exist!');
    }

    const room: Room = JSON.parse(roomJSON as string);

    if (room.isPersistent) {
      throw new ForbiddenException(
        'You can not invite other player to this room!',
      );
    }

    if (guestSId == null) {
      throw new BadRequestException('Please only invite online player!');
    }

    if (!room.memberIds.includes(inviter)) {
      throw new ForbiddenException('You are not in this room!');
    }

    if (room.memberIds.includes(guestId) || room.waitingIds.includes(guestId)) {
      throw new BadRequestException('This player has been invited!');
    }

    room.waitingIds.push(guestId);
    room.refusedIds.splice(room.waitingIds.indexOf(guestId), 1);

    await this._redis.set(
      `${RedisNamespace.Room}${room.id}`,
      JSON.stringify(room),
    );

    return { room, guestSocketId: guestSId };
  }

  /**
   * Respond to room invitation. There are 2 options:
   * accept and refuse. Leave the current room after
   * accepting if multi-room join is disabled.
   *
   * @param guestId
   * @param isAccpeted
   * @param roomId
   * @returns updated room and left rooms.
   */
  async respondInvitation(
    guestId: PlayerId,
    isAccpeted: boolean,
    roomId: string,
  ) {
    let leftRooms: Room[] = [];
    const room = await this.get(roomId);
    const deletedWaitingIndex = room.waitingIds.indexOf(guestId);

    if (deletedWaitingIndex === -1) {
      throw new BadRequestException('You are not invited to this room!');
    }

    const redisPipe = this._redis.pipeline();

    room.waitingIds.splice(deletedWaitingIndex, 1);

    if (isAccpeted) {
      // Leave all current rooms if multi-room join is disabled
      if (await this.isMemberOfAny(guestId)) {
        leftRooms = await this.leaveMany(guestId);
      }

      room.memberIds.push(guestId);
      redisPipe.lpush(`${RedisNamespace.UId2RIds}${guestId}`, room.id);
    } else {
      room.refusedIds.push(guestId);
    }

    await redisPipe
      .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room))
      .exec();

    return { room, leftRooms };
  }
}
