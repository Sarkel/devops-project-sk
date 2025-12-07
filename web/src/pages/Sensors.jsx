import { useState } from 'react'
import { SimpleGrid } from '@chakra-ui/react'
import LocationSelector from '../components/LocationSelector'
import SensorCard from '../components/SensorCard'
import ChartFilters from '../components/ChartFilters'
import TemperatureChart from '../components/TemperatureChart'
import { useLocations } from '../hooks/useLocations'
import { useSensorSummary } from '../hooks/useSensorSummary'
import { useSensorData } from '../hooks/useSensorData'

export default function Sensors() {
  const { locations, selectedLocation, setSelectedLocation } = useLocations()
  const [refreshTrigger, setRefreshTrigger] = useState(0)
  const summary = useSensorSummary(selectedLocation, refreshTrigger)

  const [startDate, setStartDate] = useState('2025-09-15T00:00')
  const [endDate, setEndDate] = useState(() =>
    new Date().toISOString().slice(0, 16)
  )
  const [aggregation, setAggregation] = useState('')
  const [types, setTypes] = useState(['api', 'local'])

  const chartData = useSensorData(selectedLocation, startDate, endDate, aggregation, types, refreshTrigger)

  const toggleType = (type) => {
    setTypes(prev =>
      prev.includes(type)
        ? prev.filter(t => t !== type)
        : [...prev, type]
    )
  }

  const handleRefresh = () => {
    setRefreshTrigger(prev => prev + 1)
  }

  const handleResetFilters = () => {
    setStartDate('2025-09-15T00:00')
    setEndDate(new Date().toISOString().slice(0, 16))
    setAggregation('')
    setTypes(['api', 'local'])
  }

  return (
    <>
      <LocationSelector
        locations={locations}
        selectedLocation={selectedLocation}
        onLocationChange={setSelectedLocation}
        onRefresh={handleRefresh}
      />

      <SimpleGrid columns={{ base: 1, md: 2 }} gap={6} mb={8}>
        <SensorCard title="API Sensor" data={summary?.api} colorScheme="blue" />
        <SensorCard title="Local Sensor" data={summary?.local} colorScheme="green" />
      </SimpleGrid>

      <ChartFilters
        startDate={startDate}
        endDate={endDate}
        aggregation={aggregation}
        types={types}
        onStartDateChange={setStartDate}
        onEndDateChange={setEndDate}
        onAggregationChange={setAggregation}
        onTypeToggle={toggleType}
        onResetFilters={handleResetFilters}
      />

      <TemperatureChart data={chartData} />
    </>
  )
}
