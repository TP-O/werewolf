import { Injectable } from '@nestjs/common';
import Redis from 'ioredis';
import { Player } from '@prisma/client';
import { PrismaService, RedisService } from '../common';
import { RoomService } from '../room';
import { PlayerId } from './player.type';
import { RedisNamespace } from 'src/common/enum';
import { PlayerStatus } from './player.enum';

@Injectable()
export class PlayerService {
  private readonly _redis: Redis;

  constructor(
    private prismaService: PrismaService,
    private roomService: RoomService,
    redisService: RedisService,
  ) {
    this._redis = redisService.client;
  }

  /**
   * Get player by id.
   *
   * @param playerId
   * @returns
   */
  async getById(playerId: PlayerId) {
    const player = await this.prismaService.player.findUnique({
      where: { id: playerId },
    });

    return player;
  }

  /**
   * Get player by socket id.
   *
   * @param socketId
   * @returns
   */
  async getBySocketId(socketId: string) {
    const playerId = await this._redis.get(
      `${RedisNamespace.SId2UId}${socketId}`,
    );

    if (playerId == null) {
      return null;
    }

    const player = await this.getById(playerId);

    return player;
  }

  /**
   * Get player id by socket id.
   *
   * @param socketId
   * @returns
   */
  async getIdBySocketId(socketId: string) {
    const playerId = await this._redis.get(
      `${RedisNamespace.SId2UId}${socketId}`,
    );

    return playerId;
  }

  /**
   * Get socket id by player id.
   *
   * @param playerId
   * @returns
   */
  async getSocketIdByplayerId(playerId: PlayerId) {
    const socketId = await this._redis.get(
      `${RedisNamespace.UID2SId}${playerId}`,
    );

    return socketId;
  }

  /**
   * Get socket ids by player ids.
   *
   * @param playerIds
   * @returns
   */
  async getSocketIdsByplayerIds(playerIds: PlayerId[]) {
    const sIdKeys = playerIds.map((uid) => `${RedisNamespace.UID2SId}${uid}`);
    const sIds = await this._redis.mget(...sIdKeys);

    return sIds;
  }

  /**
   * Get player's joined room id list.
   *
   * @param playerId
   * @returns
   */
  async getJoinedRoomIds(playerId: number) {
    const roomIds = await this._redis.lrange(
      `${RedisNamespace.UId2RIds}${playerId}`,
      0,
      -1,
    );

    return roomIds;
  }

  /**
   * Get socket id list of the player's online friends.
   *
   * @param playerId
   * @returns
   */
  async getOnlineFriendsSocketIds(playerId: PlayerId) {
    const onlineFriends = await this.prismaService.player.findMany({
      select: {
        id: true,
      },
      where: {
        OR: [
          {
            acceptedFriends: {
              some: {
                acceptorId: playerId,
              },
            },
          },
          {
            requestedFriends: {
              some: {
                senderId: playerId,
              },
            },
          },
        ],
        NOT: {
          statusId: null,
        },
      },
    });

    if (onlineFriends.length === 0) {
      return [];
    }

    const onlineFriendsIds = onlineFriends.map((friend) => friend.id);
    const onlineFriendsSIds = await this.getSocketIdsByplayerIds(
      onlineFriendsIds,
    );

    return onlineFriendsSIds;
  }

  /**
   * Get player's friend list.
   *
   * @param playerId
   * @returns
   */
  async getFriendList(playerId: PlayerId) {
    const friendList = await this.prismaService.player.findMany({
      where: {
        OR: {
          acceptedFriends: {
            every: {
              acceptorId: playerId,
            },
          },
          requestedFriends: {
            every: {
              senderId: playerId,
            },
          },
        },
      },
    });

    return friendList;
  }

  /**
   * Check if two players are friends.
   *
   * @param stplayerId
   * @param ndplayerId
   * @returns
   */
  async areFriends(stplayerId: PlayerId, ndplayerId: PlayerId) {
    const relationship = await this.prismaService.friendRelationship.findFirst({
      where: {
        OR: [
          { senderId: stplayerId, acceptorId: ndplayerId },
          { senderId: ndplayerId, acceptorId: stplayerId },
        ],
      },
    });

    return relationship != null;
  }

  /**
   * Change the player status to online and then
   * bind socket id to player id.
   *
   * @param player
   * @param socketId conntected socket id.
   * @returns updated player.
   */
  async connect(player: Player, socketId: string) {
    player.statusId = PlayerStatus.Online;

    await this.prismaService.player.update({
      data: {
        statusId: player.statusId,
      },
      where: {
        id: player.id,
      },
    });

    await this._redis
      .pipeline()
      .set(`${RedisNamespace.SId2UId}${socketId}`, player.id)
      .set(`${RedisNamespace.UID2SId}${player.id}`, socketId)
      .exec();

    return player;
  }

  /**
   * Unbind socket id from player id after that change
   * the player status to offline and leave all joined
   * rooms.
   *
   * @param player
   * @return updated player and left rooms.
   */
  async disconnect(player: Player) {
    const sId = await this.getSocketIdByplayerId(player.id);
    const leftRooms = await this.roomService.leaveMany(player.id);

    player.statusId = null;

    await this._redis
      .pipeline()
      .del(`${RedisNamespace.SId2UId}${sId}`)
      .del(`${RedisNamespace.UID2SId}${player.id}`)
      .exec();

    await this.prismaService.player.update({
      data: {
        statusId: player.statusId,
      },
      where: {
        id: player.id,
      },
    });

    return { player, leftRooms, disconnectedId: sId };
  }
}
