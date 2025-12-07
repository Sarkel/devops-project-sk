import { Box, Heading } from '@chakra-ui/react'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js'
import { Line } from 'react-chartjs-2'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
)

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top',
    },
    title: {
      display: false,
    },
  },
  scales: {
    y: {
      title: {
        display: true,
        text: 'Temperature (°C)'
      }
    }
  }
}

export default function TemperatureChart({ data }) {
  return (
    <Box bg="white" p={6} rounded="lg" shadow="md" border="1px" borderColor="gray.200">
      <Heading size="xl" color="gray.800" mb={4}>Temperature History</Heading>
      <Box h="96">
        <Line options={chartOptions} data={data} />
      </Box>
    </Box>
  )
}
