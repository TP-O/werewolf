import { ListenEvent } from 'src/module/chat';

export type ErrorMessage = string | string[];

export type HttpErrorResponse = {
  statusCode: number;
  message: ErrorMessage;
};

export type WsErrorResponse = {
  event: ListenEvent;
  message: ErrorMessage;
};

export type LoggedError = Error & {
  hostType: string;
  event?: ListenEvent;
  url?: string;
  payload: any;
};
