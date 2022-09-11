import { Injectable } from '@nestjs/common';
import { User } from '@prisma/client';
import Redis from 'ioredis';
import { RedisClient } from 'src/decorator/redis.decorator';
import { CacheNamespace } from 'src/enum/cache.enum';
import { PrismaService } from './prisma.service';

@Injectable()
export class UserService {
  @RedisClient()
  private readonly redis: Redis;

  constructor(private prismaService: PrismaService) {}

  async getBySocketId(socketId: string) {
    const userId = await this.redis.get(`${CacheNamespace.SId2UId}${socketId}`);

    if (userId == null) {
      return null;
    }

    const user = await this.prismaService.user.findUnique({
      where: { id: parseInt(userId, 10) },
    });

    return user;
  }

  async connect(user: User, socketId: string) {
    await this.redis
      .pipeline()
      .set(`${CacheNamespace.SId2UId}${socketId}`, user.id)
      .lpush(`${CacheNamespace.UId2SIds}${user.id}`, socketId)
      .exec();

    await this.prismaService.user.update({
      data: {
        sids: [...(user.sids as string[]), socketId],
      },
      where: {
        id: user.id,
      },
    });
  }

  async disconnect(user: User, socketId?: string) {
    const sIds =
      socketId == undefined
        ? []
        : (user.sids as string[]).filter((sid) => sid !== socketId);

    await this.redis
      .pipeline()
      .del(`${CacheNamespace.SId2UId}${socketId}`)
      .del(`${CacheNamespace.UId2SIds}${user.id}`)
      .lpush(`${CacheNamespace.UId2SIds}${user.id}`, ...sIds)
      .exec();

    await this.prismaService.user.update({
      data: {
        sids: sIds,
      },
      where: {
        id: user.id,
      },
    });
  }

  async getOnlineFriendSocketIds(userId: number) {
    const sIds: string[] = [];
    const onlineFriends = await this.prismaService.user.findMany({
      where: {
        acceptedFriends: {
          every: {
            inviterId: userId,
          },
        },
        AND: {
          invitedFriends: {
            every: {
              acceptorId: userId,
            },
          },
        },
        NOT: {
          sids: {
            equals: [],
          },
        },
      },
    });

    onlineFriends.forEach((friend) => sIds.push(...(friend.sids as string[])));

    return sIds;
  }
}
