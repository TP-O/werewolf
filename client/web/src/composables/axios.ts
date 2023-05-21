import axios from 'axios'

const communicationClient = axios.create({
  baseURL: 'http://127.0.0.1:8079/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

communicationClient.interceptors.request.use(
  async (config) => {
    const token = await firebaseAuth.currentUser?.getIdToken()
    if (token)
      config.headers!.Authorization = `Bearer ${token}`

    return config
  },
)

export { communicationClient }
