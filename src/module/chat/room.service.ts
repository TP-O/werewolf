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
  private async isMemeber(memberId: number) {
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
   * @param ownerId
   * @returns room value.
   */
  async bookRoom(ownerId: number) {
    if (!AppConfig.allowJoinMultipleRooms && (await this.isMemeber(ownerId))) {
      throw new WsException(
        'Please leave current room before creating a new one',
      );
    }

    const id = String(Date.now());
    const room: Room = {
      id,
      ownerId,
      memberIds: [ownerId],
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
   * @param id
   * @returns
   */
  async getRoom(id: string) {
    const roomJson = await this.redis.get(`${CacheNamespace.Room}${id}`);

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
   * @param id room id.
   * @param joinerId
   * @returns updated room value.
   */
  async joinRoom(id: string, joinerId: number) {
    if (!AppConfig.allowJoinMultipleRooms && (await this.isMemeber(joinerId))) {
      throw new WsException(
        'Please leave current room before joining another one!',
      );
    }

    const room = await this.getRoom(id);
    room.memberIds.push(joinerId);

    await this.redis
      .pipeline()
      .set(`${CacheNamespace.Room}${id}`, JSON.stringify(room))
      .lpush(`${CacheNamespace.UId2RIds}`, room.id)
      .exec();

    return room;
  }

  /**
   * Leave the room. Delete the room if it is empty.
   *
   * @param id room id.
   * @param leaverId
   * @returns updated room value.
   */
  async leaveRoom(id: string, leaverId: number) {
    const room = await this.getRoom(id);
    const deletedMemeberIndex = room.memberIds.indexOf(leaverId);

    if (deletedMemeberIndex === -1) {
      throw new WsException('You are not in this room!');
    } else {
      room.memberIds.splice(deletedMemeberIndex, 1);
    }

    const redisPipe = this.redis.pipeline();

    // Delete room if all members have left
    if (room.memberIds.length === 0) {
      redisPipe.del(`${CacheNamespace.Room}${id}`);
    } else {
      redisPipe.set(`${CacheNamespace.Room}${id}`, JSON.stringify(room));
    }

    await redisPipe.lrem(`${CacheNamespace.UId2RIds}`, 1, id).exec();

    return room;
  }

  /**
   * Kick memeber out of room.
   *
   * @param id room id.
   * @param kickerId
   * @param memberId
   * @returns updated room value.
   */
  async kickMember(id: string, kickerId: number, memberId: number) {
    const room = await this.getRoom(id);

    if (room.ownerId !== kickerId) {
      throw new WsException('You are not owner of this room!');
    }

    const deletedMemeberIndex = room.memberIds.indexOf(memberId);

    if (deletedMemeberIndex === -1) {
      throw new WsException('You are not in this room!');
    }

    room.memberIds.splice(deletedMemeberIndex, 1);

    await this.redis
      .pipeline()
      .set(`${CacheNamespace.Room}${id}`, JSON.stringify(room))
      .lrem(`${CacheNamespace.UId2RIds}`, 1, id)
      .exec();

    return room;
  }
}
