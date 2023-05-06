export enum ListenEvent {
  Connect = 'connect',
  PrivateMessage = 'private_message',
  RoomMessage = 'room_message',
}

export enum EmitEvent {
  Error = 'error',
  Success = 'success',
  FriendStatus = 'friend_status',
  PrivateMessage = 'private_message',
  RoomMessage = 'room_message',
  RoomChange = 'room_change',
}

export enum RoomChangeType {
  Join,
  Leave,
  Owner,
  Setting,
}
