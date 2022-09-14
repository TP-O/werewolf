import { Injectable } from '@nestjs/common';
import { PrismaService } from '../common/prisma.service';
import { SendPrivateMessageDto } from './dto/send-private-message.dto';

@Injectable()
export class MessageService {
  constructor(private prismaService: PrismaService) {}

  /**
   * Store new private message.
   *
   * @param senderId
   * @param privateMessageDto
   * @returns
   */
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
