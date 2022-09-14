export enum ListenEvent {
  Connect = 'connect',
  SendPrivateMessage = 'send_private_message',
  SendGroupMessage = 'send_group_message',
  CreateRoom = 'create_room',
  JoinToRoom = 'join_to_room',
  LeaveRoom = 'leave_room',
  InviteToRoom = 'invite_to_room',
  KickOutOfRoom = 'kick_out_of_room',
}

export enum EmitEvent {
  Error = 'error',
  Success = 'success',
  UpdateFriendStatus = 'update_friend_status',
  ReceivePrivateMessage = 'receive_private_message',
  ReceiveGroupMessage = 'receive_group_message',
  ReceiveRoomChanges = 'receive_room_changes',
  ReceiveRoomInvitation = 'receive_room_invitation',
  KickedOutOfRoom = 'kicked_out_of_room',
}
