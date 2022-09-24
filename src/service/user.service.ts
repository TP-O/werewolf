import { Injectable } from '@nestjs/common';
import { User } from '@prisma/client';
import Redis from 'ioredis';
import { RedisClient } from 'src/decorator';
import { ActiveStatus, CacheNamespace } from 'src/enum';
import { RoomService } from 'src/module/chat/service/room.service';
import { Room } from 'src/module/chat/type';
import { PrismaService } from './prisma.service';

@Injectable()
export class UserService {
  @RedisClient()
  private readonly redis: Redis;

  constructor(
    private prismaService: PrismaService,
    private roomService: RoomService,
  ) {}

  /**
   * Get user by id.
   *
   * @param userId
   * @returns
   */
  async getById(userId: number) {
    const user = await this.prismaService.user.findUnique({
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
    const userId = await this.redis.get(`${CacheNamespace.SId2UId}${socketId}`);

    if (userId == null) {
      return null;
    }

    const user = await this.getById(parseInt(userId, 10));

    return user;
  }

  /**
   * Get user id by socket id.
   *
   * @param socketId
   * @returns
   */
  async getId(socketId: string) {
    const userId = await this.redis.get(`${CacheNamespace.SId2UId}${socketId}`);

    return parseInt(userId, 10);
  }

  /**
   * Get socket id list by user id.
   *
   * @param userId
   * @returns
   */
  async getSocketIds(userId: number) {
    const sIdsJSON = await this.redis.get(
      `${CacheNamespace.UID2SIds}${userId}`,
    );
    const sIds: string[] = JSON.parse(sIdsJSON || '[]');

    return sIds;
  }

  /**
   * Get socket id list by user ids.
   *
   * @param userIds
   * @returns
   */
  async getSocketIdsOfMany(userIds: number[]) {
    const keys = userIds.map((uid) => `${CacheNamespace.UID2SIds}${uid}`);
    const sIdsJSONs = await this.redis.mget(...keys);
    const sIdsOfMany: string[][] = sIdsJSONs.map((sIdsJSON) =>
      JSON.parse(sIdsJSON || '[]'),
    );

    return sIdsOfMany;
  }

  /**
   * Get user's joined room id list.
   *
   * @param userId
   * @returns
   */
  async getJoinedRoomIds(userId: number) {
    const roomIds = await this.redis.lrange(
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
  async getOnlineFriendsSocketIds(userId: number) {
    const onlineFriends = await this.prismaService.user.findMany({
      select: {
        id: true,
      },
      where: {
        OR: {
          acceptedFriends: {
            every: {
              inviterId: userId,
            },
          },
          invitedFriends: {
            every: {
              acceptorId: userId,
            },
          },
        },
        NOT: {
          statusId: null,
        },
      },
    });
    const onlineFriendsIds = onlineFriends.map((friend) => friend.id);
    const onlineFriendsSIds = await this.getSocketIdsOfMany(onlineFriendsIds);

    return onlineFriendsSIds;
  }

  /**
   * Check if two users are friends.
   *
   * @param stUserId
   * @param ndUserId
   * @returns
   */
  async areFriends(stUserId: number, ndUserId: number) {
    const relationship = await this.prismaService.friendRelationship.findFirst({
      where: {
        OR: [
          { inviterId: stUserId, acceptorId: ndUserId },
          { inviterId: ndUserId, acceptorId: stUserId },
        ],
      },
    });

    return relationship != null;
  }

  /**
   * Add new socket id to user's socket id list. Change the
   * user status to online if it's the first connetected
   * socket corresponding to the user.
   *
   * @param user
   * @param socketId conntected socket id.
   * @returns updated user.
   */
  async connect(user: User, socketId: string) {
    // Change status to online if user is offline
    if (user.statusId === null) {
      user.statusId = ActiveStatus.Online;
    }

    const sIds = await this.getSocketIds(user.id);
    sIds.push(socketId);

    await this.prismaService.user.update({
      data: user,
      where: {
        id: user.id,
      },
    });
    await this.redis
      .pipeline()
      .set(`${CacheNamespace.SId2UId}${socketId}`, user.id)
      .set(`${CacheNamespace.UID2SIds}${user.id}`, JSON.stringify(sIds))
      .exec();

    return user;
  }

  /**
   * Remove disconnected socket ids from the user record.
   * Change the user status to offline and leave all joined
   * rooms if socket id list is empty.
   *
   * @param user
   * @param socketIds disconnected socket id list.
   * @return updated user, left rooms and disconnected socket ids.
   */
  async disconnect(user: User, ...socketIds: string[]) {
    let leftRooms: Room[] = [];
    const disconnectedSIds: string[] = [];
    const redisPipe = this.redis.pipeline();
    const allSIds = await this.getSocketIds(user.id);
    const removedSocketIds = socketIds.length !== 0 ? socketIds : allSIds;

    removedSocketIds.forEach((sid) => {
      removedSocketIds.splice(removedSocketIds.indexOf(sid), 1);
      disconnectedSIds.push(sid);
      redisPipe.del(`${CacheNamespace.SId2UId}${sid}`);
    });

    // Change status to offline and leave all rooms
    // if there are no sockets is connected.
    if (allSIds.length === 0) {
      user.statusId = null;
      leftRooms = await this.roomService.leaveMany(user.id);
    } else {
      redisPipe.set(`${CacheNamespace.UID2SIds}`, JSON.stringify(allSIds));
    }

    await redisPipe.exec();
    await this.prismaService.user.update({
      data: user,
      where: {
        id: user.id,
      },
    });

    return { user, leftRooms, disconnectedSIds };
  }
}
