import axios from 'axios'

const commApi = axios.create({
  baseURL: `${import.meta.env.VITE_COMMUNICATION_SERVER}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
})

commApi.interceptors.request.use(
  async (config) => {
    const token = await auth.getIdToken()
    if (token)
      config.headers!.Authorization = `Bearer ${token}`

    return config
  },
)

export { commApi }
