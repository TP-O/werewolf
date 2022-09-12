import { Module } from '@nestjs/common';
import { AuthService } from './auth.service';
import { PrismaService } from './prisma.service';
import { UserService } from './user.service';

@Module({
  exports: [AuthService, PrismaService, UserService],
  providers: [AuthService, PrismaService, UserService],
})
export class CommonModule {}
