import { ListenEvent } from 'src/enum/event.enum';

type ErrorMessage = string | string[];

export type HttpErrorResponse = {
  statusCode: number;
  message: ErrorMessage;
};

export type WsErrorResponse = {
  event: ListenEvent | null;
  message: ErrorMessage;
};
