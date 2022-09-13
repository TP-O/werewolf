import { Injectable } from '@nestjs/common';
import { User } from '@prisma/client';
import Redis from 'ioredis';
import { RedisClient } from 'src/decorator/redis.decorator';
import { CacheNamespace } from 'src/enum/cache.enum';
import { ActiveStatus } from 'src/enum/user.enum';
import { PrismaService } from './prisma.service';

@Injectable()
export class UserService {
  @RedisClient()
  private readonly redis: Redis;

  constructor(private prismaService: PrismaService) {}

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

  async disconnect(user: User, ...socketIds: string[]) {
    const redisPipe = this.redis.pipeline();
    const removedSocketIds = socketIds.length === 0 ? user.sids : socketIds;

    removedSocketIds.forEach((sid) => {
      redisPipe.del(`${CacheNamespace.SId2UId}${sid}`);
    });

    // Change status to offline if there are no sockets is connected
    if (removedSocketIds.length === user.sids.length) {
      user.statusId = null;
    }

    await redisPipe.exec();
    await this.prismaService.user.update({
      data: user,
      where: {
        id: user.id,
      },
    });
  }

  async getById(userId: number) {
    const user = await this.prismaService.user.findUnique({
      where: { id: userId },
    });

    return user;
  }

  async getBySocketId(socketId: string) {
    const userId = await this.redis.get(`${CacheNamespace.SId2UId}${socketId}`);

    if (userId == null) {
      return null;
    }

    const user = await this.getById(parseInt(userId, 10));

    return user;
  }

  async getId(socketId: string) {
    const userId = await this.redis.get(`${CacheNamespace.SId2UId}${socketId}`);

    return parseInt(userId, 10);
  }

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
