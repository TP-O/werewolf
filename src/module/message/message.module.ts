import { Module } from '@nestjs/common';
import { PrismaService } from 'src/common/service/prisma.service';
import { MessageService } from './message.service';

@Module({
  providers: [MessageService, PrismaService],
  exports: [MessageService],
})
export class MessageModule {
  //
}
