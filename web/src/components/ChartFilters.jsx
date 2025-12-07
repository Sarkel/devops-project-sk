import { Box, Heading, SimpleGrid, Input, HStack, Text, IconButton } from '@chakra-ui/react'
import { IoReload } from 'react-icons/io5'

export default function ChartFilters({
  startDate,
  endDate,
  aggregation,
  types,
  onStartDateChange,
  onEndDateChange,
  onAggregationChange,
  onTypeToggle,
  onResetFilters
}) {
  return (
    <Box bg="white" p={6} rounded="lg" shadow="md" border="1px" borderColor="gray.200" mb={6}>
      <HStack justifyContent="space-between" mb={4}>
        <Heading size="xl" color="gray.800">Chart Filters</Heading>
        <IconButton
          aria-label="Reset filters"
          size="sm"
          onClick={onResetFilters}
          variant="outline"
          borderColor="gray.300"
          color="gray.700"
          _hover={{ bg: 'gray.100' }}
        >
          <IoReload />
        </IconButton>
      </HStack>
      <SimpleGrid columns={{ base: 1, md: 2, lg: 4 }} gap={4}>
        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>Start Date:</Text>
          <Input
            type="datetime-local"
            value={startDate}
            onChange={(e) => onStartDateChange(e.target.value)}
          />
        </Box>
        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>End Date:</Text>
          <Input
            type="datetime-local"
            value={endDate}
            onChange={(e) => onEndDateChange(e.target.value)}
          />
        </Box>
        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>Aggregation:</Text>
          <Box
            as="select"
            value={aggregation}
            onChange={(e) => onAggregationChange(e.target.value)}
            px={3}
            py={2}
            borderRadius="md"
            borderWidth="1px"
            borderColor="gray.300"
            w="full"
            _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px blue.500' }}
          >
            <option value="">----</option>
            <option value="day">Day</option>
          </Box>
        </Box>
        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>Sensor Types:</Text>
          <HStack gap={4} h={10} alignItems="center">
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
              <input
                type="checkbox"
                checked={types.includes('api')}
                onChange={() => onTypeToggle('api')}
                style={{ width: '16px', height: '16px' }}
              />
              <Text fontSize="sm" color="gray.700">API</Text>
            </label>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
              <input
                type="checkbox"
                checked={types.includes('local')}
                onChange={() => onTypeToggle('local')}
                style={{ width: '16px', height: '16px' }}
              />
              <Text fontSize="sm" color="gray.700">Local</Text>
            </label>
          </HStack>
        </Box>
      </SimpleGrid>
    </Box>
  )
}
