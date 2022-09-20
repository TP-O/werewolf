import { Injectable } from '@nestjs/common';
import { WsException } from '@nestjs/websockets';
import Redis from 'ioredis';
import { AppConfig } from 'src/config';
import { RedisClient } from 'src/decorator';
import { CacheNamespace } from 'src/enum';
import { Room } from '../type';

@Injectable()
export class RoomService {
  @RedisClient()
  private readonly redis: Redis;

  /**
   * Check if user is in any room.
   *
   * @param memberId
   * @returns
   */
  private async isMemberOfAny(memberId: number) {
    const roomIds = await this.redis.llen(
      `${CacheNamespace.UId2RIds}${memberId}`,
    );

    return roomIds > 0;
  }

  /**
   * Create a room and add the booker to its member list.
   * If multi-room join is disabled, the booker must not
   * enter any room before creating the room.
   *
   * @param bookerId
   * @returns updated room.
   */
  async book(bookerId: number) {
    if (
      !AppConfig.allowJoinMultipleRooms &&
      (await this.isMemberOfAny(bookerId))
    ) {
      throw new WsException(
        'Please leave current room before creating a new one',
      );
    }

    const id = String(Date.now());
    const room: Room = {
      id,
      ownerId: bookerId,
      memberIds: [bookerId],
      waitingIds: [],
      refusedIds: [],
    };

    await this.redis
      .pipeline()
      .set(`${CacheNamespace.Room}${id}`, JSON.stringify(room))
      .lpush(`${CacheNamespace.UId2RIds}`, room.id)
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
    const roomJson = await this.redis.get(`${CacheNamespace.Room}${roomId}`);

    if (roomJson === null) {
      throw new WsException('Room does not exist!');
    }

    const room: Room = JSON.parse(roomJson);

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
    const redisPipe = this.redis.pipeline();
    roomIds.forEach((rId) => redisPipe.get(`${CacheNamespace.Room}${rId}`));

    (await redisPipe.exec()).forEach((val) => {
      if (typeof val[1] === 'string') {
        rooms.push(JSON.parse(val[1] as string));
      }
    });

    return rooms;
  }

  /**
   * Join to a new room. If multi-room join is disabled,
   * the booker must not enter any room before creating
   * the room.
   *
   * @param joinerId
   * @param roomId
   * @returns updated room.
   */
  async join(joinerId: number, roomId: string) {
    if (
      !AppConfig.allowJoinMultipleRooms &&
      (await this.isMemberOfAny(joinerId))
    ) {
      throw new WsException(
        'Please leave current room before joining another one!',
      );
    }

    const room = await this.get(roomId);
    room.memberIds.push(joinerId);
    room.waitingIds.splice(room.waitingIds.indexOf(joinerId), 1);
    room.refusedIds.splice(room.waitingIds.indexOf(joinerId), 1);

    await this.redis
      .pipeline()
      .set(`${CacheNamespace.Room}${room.id}`, JSON.stringify(room))
      .lpush(`${CacheNamespace.UId2RIds}`, room.id)
      .exec();

    return room;
  }

  /**
   * Leave the room. Transfer ownership for to a member
   * in room if leaver is owner. Empty room will be deleted.
   *
   * @param leaverId
   * @param roomId
   * @returns updated room.
   */
  async leave(leaverId: number, roomId: string) {
    const room = await this.get(roomId);
    const deletedMemberIndex = room.memberIds.indexOf(leaverId);

    if (deletedMemberIndex === -1) {
      throw new WsException('You are not in this room!');
    } else {
      room.memberIds.splice(deletedMemberIndex, 1);
      room.refusedIds.push(leaverId);
    }

    const redisPipe = this.redis.pipeline();

    // Delete room if all members have left
    if (room.memberIds.length === 0) {
      redisPipe.del(`${CacheNamespace.Room}${room.id}`);
    } else {
      // Assign owner to the first member
      if (leaverId === room.ownerId) {
        room.ownerId = room.memberIds[0];
      }

      redisPipe.set(`${CacheNamespace.Room}${room.id}`, JSON.stringify(room));
    }

    await redisPipe.lrem(`${CacheNamespace.UId2RIds}`, 1, room.id).exec();

    return room;
  }

