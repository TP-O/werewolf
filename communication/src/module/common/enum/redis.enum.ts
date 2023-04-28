export enum RedisNamespace {
  // Map player ID to socket ID
  Id2Sid = 'id_to_sid:',

  // Map player ID to list of room ID
  Id2Rids = 'id_to_rids:',

  // Map room ID to room
  Room = 'room:',
}
