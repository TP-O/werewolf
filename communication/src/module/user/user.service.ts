import { Injectable } from '@nestjs/common';
import Redis from 'ioredis';
import { PrismaService } from 'src/common/service/prisma.service';
import { ActiveStatus, CacheNamespace } from 'src/enum';
import { RoomService } from '../room/room.service';
import { RedisService } from 'src/common/service/redis.service';
import { Player } from '@prisma/client';
import { PlayerId } from './player.type';

@Injectable()
export class UserService {
  private readonly _redis: Redis;

  constructor(
    private prismaService: PrismaService,
    private roomService: RoomService,
    private redisService: RedisService,
  ) {
    this._redis = redisService.client;
  }

  /**
   * Get user by id.
   *
   * @param userId
   * @returns
   */
  async getById(userId: PlayerId) {
    const user = await this.prismaService.player.findUnique({
      where: { id: userId },
    });

    return user;
  }

  /**
   * Get user by socket id.
   *
   * @param socketId
   * @returns
   */
  async getBySocketId(socketId: string) {
    const userId = await this._redis.get(
      `${CacheNamespace.SId2UId}${socketId}`,
    );

    if (userId == null) {
      return null;
    }

    const user = await this.getById(userId);

    return user;
  }

  /**
   * Get user id by socket id.
   *
   * @param socketId
   * @returns
   */
  async getIdBySocketId(socketId: string) {
    const userId = await this._redis.get(
      `${CacheNamespace.SId2UId}${socketId}`,
    );

    return userId;
  }

  /**
   * Get socket id by user id.
   *
   * @param userId
   * @returns
   */
  async getSocketIdByUserId(userId: PlayerId) {
    const socketId = await this._redis.get(
      `${CacheNamespace.UID2SId}${userId}`,
    );

    return socketId;
  }

  /**
   * Get socket ids by user ids.
   *
   * @param userIds
   * @returns
   */
  async getSocketIdsByUserIds(userIds: PlayerId[]) {
    const sIdKeys = userIds.map((uid) => `${CacheNamespace.UID2SId}${uid}`);
    const sIds = await this._redis.mget(...sIdKeys);

    return sIds;
  }

  /**
   * Get user's joined room id list.
   *
   * @param userId
   * @returns
   */
  async getJoinedRoomIds(userId: number) {
    const roomIds = await this._redis.lrange(
      `${CacheNamespace.UId2RIds}${userId}`,
      0,
      -1,
    );

    return roomIds;
  }

  /**
   * Get socket id list of the user's online friends.
   *
   * @param userId
   * @returns
   */
  async getOnlineFriendsSocketIds(userId: PlayerId) {
    const onlineFriends = await this.prismaService.player.findMany({
      select: {
        id: true,
      },
      where: {
        OR: [
          {
            acceptedFriends: {
              some: {
                acceptorId: userId,
              },
            },
          },
          {
            requestedFriends: {
              some: {
                senderId: userId,
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
    const onlineFriendsSIds = await this.getSocketIdsByUserIds(
      onlineFriendsIds,
    );

    return onlineFriendsSIds;
  }

  /**
   * Get user's friend list.
   *
   * @param userId
   * @returns
   */
  async getFriendList(userId: PlayerId) {
    const friendList = await this.prismaService.player.findMany({
      where: {
        OR: {
          acceptedFriends: {
            every: {
              acceptorId: userId,
            },
          },
          requestedFriends: {
            every: {
              senderId: userId,
            },
          },
        },
      },
    });

    return friendList;
  }

  /**
   * Check if two users are friends.
   *
   * @param stUserId
   * @param ndUserId
   * @returns
   */
  async areFriends(stUserId: PlayerId, ndUserId: PlayerId) {
    const relationship = await this.prismaService.friendRelationship.findFirst({
      where: {
        OR: [
          { senderId: stUserId, acceptorId: ndUserId },
          { senderId: ndUserId, acceptorId: stUserId },
        ],
      },
    });

    return relationship != null;
  }

  /**
   * Change the user status to online and then
   * bind socket id to user id.
   *
   * @param user
   * @param socketId conntected socket id.
   * @returns updated user.
   */
  async connect(user: Player, socketId: string) {
    user.statusId = ActiveStatus.Online;

    await this.prismaService.player.update({
      data: {
        statusId: user.statusId,
      },
      where: {
        id: user.id,
      },
    });

    await this._redis
      .pipeline()
      .set(`${CacheNamespace.SId2UId}${socketId}`, user.id)
      .set(`${CacheNamespace.UID2SId}${user.id}`, socketId)
      .exec();

    return user;
  }

  /**
   * Unbind socket id from user id after that change
   * the user status to offline and leave all joined
   * rooms.
   *
   * @param user
   * @return updated user and left rooms.
   */
  async disconnect(user: Player) {
    const sId = await this.getSocketIdByUserId(user.id);
    const leftRooms = await this.roomService.leaveMany(user.id);

    user.statusId = null;

    await this._redis
      .pipeline()
      .del(`${CacheNamespace.SId2UId}${sId}`)
      .del(`${CacheNamespace.UID2SId}${user.id}`)
      .exec();

    await this.prismaService.player.update({
      data: {
        statusId: user.statusId,
      },
      where: {
        id: user.id,
      },
    });

    return { user, leftRooms, disconnectedId: sId };
  }
}
