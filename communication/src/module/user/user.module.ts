import { Module } from '@nestjs/common';
import { AuthService } from 'src/common/service/auth.service';
import { PrismaService } from 'src/common/service/prisma.service';
import { RoomModule } from '../room/room.module';
import { UserController } from './user.controller';
import { UserService } from './user.service';
import { RedisService } from 'src/common/service/redis.service';
import { FirebaseService } from 'src/common/service/firebase.service';

@Module({
  imports: [RoomModule],
  controllers: [UserController],
  providers: [
    UserService,
    PrismaService,
    AuthService,
    RedisService,
    FirebaseService,
  ],
  exports: [UserService],
})
export class UserModule {
  //
}
