import { Body, Controller, Delete, Post, Res, UseGuards } from '@nestjs/common';
import { FastifyReply } from 'fastify';
import {
  AddToRoomDto,
  CreatePersistentRoomsDto,
  CreateTemporaryRoomsDto,
  MuteRoomDto,
  RemoveFromRoomDto,
  RemoveRoomsDto,
} from './dto';
import { PlayerId } from '../player';
import { ChatGateway, EmitEvent, RoomEvent } from '../chat';
import { RoomService } from './room.service';
import { Room } from './room.type';
import { HmacGuard } from 'src/common/guard';

@Controller('rooms')
@UseGuards(HmacGuard)
export class RoomController {
  constructor(
    private roomService: RoomService,
    private chatGateway: ChatGateway,
  ) {}

  /**
   * Notify room joins to all members.
   *
   * @param socketIdsList
   * @param rooms
   * @param joinerIdsLst
   */
  private notifyRoomJoins(
    socketIdsList: string[][],
    rooms: Room[],
    joinerIdsLst: PlayerId[][],
  ) {
    socketIdsList.forEach((sIds, i) => {
      this.chatGateway.server.to(sIds).socketsJoin(rooms[i].id);
      this.chatGateway.server.to(sIds).emit(EmitEvent.ReceiveRoomChanges, {
        event: RoomEvent.Join,
        actorIds: joinerIdsLst[i],
        room: rooms[i],
      });
    });
  }

  /**
   * Notify room leaves to all members.
   *
   * @param socketIdsList
   * @param rooms
   * @param leaverIdsList
   */
  private notifyRoomLeaves(
    socketIdsList: string[][],
    rooms: Room[],
    leaverIdsList: PlayerId[][],
  ) {
    socketIdsList.forEach((sIds, i) => {
      this.chatGateway.server.to(sIds).socketsLeave(rooms[i].id);
      this.chatGateway.server.to(sIds).emit(EmitEvent.ReceiveRoomChanges, {
        event: RoomEvent.Leave,
        actorIds: leaverIdsList[i],
        room: rooms[i],
      });
    });
  }

  /**
   * Create many temporary rooms at once.
   *
   * @param payload
   * @param response
   */
  @Post('temporary')
  async createTemporarily(
    @Body() payload: CreateTemporaryRoomsDto,
    @Res() response: FastifyReply,
  ) {
    const { rooms, socketIdsList, joinerIdsList } =
      await this.roomService.createTemporarily(payload);

    this.notifyRoomJoins(socketIdsList, rooms, joinerIdsList);
    response.code(201).send({
      data: rooms,
    });
  }

  /**
   * Create many persistent rooms at once.
   *
   * @param payload
   * @param response
   */
  @Post('persistent')
  async createPersistently(
    @Body() payload: CreatePersistentRoomsDto,
    @Res() response: FastifyReply,
  ) {
    const { rooms, socketIdsList, joinerIdsList } =
      await this.roomService.createPersistently(payload);

    this.notifyRoomJoins(socketIdsList, rooms, joinerIdsList);
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
  async remove(@Body() payload: RemoveRoomsDto, @Res() response: FastifyReply) {
    const { rooms, socketIdsList, leaverIdsList } =
      await this.roomService.remove(payload.ids);

    this.notifyRoomLeaves(socketIdsList, rooms, leaverIdsList);
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

    this.notifyRoomJoins([socketIds], [room], [payload.memberIds]);
    response.code(200).send({
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

    this.notifyRoomLeaves([socketIds], [room], [payload.memberIds]);
    response.code(200).send({
      data: room,
    });
  }

  @Post('mute')
  async mute(@Body() payload: MuteRoomDto, @Res() response: FastifyReply) {
    const room = await this.roomService.allowChat(payload.roomId, payload.mute);

    this.chatGateway.server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      event: RoomEvent.Mute,
      actorIds: [],
      room,
    });
    response.code(201).send({
      data: room,
    });
  }
}
