import { Module } from '@nestjs/common';
import { RoomService } from '../chat/room.service';
import { AuthService } from './auth.service';
import { PrismaService } from './prisma.service';
import { UserService } from './user.service';

@Module({
  exports: [AuthService, PrismaService, UserService],
  providers: [AuthService, PrismaService, UserService, RoomService],
})
export class CommonModule {}
