import { NextPageWithLayout } from '@/models'
import { Box } from '@mui/system'
import { MainLayout } from '@/layouts'
import { Example } from '@/components/Example'

const Home: NextPageWithLayout = () => {
    return (
        <Box>
            <h1>Home</h1>
            <Example />
        </Box>
    )
}

Home.Layout = MainLayout

export default Home
