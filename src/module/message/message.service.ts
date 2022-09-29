import { BadRequestException, Injectable } from '@nestjs/common';
import { PrismaService } from 'src/common/service/prisma.service';
import { UserService } from '../user/user.service';
import { SendPrivateMessageDto, SendRoomMessageDto } from './dto';

@Injectable()
export class MessageService {
  constructor(
    private prismaService: PrismaService,
    private userService: UserService,
  ) {}

  /**
   * Store new private message.
   *
   * @param senderId
   * @param privateMessageDto
   * @returns
   */
  async createPrivateMessage(
    senderId: number,
    privateMessageDto: SendPrivateMessageDto,
  ) {
    if (
      !(await this.userService.areFriends(
        senderId,
        privateMessageDto.receiverId,
      ))
    ) {
      throw new BadRequestException(
        'Only friends can send messages to each other!',
      );
    }

    return this.prismaService.privateMessage.create({
      data: {
        ...privateMessageDto,
        senderId,
      },
    });
  }

  /**
   * Store new room message.
   *
   * @param senderId
   * @param roomMessageDto
   */
  async createRoomMessage(
    senderId: number,
    roomMessageDto: SendRoomMessageDto,
  ) {
    return this.prismaService.roomMessage.create({
      data: {
        roomId: roomMessageDto.roomId,
        senderId,
        content: roomMessageDto.content,
      },
    });
  }
}
