import {
  BadRequestException,
  ForbiddenException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import Redis, { ChainableCommander } from 'ioredis';
import { Room, RoomId } from './room.type';
import { RedisService } from '../common';
import { PlayerId } from '../player';
import { RedisNamespace } from '../common/enum/redis.enum';
import * as randomstring from 'randomstring';
import _merge from 'just-merge';

const ROOM_ID_LENGTH = 5;

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

    const redisPipe = this._redis.pipeline();
    ids.forEach((id) => redisPipe.get(`${RedisNamespace.Room}${id}`));
    (await redisPipe.exec())?.forEach(([err, json]) => {
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
   * Create the room..
   *
   * @param room The created room.
   */
  async create(room: Omit<Room, 'id'>): Promise<Room> {
    const id = randomstring.generate(ROOM_ID_LENGTH);
    await this._redis
      .pipeline()
      .set(`${RedisNamespace.Room}${id}`, JSON.stringify({ ...room, id }))
      .lpush(`${RedisNamespace.Id2Rids}${id}`, id)
      .exec();

    return { ...room, id };
  }

  /**
   * Create new rooms and replace the old one if available.
   *
   * @param rooms The list of new room.
   */
  async forceCreateMany(rooms: Room[]): Promise<Room[]> {
    const redisPipe = this._redis.pipeline();
    rooms.forEach((room) => {
      room.memberIds.forEach((id) => {
        redisPipe.lpush(`${RedisNamespace.Id2Rids}${id}`, room.id);
      });
      redisPipe.set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room));
    });
    await redisPipe.exec();

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

    const redisPipe = this._redis.pipeline();
    rooms.forEach((room) => {
      room.memberIds.forEach((id) => {
        redisPipe.lrem(`${RedisNamespace.Id2Rids}${id}`, 1, room.id);
      });
      redisPipe.del(`${RedisNamespace.Room}${room.id}`);
    });
    await redisPipe.exec();

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

    if (room.password && room.password !== password) {
      throw new BadRequestException('Incorrect password!');
    }

    const redisPipe = this._redis.pipeline();
    this._createAddMembersPipepline(redisPipe, room, [memberId]);
    await redisPipe
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

    const redisPipe = this._redis.pipeline();
    this._createAddMembersPipepline(redisPipe, room, memberIds);
    if (redisPipe.length) {
      await redisPipe
        .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room))
        .exec();
    }

    return room;
  }

  /**
   * Create redis pipeline to add members to room.
   *
   * @param pipe The redis pipeline instance.
   * @param room The room.
   * @param memberIds The list of added member ID.
   */
  _createAddMembersPipepline(
    pipe: ChainableCommander,
    room: Room,
    memberIds: PlayerId[],
  ): ChainableCommander {
    memberIds.forEach((mid) => {
      if (room.memberIds.includes(mid)) {
        throw new BadRequestException('Player is already in this room!');
      }

      room.memberIds.push(mid);
      pipe.lpush(`${RedisNamespace.Id2Rids}${mid}`, room.id);
    });

    return pipe;
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

    const redisPipe = this._redis.pipeline();
    this._createRemoveMembersPipeline(redisPipe, room, memberIds);
    if (redisPipe.length) {
      await redisPipe
        .set(`${RedisNamespace.Room}${room.id}`, JSON.stringify(room))
        .exec();
    }

    return room;
  }

  /**
   * Remove member from many rooms.
   *
   * @param ids The list of room ID. All rooms if empty.
   * @param memberId The removed member ID.
   */
  async removeFromRooms(ids: RoomId[], memberId: PlayerId): Promise<Room[]> {
    // Remove from all rooms if
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

    const redisPipe = this._redis.pipeline();
    const rooms = (await this.getMany(ids)).map((room) => {
      this._createRemoveMembersPipeline(redisPipe, room, [memberId]);
      return room;
    });

    if (redisPipe.length) {
      await redisPipe.exec();
    }

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
    memberIds.forEach((mid) => {
      const removedMemberIndex = room.memberIds.indexOf(mid);
      if (removedMemberIndex === -1) {
        throw new BadRequestException("Player isn't in this room!");
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
      throw new BadRequestException('Nothing changes!');
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
   * @param ownerId The current owner player ID. Force set owner if ignore.
   */
  async transferOwnership(
    id: RoomId,
    newOwnerId: PlayerId,
    ownerId?: PlayerId,
  ): Promise<Room> {
    const room = await this.get(id);
    if (!room) {
      throw new NotFoundException("Room doesn't exist!");
    }

    ownerId = ownerId ? ownerId : room.ownerId;
    if (room.ownerId !== ownerId) {
      throw new ForbiddenException('You are not owner of this room!');
    }

    if (ownerId === newOwnerId) {
      throw new BadRequestException('Already own this room!');
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
  async _mergeAndStoreRoom(
    room: Room,
    data: Partial<Omit<Room, 'id'>>,
  ): Promise<Room> {
    _merge(room, data);
    await this._redis.set(
      `${RedisNamespace.Room}${room.id}`,
      JSON.stringify(room),
    );

    return room;
  }
}
