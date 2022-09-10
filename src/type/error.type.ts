import { ListenedEvent } from 'src/enum/event.enum';

export type WsError = {
  event: ListenedEvent | null;
  error: string | object;
};
