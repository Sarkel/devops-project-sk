import { Box, Heading, Text } from '@chakra-ui/react'
import { IoTrendingUp, IoTrendingDown } from 'react-icons/io5'

function getTrendIcon(trend) {
  if (trend > 0) return <IoTrendingUp />
  if (trend < 0) return <IoTrendingDown />
}

export default function SensorCard({ title, data, colorScheme }) {
  return (
    <Box bg="white" p={6} rounded="lg" shadow="md" border="1px" borderColor="gray.200">
      <Heading size="xl" color="gray.800" mb={4}>{title}</Heading>
      {data ? (
        <>
          <Text fontSize="4xl" fontWeight="bold" color={`${colorScheme}.600`} mb={2} display="flex" alignItems="center" gap={2}>
            {data.temperature.toFixed(1)}°C
            <Box as="span" fontSize="3xl" display="inline-flex" alignItems="center">
              {getTrendIcon(data.trend)}
            </Box>
          </Text>
          <Text fontSize="sm" color="gray.600">
            {new Date(data.timestamp).toLocaleString()}
          </Text>
        </>
      ) : (
        <Text color="gray.400">No data available</Text>
      )}
    </Box>
  )
}
