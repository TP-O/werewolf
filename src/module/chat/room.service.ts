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

  private async isMemeber(memberId: number) {
    const roomIds = await this.redis.llen(
      `${CacheNamespace.UId2RIds}${memberId}`,
    );

    return roomIds > 0;
  }

  async bookRoom(ownerId: number) {
    if (!AppConfig.allowJoinMultipleRooms && (await this.isMemeber(ownerId))) {
      throw new WsException(
        'Please leave current room before creating a new one',
      );
    }

    const id = String(Date.now());
    const room: Room = {
      id,
      memberIds: [ownerId],
    };

    await this.redis
      .pipeline()
      .set(`${CacheNamespace.Room}${id}`, JSON.stringify(room))
      .lpush(`${CacheNamespace.UId2RIds}`, room.id)
      .exec();

    return room;
  }

  async getRoom(id: string) {
    const roomJson = await this.redis.get(`${CacheNamespace.Room}${id}`);

    if (roomJson === null) {
      throw new WsException('Room does not exist!');
    }

    const room: Room = JSON.parse(roomJson);

    return room;
  }

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

  async leaveRoom(id: string, leaverId: number) {
    const room = await this.getRoom(id);
    const deletedMemeberIndex = room.memberIds.indexOf(leaverId);

    if (deletedMemeberIndex === -1) {
      throw new WsException('Join before leaving a room!');
    }

    const redisPipe = this.redis.pipeline();

    // Delete room if all members have left
    if (room.memberIds.length === 0) {
      redisPipe.del(`${CacheNamespace.Room}${id}`);
    } else {
      room.memberIds = room.memberIds.filter((id) => id !== leaverId);
      redisPipe.set(`${CacheNamespace.Room}${id}`, JSON.stringify(room));
    }

    await redisPipe.lrem(`${CacheNamespace.UId2RIds}`, 1, id).exec();
  }
}
