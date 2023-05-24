export enum CommListenEvent {
  Error = 'error',
  Success = 'success',
  FriendStatus = 'friend_status',
  PrivateMessage = 'private_message',
  RoomMessage = 'room_message',
  RoomChange = 'room_change',
}

export enum CommEmitEvent {
  RoomMessage = 'room_message',
}
