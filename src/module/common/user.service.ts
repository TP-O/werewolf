import { Injectable } from '@nestjs/common';
import { User } from '@prisma/client';
import Redis from 'ioredis';
import { RedisClient } from 'src/decorator/redis.decorator';
import { CacheNamespace } from 'src/enum/cache.enum';
import { ActiveStatus } from 'src/enum/user.enum';
import { RoomService } from '../chat/room.service';
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
   * Add new socket id to user's socket id list. Change the
   * user status to online if it's the first connetected
   * socket corresponding to the user.
   *
   * @param user user record.
   * @param socketId conntected socket id.
   * @returns updated user record.
   */
  async connect(user: User, socketId: string) {
    user.sids.push(socketId);

    // Change status to online if user is offline
    if (user.statusId === null) {
      user.statusId = ActiveStatus.Online;
    }

    await this.prismaService.user.update({
      data: user,
      where: {
        id: user.id,
      },
    });
    await this.redis.set(`${CacheNamespace.SId2UId}${socketId}`, user.id);

    return user;
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
   * Remove disconnected socket ids from the user record.
   * Update the user status to offline and leave all joined
   * rooms if socket id list is empty.
   *
   * @param user user record.
   * @param socketIds disconnected socket id list.
   * @return updated user record.
   */
  async disconnect(user: User, ...socketIds: string[]) {
    const redisPipe = this.redis.pipeline();
    const removedSocketIds = socketIds.length === 0 ? user.sids : socketIds;

    removedSocketIds.forEach((sid) => {
      user.sids.splice(user.sids.indexOf(sid), 1);
      redisPipe.del(`${CacheNamespace.SId2UId}${sid}`);
    });

    // Change status to offline and leave all rooms
    // if there are no sockets is connected.
    if (removedSocketIds.length === user.sids.length) {
      const roomIds = await this.getJoinedRoomIds(user.id);
      roomIds.forEach((rid) => this.roomService.leaveRoom(rid, user.id));

      user.statusId = null;
    }

    user.sids = user.sids.filter((sid) => !socketIds.includes(sid));

    await redisPipe.exec();
    await this.prismaService.user.update({
      data: user,
      where: {
        id: user.id,
      },
    });

    return user;
  }

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
    const user = await this.prismaService.user.findUnique({
      select: {
        sids: true,
      },
      where: {
        id: userId,
      },
    });

    return user.sids;
  }

  /**
   * Get socket id list of the user's online friends.
   *
   * @param userId
   * @returns
   */
  async getOnlineFriendsSocketIds(userId: number) {
    const onlineFriends = await this.prismaService.user.findMany({
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
    const onlineFriendsSids = onlineFriends.map((friend) => friend.sids);

    return onlineFriendsSids;
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
}
