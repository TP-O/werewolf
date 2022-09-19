import { Injectable } from '@nestjs/common';
import { WsException } from '@nestjs/websockets';
import Redis from 'ioredis';
import { AppConfig } from 'src/config/app.config';
import { RedisClient } from 'src/decorator/redis.decorator';
import { CacheNamespace } from 'src/enum/cache.enum';
import { Room } from './type/room.type';

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
  private async isMemeberOfAny(memberId: number) {
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
   * @returns room value.
   */
  async book(bookerId: number) {
    if (
      !AppConfig.allowJoinMultipleRooms &&
      (await this.isMemeberOfAny(bookerId))
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
      rejectedIds: [],
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
   * Join to a new room. If multi-room join is disabled,
   * the booker must not enter any room before creating
   * the room.
   *
   * @param roomId
   * @param joinerId
   * @returns updated room value.
   */
  async join(roomId: string, joinerId: number) {
    if (
      !AppConfig.allowJoinMultipleRooms &&
      (await this.isMemeberOfAny(joinerId))
    ) {
      throw new WsException(
        'Please leave current room before joining another one!',
      );
    }

    const room = await this.get(roomId);
    room.memberIds.push(joinerId);

    await this.redis
      .pipeline()
      .set(`${CacheNamespace.Room}${room.id}`, JSON.stringify(room))
      .lpush(`${CacheNamespace.UId2RIds}`, room.id)
      .exec();

    return room;
  }

  /**
   * Leave the room. Delete the room if it is empty.
   *
   * @param roomId
   * @param leaverId
   * @returns updated room value.
   */
  async leave(roomId: string, leaverId: number) {
    const room = await this.get(roomId);
    const deletedMemeberIndex = room.memberIds.indexOf(leaverId);

    if (deletedMemeberIndex === -1) {
      throw new WsException('You are not in this room!');
    } else {
      room.memberIds.splice(deletedMemeberIndex, 1);
      room.rejectedIds.push(leaverId);
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
   * Kick memeber out of room.
   *
   * @param roomId
   * @param kickerId
   * @param memberId
   * @returns updated room value.
   */
  async kick(roomId: string, kickerId: number, memberId: number) {
    const room = await this.get(roomId);

    if (room.ownerId !== kickerId) {
      throw new WsException('You are not owner of this room!');
    }

    const deletedMemeberIndex = room.memberIds.indexOf(memberId);

    if (deletedMemeberIndex === -1) {
      throw new WsException('Member is not in this room!');
    }

    room.memberIds.splice(deletedMemeberIndex, 1);
    room.rejectedIds.push(memberId);

    await this.redis
      .pipeline()
      .set(`${CacheNamespace.Room}${room.id}`, JSON.stringify(room))
      .lrem(`${CacheNamespace.UId2RIds}`, 1, room.id)
      .exec();

    return room;
  }

  /**
   *
   * @param roomId
   * @param transfererId
   * @param candidateId
   * @returns
   */
  async transferOwnership(
    roomId: string,
    transfererId: number,
    candidateId: number,
  ) {
    const room = await this.get(roomId);

    if (room.ownerId !== transfererId) {
      throw new WsException('You are not owner of this room!');
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
   *
   * @param roomId
   * @param inviter
   * @param guestId
   * @returns
   */
  async invite(roomId: string, inviter: number, guestId: number) {
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
    room.rejectedIds = room.rejectedIds.filter((id) => id !== guestId);

    await this.redis.set(
      `${CacheNamespace.Room}${room.id}`,
      JSON.stringify(room),
    );

    return { room, guestSIds };
  }

  /**
   *
   * @param roomId
   * @param guestId
   * @param isAccpeted
   * @returns
   */
  async replyInvitation(roomId: string, guestId: number, isAccpeted: boolean) {
    const room = await this.get(roomId);
    const deletedWaitingIndex = room.waitingIds.indexOf(guestId);

    if (deletedWaitingIndex === -1) {
      throw new WsException('You are not invited to this room!');
    }

    room.waitingIds.splice(deletedWaitingIndex, 1);

    if (isAccpeted) {
      room.memberIds.push(guestId);
    } else {
      room.rejectedIds.push(guestId);
    }

    await this.redis.set(
      `${CacheNamespace.Room}${room.id}`,
      JSON.stringify(room),
    );

    return room;
  }
}
