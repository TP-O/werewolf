import axios from 'axios'

const axiosClient = axios.create({
  baseURL: 'http://127.0.0.1:8079/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

axiosClient.interceptors.request.use(
  async (config) => {
    const jwt = await firebaseAuth.currentUser?.getIdToken()

    if (jwt !== undefined && jwt !== '')
      config.headers!.Authorization = `Bearer ${jwt}`

    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

axiosClient.interceptors.response.use(
  (response) => {
    return {
      ...response,
      data: response.data.data,
    }
  },
  (error) => {
    return {
      error: error.response.data.error || 'Unknown error!',
    }
  },
)

export { axiosClient }