  /**
   * Leave many rooms. Logic is the same as `leave` method.
   *
   * @param leaverId
   * @param roomIds empty if leaving all rooms.
   * @returns updated rooms.
   */
  async leaveMany(leaverId: number, ...roomIds: string[]) {
    if (roomIds.length === 0) {
      roomIds = await this.redis.lrange(
        `${CacheNamespace.UId2RIds}${leaverId}`,
        0,
        -1,
      );
    }

    const rooms = await this.getMany(roomIds);
    const redisPipe = this.redis.pipeline();

    rooms.map((room) => {
      room.memberIds.splice(room.memberIds.indexOf(leaverId), 1);
      room.refusedIds.push(leaverId);

      // Delete room if all members have left
      if (room.memberIds.length === 0) {
        redisPipe.del(`${CacheNamespace.Room}${room.id}`);
      } else {
        // Assign owner to the first member
        if (leaverId === room.ownerId) {
          room.ownerId = room.memberIds[0];
        }

        redisPipe.set(`${CacheNamespace.Room}${room.id}`, JSON.stringify(room));
      }

      redisPipe.lrem(`${CacheNamespace.UId2RIds}`, 1, room.id).exec();

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
   * @returns updated room.
   */
  async kick(kickerId: number, memberId: number, roomId: string) {
    const room = await this.get(roomId);

    if (room.ownerId !== kickerId) {
      throw new WsException('You are not owner of this room!');
    }

    const deletedMemberIndex = room.memberIds.indexOf(memberId);

    if (deletedMemberIndex === -1) {
      throw new WsException('Member is not in this room!');
    }

    room.memberIds.splice(deletedMemberIndex, 1);
    room.refusedIds.push(memberId);

    await this.redis
      .pipeline()
      .set(`${CacheNamespace.Room}${room.id}`, JSON.stringify(room))
      .lrem(`${CacheNamespace.UId2RIds}`, 1, room.id)
      .exec();

    return room;
  }

  /**
   * Transfer ownership to another member. Decline action
   * if room is empty, actor is not owner, or choosed member
   * does not exist in the room.
   *
   * @param transfererId
   * @param candidateId
   * @param roomId
   * @returns update room.
   */
  async transferOwnership(
    transfererId: number,
    candidateId: number,
    roomId: string,
  ) {
    const room = await this.get(roomId);

    if (room.ownerId !== transfererId) {
      throw new WsException('You are not owner of this room!');
    }

    if (room.memberIds.length === 1) {
      throw new WsException('Your room is empty!');
    }

    if (!room.memberIds.includes(candidateId)) {
      throw new WsException('New owner must in this room!');
    }

    room.ownerId = candidateId;

    await this.redis.set(
      `${CacheNamespace.Room}${room.id}`,
      JSON.stringify(room),
    );

    return room;
  }

  /**
   * Invite a guest into room. Only invite online user and
   * non-exist in room user.
   *
   * @param inviter
   * @param guestId
   * @param roomId
   * @returns updated room.
   */
  async invite(inviter: number, guestId: number, roomId: string) {
    const [[, roomJson], [, guestSIds]] = (await this.redis
      .pipeline()
      .get(`${CacheNamespace.Room}${roomId}`)
      .lrange(`${CacheNamespace.UId2RIds}${guestId}`, 0, -1)
      .exec()) as [error: any, result: string | string[]][];
    const room: Room = JSON.parse(roomJson as string);

    if (guestSIds.length === 0) {
      throw new WsException('Please only invite online user!');
    }

    if (!room.memberIds.includes(inviter)) {
      throw new WsException('You are not in this room!');
    }

    if (room.memberIds.includes(guestId) || room.waitingIds.includes(guestId)) {
      throw new WsException('This user has been invited!');
    }

    room.waitingIds.push(guestId);
    room.refusedIds.splice(room.waitingIds.indexOf(guestId), 1);

    await this.redis.set(
      `${CacheNamespace.Room}${room.id}`,
      JSON.stringify(room),
    );

    return { room, guestSIds };
  }

  /**
   * Respond to room invitation. There are 2 options:
   * accept and refuse. Leave the current room after
   * accepting if multi-room join is disabled.
   *
   * @param guestId
   * @param isAccpeted
   * @param roomId
   * @returns update room.
   */
  async respondInvitation(
    guestId: number,
    isAccpeted: boolean,
    roomId: string,
  ) {
    const room = await this.get(roomId);
    const deletedWaitingIndex = room.waitingIds.indexOf(guestId);

    if (deletedWaitingIndex === -1) {
      throw new WsException('You are not invited to this room!');
    }

    room.waitingIds.splice(deletedWaitingIndex, 1);

    if (isAccpeted) {
      // Leave all current rooms if multi-room join is disabled
      if (
        !AppConfig.allowJoinMultipleRooms &&
        (await this.isMemberOfAny(guestId))
      ) {
        await this.leaveMany(guestId);
      }

      room.memberIds.push(guestId);
    } else {
      room.refusedIds.push(guestId);
    }

    await this.redis.set(
      `${CacheNamespace.Room}${room.id}`,
      JSON.stringify(room),
    );

    return room;
  }
}
