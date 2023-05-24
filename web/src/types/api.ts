export interface ResponseData<T> {
  data: T
}

export type ResponseError = ResponseData<{
  statusCode: number
  message: string | string[]
}>
