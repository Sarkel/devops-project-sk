import { Box, Heading, HStack, Text, IconButton } from '@chakra-ui/react'
import { IoRefresh } from 'react-icons/io5'

export default function LocationSelector({ locations, selectedLocation, onLocationChange, onRefresh }) {
  return (
    <Box as="header" mb={8} pb={6} borderBottom="2px" borderColor="gray.300">
      <Heading size="4xl" mb={6}>Temperature Checker</Heading>
      <HStack gap={3}>
        <Text fontWeight="semibold" color="gray.700">
          Location:
        </Text>
        <Box
          as="select"
          value={selectedLocation}
          onChange={(e) => onLocationChange(e.target.value)}
          px={3}
          py={2}
          borderRadius="md"
          borderWidth="1px"
          borderColor="gray.300"
          _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px blue.500' }}
        >
          {locations.map(loc => (
            <option key={loc.sid} value={loc.sid}>{loc.name}</option>
          ))}
        </Box>
        <IconButton
          aria-label="Refresh"
          size="sm"
          onClick={onRefresh}
          variant="outline"
          borderColor="gray.300"
          color="gray.700"
          _hover={{ bg: 'gray.100' }}
        >
          <IoRefresh />
        </IconButton>
      </HStack>
    </Box>
  )
}
