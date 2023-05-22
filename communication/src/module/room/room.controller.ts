import {
  Body,
  Controller,
  Delete,
  HttpStatus,
  Inject,
  Post,
  Res,
  UseGuards,
  forwardRef,
} from '@nestjs/common';
import { FastifyReply } from 'fastify';
import {
  AddRoomMembersDto,
  BookRoomDto,
  ForceCreateRoomsDto,
  JoinRoomDto,
  KickOutOfRoomDto,
  LeaveRoomDto,
  MuteRoomDto,
  RemoveMemberRoomsDto,
  RemoveRoomsDto,
  TransferOwnershipDto,
} from './dto';
import { RoomService } from './room.service';
import { HttpPlayer } from 'src/common/decorator';
import { Player } from '@prisma/client';
import { HmacGuard, RequireActiveGuard, TokenGuard } from 'src/common/guard';
import { PlayerService } from '../player/player.service';
import { ChatGateway } from '../chat/chat.gateway';
import { OnlinePlayer, PlayerId, SocketId } from '../player/player.type';
import { EmitEvent, RoomChangeType } from '../chat/chat.enum';

@Controller('rooms')
export class RoomController {
  constructor(
    private readonly roomService: RoomService,
    @Inject(forwardRef(() => PlayerService))
    private readonly playerService: PlayerService,
    private readonly chatGateway: ChatGateway,
  ) {}

  /**
   * Book a room for the player.
   *
   * @param payload
   * @param player
   * @param response
   */
  @Post('/')
  @UseGuards(TokenGuard, RequireActiveGuard)
  async bookRoom(
    @Body() payload: BookRoomDto,
    @HttpPlayer() player: OnlinePlayer,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const room = await this.roomService.create({
      ownerId: player.id,
      password: payload.password,
      isMuted: false,
      memberIds: [player.id],
    });
    this.chatGateway.server.to(player.sid).socketsJoin(room.id);

    response.code(HttpStatus.CREATED).send({
      data: room,
    });
  }

  /**
   * Force create the rooms.
   *
   * @param payload
   * @param response
   */
  @Post('/many')
  @UseGuards(HmacGuard)
  async foreceCreateRooms(
    @Body() payload: ForceCreateRoomsDto,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const rooms = await this.roomService.forceCreateMany(payload.rooms);
    const playerIds: PlayerId[] = [];
    rooms.forEach((room) => {
      room.memberIds.forEach((mid) => {
        if (!playerIds.includes(mid)) {
          playerIds.push(mid);
        }
      });
    });
    const id2Sid = await this.playerService.getSocketIds(playerIds);

    rooms.forEach((room) => {
      const sids = room.memberIds
        .map((mid) => id2Sid[mid])
        .filter((sid) => !!sid) as SocketId[];

      // Clear old room if exist
      this.chatGateway.server.to(room.id).emit(EmitEvent.RoomChange, {
        changeType: RoomChangeType.Leave,
        room: {
          id: room.id,
          memberIds: room.memberIds,
        },
      });
      this.chatGateway.server.in(room.id).socketsLeave(room.id);

      // Add players to the room
      this.chatGateway.server.to(sids).socketsJoin(room.id);
      this.chatGateway.server.to(room.id).emit(EmitEvent.RoomChange, {
        changeType: RoomChangeType.Join,
        room,
      });
    });

    response.code(HttpStatus.CREATED).send({
      data: rooms,
    });
  }

  /**
   * Remove the rooms.
   *
   * @param payload
   * @param response
   */
  @Delete('/')
  @UseGuards(HmacGuard)
  async removeRooms(
    @Body() payload: RemoveRoomsDto,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const rooms = await this.roomService.removeMany(payload.ids);
    rooms.forEach((room) => {
      // Clear old room if exist
      this.chatGateway.server.to(room.id).emit(EmitEvent.RoomChange, {
        changeType: RoomChangeType.Leave,
        room: {
          id: room.id,
          memberIds: room.memberIds,
        },
      });
      this.chatGateway.server.in(room.id).socketsLeave(room.id);
    });

    response.code(HttpStatus.OK).send({
      data: rooms,
    });
  }

  /**
   * Join the room.
   *
   * @param payload
   * @param player
   * @param response
   */
  @Post('/join')
  @UseGuards(TokenGuard, RequireActiveGuard)
  async joinRoom(
    @Body() payload: JoinRoomDto,
    @HttpPlayer() player: OnlinePlayer,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const room = await this.roomService.join(
      payload.id,
      player.id,
      payload.password,
    );
    this.chatGateway.server.to(room.id).emit(EmitEvent.RoomChange, {
      changeType: RoomChangeType.Join,
      room: {
        id: room.id,
        memberIds: [player.id],
      },
    });
    this.chatGateway.server.to(player.sid).socketsJoin(room.id);

    response.code(HttpStatus.OK).send({
      data: room,
    });
  }

