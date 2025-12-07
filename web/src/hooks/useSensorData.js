import { useState, useEffect } from 'react'

export function useSensorData(
  locationSid,
  startDate,
  endDate,
  aggregation,
  types,
  refreshTrigger = 0
) {
  const [chartData, setChartData] = useState({ labels: [], datasets: [] })

  useEffect(() => {
    if (!locationSid) {
      return
    }

    let isMounted = true

    console.log('Fetching chart data for:', locationSid, startDate, endDate, aggregation, types)

    const typesParam = types.map(t => `types=${t}`).join('&')
    const aggregationParam = aggregation ? `&aggregation=${aggregation}` : ''
    const url = `api/v1/sensors/data?location_sid=${locationSid}&start_datetime=${startDate}:00Z&end_datetime=${endDate}:59Z${aggregationParam}&${typesParam}`

    fetch(url)
      .then(res => res.json())
      .then(data => {
        if (!isMounted) return

        // Group data by timestamp and type
        const grouped = {}
        data.forEach(point => {
          const timestamp = point.timestamp
          if (!grouped[timestamp]) {
            grouped[timestamp] = {}
          }
          grouped[timestamp][point.type] = point.temperature
        })

        const labels = Object.keys(grouped).map(ts =>
          new Date(ts).toLocaleString()
        )
        const datasets = []

        if (types.includes('api')) {
          datasets.push({
            label: 'API Sensor',
            data: Object.keys(grouped).map(ts => grouped[ts].api || null),
            borderColor: 'rgb(59, 130, 246)',
            backgroundColor: 'rgba(59, 130, 246, 0.5)',
            tension: 0.1
          })
        }

        if (types.includes('local')) {
          datasets.push({
            label: 'Local Sensor',
            data: Object.keys(grouped).map(ts => grouped[ts].local || null),
            borderColor: 'rgb(34, 197, 94)',
            backgroundColor: 'rgba(34, 197, 94, 0.5)',
            tension: 0.1
          })
        }

        setChartData({ labels, datasets })
      })
      .catch(err => {
        console.error('Error fetching chart data:', err)
        if (isMounted) {
          setChartData({ labels: [], datasets: [] })
        }
      })

    return () => {
      isMounted = false
    }
  }, [locationSid, startDate, endDate, aggregation, types, refreshTrigger])

  return chartData
}
