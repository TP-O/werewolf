import { Module } from '@nestjs/common';
import { PrismaService } from 'src/common/service/prisma.service';
import { RoomModule } from '../room/room.module';
import { UserService } from './user.service';

@Module({
  imports: [RoomModule],
  providers: [UserService, PrismaService],
  exports: [UserService],
})
export class UserModule {
  //
}
