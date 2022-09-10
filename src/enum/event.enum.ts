export enum ListenedEvent {
  Connect = 'CONNECT',
  PrivateMessage = 'SEND_PRIVATE_MESSAGE',
  GroupMessage = 'SEND_GROUP_MESSAGE',
  CreateRoom = 'CREATE_ROOM',
  JoinRoom = 'JOIN_ROOM',
  LeaveRoom = 'LEAVE_ROOM',
}

export enum EmitedEvent {
  Error = 'ERROR',
  Success = 'SUCCESS',
  FriendStatus = 'UPDATE_FRIEND_STATUS',
  PrivateMessage = 'RECEIVE_PRIVATE_MESSAGE',
  GroupMessage = 'RECEIVE_GROUP_MESSAGE',
  GroupMemeber = 'UPDATE_GROUP_MEMBER',
}
