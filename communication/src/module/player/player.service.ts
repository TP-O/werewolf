import { Injectable } from '@nestjs/common';
import Redis from 'ioredis';
import { Player } from '@prisma/client';
import { PrismaService, RedisService } from '../common';
import { PlayerId, SocketId } from './player.type';
import { PlayerStatus } from './player.enum';
import { RedisNamespace } from '../common/enum/redis.enum';

@Injectable()
export class PlayerService {
  private readonly _redis: Redis;

  constructor(
    redisService: RedisService,
    private readonly prismaService: PrismaService,
  ) {
    this._redis = redisService.client;
  }

  /**
   * Get player by ID.
   *
   * @param id The player ID.
   */
  getById(id: PlayerId): Promise<Player | null> {
    return this.prismaService.player.findUnique({
      where: { id },
    });
  }

  /**
   * Get socket ID by player ID.
   *
   * @param id The player ID.
   */
  getSocketId(id: PlayerId): Promise<SocketId | null> {
    return this._redis.get(`${RedisNamespace.Id2Sid}${id}`);
  }

  /**
   * Get socket ids by player ids.
   *
   * @param ids The list of player ID.
   */
  async getSocketIds(
    ids: PlayerId[],
  ): Promise<Record<PlayerId, SocketId | null>> {
    const sids = await this._redis.mget(
      ...ids.map((id) => `${RedisNamespace.Id2Sid}${id}`),
    );
    const id2Sid: Record<PlayerId, SocketId | null> = {};
    sids.forEach((sid, i) => {
      id2Sid[ids[i]] = sid;
    });

    return id2Sid;
  }

  /**
   * Get socket ID list of the player's online friends.
   *
   * @param id The player ID.
   */
  async getFriendsSocketIds(id: PlayerId): Promise<(SocketId | null)[]> {
    const onlineFriends = await this.prismaService.player.findMany({
      select: {
        id: true,
      },
      where: {
        OR: [
          {
            acceptedFriends: {
              some: {
                acceptorId: id,
              },
            },
          },
          {
            requestedFriends: {
              some: {
                senderId: id,
              },
            },
          },
        ],
        NOT: {
          statusId: PlayerStatus.Offline,
        },
      },
    });
    if (onlineFriends.length === 0) {
      return [];
    }

    return Object.values(
      await this.getSocketIds(onlineFriends.map((friend) => friend.id)),
    );
  }

  /**
   * Get the player's friend list.
   *
   * @param id The player ID.
   */
  getFriends(id: PlayerId): Promise<Player[]> {
    return this.prismaService.player.findMany({
      where: {
        OR: {
          acceptedFriends: {
            every: {
              acceptorId: id,
            },
          },
          requestedFriends: {
            every: {
              senderId: id,
            },
          },
        },
      },
    });
  }

  /**
   * Check if two players are friends.
   *
   * @param id1 The first player ID.
   * @param id2 The second player ID,
   */
  async areFriends(id1: PlayerId, id2: PlayerId): Promise<boolean> {
    const relationship = await this.prismaService.friendRelationship.findFirst({
      where: {
        OR: [
          { senderId: id1, acceptorId: id2 },
          { senderId: id2, acceptorId: id1 },
        ],
      },
    });
    return relationship !== null;
  }

  /**
   * Store player connection state.
   *
   * @param id The player ID.
   * @param socketId The player socket ID.
   */
  async connect(id: PlayerId, socketId: SocketId) {
    await this.prismaService.player.update({
      data: {
        statusId: PlayerStatus.Online,
      },
      where: {
        id,
      },
    });
    await this._redis.set(`${RedisNamespace.Id2Sid}${id}`, socketId);
  }

  /**
   * Remove player connection state.
   *
   * @param id The player ID.
   */
  async disconnect(id: PlayerId) {
    await this.prismaService.player.update({
      data: {
        statusId: PlayerStatus.Offline,
      },
      where: {
        id: id,
      },
    });
    await this._redis.del(`${RedisNamespace.Id2Sid}${id}`);
  }
}
