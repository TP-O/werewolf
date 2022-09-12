import { Injectable } from '@nestjs/common';
import { PrismaService } from '../common/prisma.service';
import { SendPrivateMessageDto } from './dto/send-private-message.dto';

@Injectable()
export class MessageService {
  constructor(private prismaService: PrismaService) {}

  createPrivateMessage(
    senderId: number,
    privateMessageDto: SendPrivateMessageDto,
  ) {
    return this.prismaService.privateMessage.create({
      data: {
        ...privateMessageDto,
        senderId,
      },
    });
  }
}
