import {
  BadRequestException,
  ForbiddenException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import Redis, { ChainableCommander } from 'ioredis';
import { Room, RoomId } from './room.type';
import { RedisService } from '../common/service';
import { PlayerId } from '../player/player.type';
import { RedisNamespace } from '../common/enum/redis.enum';
import * as randomstring from 'randomstring';
import merge from 'just-merge';

const ROOM_ID_LENGTH = 8;

@Injectable()
export class RoomService {
  private readonly _redis: Redis;

  constructor(redisService: RedisService) {
    this._redis = redisService.client;
  }

  /**
   * Get room by ID.
   *
   * @param id The room ID.
   */
  async get(id: RoomId): Promise<Room | null> {
    const roomJson = await this._redis.get(`${RedisNamespace.Room}${id}`);
    if (!roomJson) {
      return null;
    }

    return JSON.parse(roomJson);
  }

  /**
   * Get rooms by ID.
   *W
   * @param ids The list of room ID.
   */
  async getMany(ids: RoomId[]) {
    const rooms: Room[] = [];

    const pipe = this._redis.pipeline();
    ids.forEach((id) => pipe.get(`${RedisNamespace.Room}${id}`));
    (await pipe.exec())?.forEach(([err, json]) => {
      if (!err && typeof json === 'string') {
        rooms.push(JSON.parse(json));
      }
    });

    return rooms;
  }

  /**
   * Get list of room ID containing the given player ID.
   *
   * @param id The player ID.
   */
  getRoomIdsOfPlayer(id: PlayerId): Promise<RoomId[]> {
    return this._redis.lrange(`${RedisNamespace.Id2Rids}${id}`, 0, -1);
  }

  /**
   * Create the room.
   *
   * @param room The created room.
   */
  async create(room: Omit<Room, 'id'>): Promise<Room> {
    const id = randomstring.generate(ROOM_ID_LENGTH);
    const pipe = this._redis
      .pipeline()
      .set(`${RedisNamespace.Room}${id}`, JSON.stringify({ ...room, id }));
    room.memberIds.forEach((mId) =>
      pipe.lpush(`${RedisNamespace.Id2Rids}${mId}`, id),
    );

    await pipe.exec();
    return { ...room, id };
  }

  /**
   * Create new rooms and replace the old one if available.
   *
   * @param rooms The list of new room.
   */
  async forceCreateMany(rooms: Room[]): Promise<Room[]> {
    const pipe = this._redis.pipeline();
    rooms.forEach((room) => {
      room.memberIds.forEach((mId) => {
        pipe.lpush(`${RedisNamespace.Id2Rids}${mId}`, room.id);
      });
      pipe.set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room));
    });

    await pipe.exec();
    return rooms;
  }

  /**
   * Remove many rooms.
   *
   * @param ids The list of removed room ID.
   */
  async removeMany(ids: RoomId[]): Promise<Room[]> {
    const rooms: Room[] = (
      (
        await this._redis.mget(
          ...ids.map((id) => `${RedisNamespace.Room}${id}`),
        )
      ).filter((roomJson) => roomJson) as string[]
    ).map((roomJson) => JSON.parse(roomJson));

    const pipe = this._redis.pipeline();
    rooms.forEach((room) => {
      room.memberIds.forEach((id) => {
        pipe.lrem(`${RedisNamespace.Id2Rids}${id}`, 1, room.id);
      });
      pipe.del(`${RedisNamespace.Room}${room.id}`);
    });

    await pipe.exec();
    return rooms;
  }

  /**
   * Add member to room.
   *
   * @param id The room ID.
   * @param playerID The added player ID.
   * @param password The room password.
   */
  async join(id: RoomId, memberId: PlayerId, password?: string): Promise<Room> {
    const room = await this.get(id);
    if (!room) {
      throw new NotFoundException("Room doesn't exist!");
    }

    if (room.memberIds.includes(memberId)) {
      return room;
    }

    if (room.password && room.password !== password) {
      throw new BadRequestException('Incorrect password!');
    }

    room.memberIds.push(memberId);
    await this._redis
      .pipeline()
      .lpush(`${RedisNamespace.Id2Rids}${memberId}`, room.id)
      .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room))
      .exec();

    return room;
  }

  /**
   * Add members to room.
   *
   * @param id The room ID.
   * @param memberIds The list of added player ID.
   */
  async forceAddMembers(id: RoomId, memberIds: PlayerId[]): Promise<Room> {
    const room = await this.get(id);
    if (!room) {
      throw new NotFoundException("Room doesn't exist!");
    }

    const pipe = this._redis.pipeline();
    memberIds.forEach((mid) => {
      if (room.memberIds.includes(mid)) {
        return;
      }

      room.memberIds.push(mid);
      pipe.lpush(`${RedisNamespace.Id2Rids}${mid}`, room.id);
    });

    await pipe
      .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room))
      .exec();
    return room;
  }

  /**
   * Remove members from room.
   *
   * @param id The room ID.
   * @param memberIds The list of removed player ID.
   * @param kickerIds The kicker player ID.
   */
  async removeMembers(
    id: RoomId,
    memberIds: PlayerId[],
    kickerId?: PlayerId,
  ): Promise<Room> {
    const room = await this.get(id);
    if (!room) {
      throw new NotFoundException("Room doesn't exist!");
    }

    if (kickerId && room.ownerId !== kickerId) {
      throw new BadRequestException('Only owner is able to kick members');
    }

    const pipe = this._redis.pipeline();
    this._createRemoveMembersPipeline(pipe, room, memberIds);

    await pipe.exec();
    return room;
  }

  /**
   * Remove member from many rooms.
   *
   * @param ids The list of room ID. All rooms if empty.
   * @param memberId The removed member ID.
   */
  async removeFromRooms(ids: RoomId[], memberId: PlayerId): Promise<Room[]> {
    // Remove from all rooms if ids is empty
    if (ids.length === 0) {
      ids = await this._redis.lrange(
        `${RedisNamespace.Id2Rids}${memberId}`,
        0,
        -1,
      );
      if (ids.length === 0) {
        return [];
      }
    }

    const pipe = this._redis.pipeline();
    const rooms = (await this.getMany(ids)).map((room) => {
      this._createRemoveMembersPipeline(pipe, room, [memberId]);
      return room;
    });

    await pipe.exec();
    return rooms;
  }

  /**
   * Create redis pipeline to remove members from room.
   *
   * @param pipe The redis pipeline instance.
   * @param room The room.
   * @param memberIds The list of removed member ID.
   */
  private _createRemoveMembersPipeline(
    pipe: ChainableCommander,
    room: Room,
    memberIds: PlayerId[],
  ): ChainableCommander {
    if (!memberIds.some((mid) => memberIds.includes(mid))) {
      return pipe;
    }

    memberIds.forEach((mid) => {
      const removedMemberIndex = room.memberIds.indexOf(mid);
      if (removedMemberIndex === -1) {
        return;
      }

      room.memberIds.splice(removedMemberIndex, 1);
      pipe.lrem(`${RedisNamespace.Id2Rids}${mid}`, 1, room.id);

      // Set new owner if the owner is removed
      if (mid === room.ownerId && room.memberIds.length) {
        room.ownerId = room.memberIds[0];
      }
    });

    // Delete room if it is empty
    if (!room.memberIds.length) {
      pipe.del(`${RedisNamespace.Room}${room.id}`);
    } else {
      pipe.set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room));
    }

    return pipe;
  }

  /**
   * Mute the room.
   *
   * @param id The room ID.
   * @param isMuted Mute or unmute.
   */
  async mute(id: RoomId, isMuted = true): Promise<Room> {
    const room = await this.get(id);
    if (!room) {
      throw new NotFoundException("Room doesn't exist!");
    }

    if (room.isMuted === isMuted) {
      return room;
    }

    return this._mergeAndStoreRoom(room, { isMuted });
  }

  /**
   * Transfer ownership to another member. Decline action
   * if room is empty, actor is not owner, or choosed member
   * does not exist in the room.
   *
   * @param id The room ID.
   * @param newOwnerId The new owner player ID.
   * @param currentOwnerId The current owner player ID. Force set owner if ignore.
   */
  async transferOwnership(
    id: RoomId,
    newOwnerId: PlayerId,
    currentOwnerId?: PlayerId,
  ): Promise<Room> {
    const room = await this.get(id);
    if (!room) {
      throw new NotFoundException("Room doesn't exist!");
    }

    if (!currentOwnerId) {
      currentOwnerId = room.ownerId;
    }

    if (room.ownerId !== currentOwnerId) {
      throw new ForbiddenException('You are not owner of this room!');
    }

    if (currentOwnerId === newOwnerId) {
      return room;
    }

    if (!room.memberIds.includes(newOwnerId)) {
      throw new BadRequestException('New owner must be a member in this room!');
    }

    return this._mergeAndStoreRoom(room, { ownerId: newOwnerId });
  }

  /**
   * Merge new data into the room.
   *
   * @param room The updated room.
   * @param data The new data.
   */
  private async _mergeAndStoreRoom(
    room: Room,
    data: Partial<Omit<Room, 'id'>>,
  ): Promise<Room> {
    merge(room, data);
    await this._redis.set(
      `${RedisNamespace.Room}${room.id}`,
      JSON.stringify(room),
    );

    return room;
  }
}
