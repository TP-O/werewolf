import { Body, Controller, Delete, Post, Res } from '@nestjs/common';
import { FastifyReply } from 'fastify';
import { EmitEvent, RoomEvent } from 'src/enum';
import { CommunicationGateway } from '../communication/communication.gateway';
import {
  AddToRoomDto,
  CreateManyRoomDto,
  RemoveFromRoomDto,
  RemoveManyRoomDto,
} from './dto';
import { RoomService } from './room.service';

@Controller('rooms')
export class RoomController {
  constructor(
    private roomService: RoomService,
    private communicationGateway: CommunicationGateway,
  ) {}

  /**
   * Create many rooms at once.
   *
   * @param payload
   * @param response
   */
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

  /**
   * Remove many room at once.
   *
   * @param payload
   * @param response
   */
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

  /**
   * Add many members to room.
   *
   * @param payload
   * @param response
   */
  @Post('members')
  async addMembers(
    @Body() payload: AddToRoomDto,
    @Res() response: FastifyReply,
  ) {
    const { room, socketIds } = await this.roomService.addMembers(
      payload.roomId,
      payload.memberIds,
    );

    this.communicationGateway.server.to(socketIds).socketsJoin(room.id);
    this.communicationGateway.server
      .to(socketIds)
      .emit(EmitEvent.ReceiveRoomChanges, {
        event: RoomEvent.Join,
        actorId: 0,
        room,
      });

    response.code(201).send({
      data: room,
    });
  }

  /**
   * Remove many members from room.
   *
   * @param payload
   * @param response
   */
  @Delete('members')
  async removeMembers(
    @Body() payload: RemoveFromRoomDto,
    @Res() response: FastifyReply,
  ) {
    const { room, socketIds } = await this.roomService.removeMembers(
      payload.roomId,
      payload.memberIds,
    );

    this.communicationGateway.server.to(socketIds).socketsLeave(room.id);
    this.communicationGateway.server
      .to(socketIds)
      .emit(EmitEvent.ReceiveRoomChanges, {
        event: RoomEvent.Join,
        actorId: 0,
        room,
      });

    response.code(201).send({
      data: room,
    });
  }
}
