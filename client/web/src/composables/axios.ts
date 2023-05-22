import type { AxiosError } from 'axios'
import axios from 'axios'
import { error, info } from 'loglevel'

const commApi = axios.create({
  baseURL: `${import.meta.env.VITE_COMMUNICATION_SERVER}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
})

commApi.interceptors.request.use(
  async (request) => {
    const token = await auth.getIdToken()
    if (token)
      request.headers!.Authorization = `Bearer ${token}`

    info(`Communication API request [${request.url}]:`, request)
    return request
  },
)

commApi.interceptors.response.use(undefined, (err: AxiosError<any>) => {
  error(`Communication API error [${err.request.responseURL}]`, err.response)

  if (err.response?.data.message)
    return Promise.reject(new Error(err.response.data.message))

  return Promise.reject(new Error('Something went wrong'))
})

export { commApi }
