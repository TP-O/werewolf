import { Body, Controller, Delete, Post, Res } from '@nestjs/common';
import { FastifyReply } from 'fastify';
import { EmitEvent, RoomEvent } from 'src/enum';
import { CommunicationGateway } from '../communication/communication.gateway';
import { CreateManyRoomDto, RemoveManyRoomDto } from './dto';
import { RoomService } from './room.service';

@Controller('rooms')
export class RoomController {
  constructor(
    private roomService: RoomService,
    private communicationGateway: CommunicationGateway,
  ) {}

  @Post()
  async createMany(
    @Body() payload: CreateManyRoomDto,
    @Res() response: FastifyReply,
  ) {
    const { rooms, socketIdsList } = await this.roomService.create(
      payload.rooms,
    );

    socketIdsList.forEach((sIds, i) => {
      this.communicationGateway.server.to(sIds).socketsJoin(rooms[i].id);
      this.communicationGateway.server
        .to(sIds)
        .emit(EmitEvent.ReceiveRoomChanges, {
          event: RoomEvent.Join,
          actorId: 0,
          room: rooms[i],
        });
    });

    response.code(201).send({
      data: rooms,
    });
  }

  @Delete()
  async remove(
    @Body() payload: RemoveManyRoomDto,
    @Res() response: FastifyReply,
  ) {
    const { removedRoomIds, socketIdsList } = await this.roomService.remove(
      payload.ids,
    );

    socketIdsList.forEach((sIds, i) => {
      this.communicationGateway.server.to(sIds).socketsLeave(removedRoomIds[i]);
      this.communicationGateway.server
        .to(sIds)
        .emit(EmitEvent.ReceiveRoomChanges, {
          event: RoomEvent.Remove,
          actorId: 0,
          room: {
            id: removedRoomIds[i],
          },
        });
    });

    response.code(200).send({
      data: true,
    });
  }
}