  /**
   * Add the members to the room.
   *
   * @param payload
   * @param response
   */
  @Post('/members')
  @UseGuards(HmacGuard)
  async addMembersToRoom(
    @Body() payload: AddRoomMembersDto,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const room = await this.roomService.forceAddMembers(
      payload.id,
      payload.memberIds,
    );
    const id2Sid = this.playerService.getSocketIds(payload.memberIds);
    const sids = payload.memberIds
      .map((mid) => id2Sid[mid])
      .filter((sid) => !!sid) as SocketId[];

    this.chatGateway.server.to(room.id).emit(EmitEvent.RoomChange, {
      changeType: RoomChangeType.Join,
      room: {
        id: room.id,
        memberIds: payload.memberIds,
      },
    });
    this.chatGateway.server.to(sids).emit(EmitEvent.RoomChange, {
      changeType: RoomChangeType.Join,
      room,
    });
    this.chatGateway.server.to(sids).socketsJoin(room.id);

    response.code(HttpStatus.OK).send({
      data: room,
    });
  }

  /**
   * Leave the room.
   *
   * @param payload
   * @param player
   * @param response
   */
  @Post('/leave')
  @UseGuards(TokenGuard, RequireActiveGuard)
  async leaveRoom(
    @Body() payload: LeaveRoomDto,
    @HttpPlayer() player: OnlinePlayer,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const room = await this.roomService.removeMembers(payload.id, [player.id]);
    this.chatGateway.server.to(player.sid).socketsLeave(room.id);
    this.chatGateway.server.to(room.id).emit(EmitEvent.RoomChange, {
      changeType: RoomChangeType.Leave,
      room: {
        id: room.id,
        memberIds: [player.id],
      },
    });

    response.code(HttpStatus.OK).send({
      data: room,
    });
  }

  /**
   * Kick the member out of the room.
   *
   * @param payload
   * @param player
   * @param response
   */
  @Post('/kick')
  @UseGuards(TokenGuard, RequireActiveGuard)
  async kickRoom(
    @Body() payload: KickOutOfRoomDto,
    @HttpPlayer() player: OnlinePlayer,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const room = await this.roomService.removeMembers(
      payload.id,
      [payload.memberId],
      player.id,
    );
    const sid = await this.playerService.getSocketId(payload.memberId);

    this.chatGateway.server.to(room.id).emit(EmitEvent.RoomChange, {
      changeType: RoomChangeType.Leave,
      room: {
        id: room.id,
        memberIds: [payload.memberId],
      },
    });

    if (sid) {
      this.chatGateway.server.to(player.sid).socketsLeave(room.id);
    }

    response.code(HttpStatus.OK).send({
      data: room,
    });
  }

  /**
   * Remove the member from the rooms.
   *
   * @param payload
   * @param response
   */
  @Delete('/members')
  @UseGuards(HmacGuard)
  async removeMemberFromRooms(
    @Body() payload: RemoveMemberRoomsDto,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const rooms = await this.roomService.removeFromRooms(
      payload.ids ?? [],
      payload.memberId,
    );
    const sid = await this.playerService.getSocketId(payload.memberId);

    rooms.forEach((room) => {
      this.chatGateway.server.to(room.id).emit(EmitEvent.RoomChange, {
        changeType: RoomChangeType.Leave,
        room: {
          id: room.id,
          memberIds: [payload.memberId],
        },
      });

      if (sid) {
        this.chatGateway.server.to(sid).socketsLeave(room.id);
      }
    });

    response.code(HttpStatus.OK).send({
      data: rooms,
    });
  }

  /**
   * Mute or unmute the room.
   *
   * @param payload
   * @param response
   */
  @Post('/mute')
  @UseGuards(HmacGuard)
  async muteRoom(
    @Body() payload: MuteRoomDto,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const room = await this.roomService.mute(payload.id, payload.isMuted);
    this.chatGateway.server.to(room.id).emit(EmitEvent.RoomChange, {
      changeType: RoomChangeType.Setting,
      room: {
        id: room.id,
        isMuted: room.isMuted,
      },
    });

    response.code(HttpStatus.OK).send({
      data: room,
    });
  }

  /**
   * Transfer room ownership.
   *
   * @param payload
   * @param player
   * @param response
   */
  @Post('/owner')
  @UseGuards(TokenGuard, RequireActiveGuard)
  async transferOwnership(
    @Body() payload: TransferOwnershipDto,
    @HttpPlayer() player: Player,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const room = await this.roomService.transferOwnership(
      payload.id,
      payload.newOwnerId,
      player.id,
    );
    this.chatGateway.server.to(room.id).emit(EmitEvent.RoomChange, {
      changeType: RoomChangeType.Owner,
      room: {
        id: room.id,
        ownerId: room.ownerId,
      },
    });

    response.code(HttpStatus.CREATED).send({
      data: room,
    });
  }
}
