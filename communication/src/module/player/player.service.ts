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
  async getSocketIds(ids: PlayerId[]): Promise<SocketId[]> {
    const sids = await this._redis.mget(
      ...ids.map((id) => `${RedisNamespace.Id2Sid}${id}`),
    );
    return sids.filter((sid) => !!sid) as SocketId[];
  }

  /**
   * Get socket ID list of the player's online friends.
   *
   * @param id The player ID.
   */
  async getFriendsSocketIds(id: PlayerId): Promise<SocketId[]> {
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

    return this.getSocketIds(onlineFriends.map((friend) => friend.id));
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

  // /**
  //  * Change the player status to online and then
  //  * bind socket id to player id.
  //  *
  //  * @param player
  //  * @param socketId conntected socket id.
  //  */
  // async connect(player: Player, socketId: string) {
  //   player.statusId = PlayerStatus.Online;

  //   await this.prismaService.player.update({
  //     data: {
  //       statusId: player.statusId,
  //     },
  //     where: {
  //       id: player.id,
  //     },
  //   });

  //   await this._redis
  //     .pipeline()
  //     .set(`${RedisNamespace.SId2UId}${socketId}`, player.id)
  //     .set(`${RedisNamespace.UID2SId}${player.id}`, socketId)
  //     .exec();

  //   return player;
  // }

  // /**
  //  * Unbind socket id from player id after that change
  //  * the player status to offline and leave all joined
  //  * rooms.
  //  *
  //  * @param player
  //  * @return updated player and left rooms.
  //  */
  // async disconnect(player: Player) {
  //   const sId = await this.getSocketIdByplayerId(player.id);
  //   const leftRooms = await this.roomService.leaveMany(player.id);

  //   player.statusId = null;

  //   await this._redis
  //     .pipeline()
  //     .del(`${RedisNamespace.SId2UId}${sId}`)
  //     .del(`${RedisNamespace.UID2SId}${player.id}`)
  //     .exec();

  //   await this.prismaService.player.update({
  //     data: {
  //       statusId: player.statusId,
  //     },
  //     where: {
  //       id: player.id,
  //     },
  //   });

  //   return { player, leftRooms, disconnectedId: sId };
  // }
}
