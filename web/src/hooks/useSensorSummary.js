import { useState, useEffect } from 'react'

export function useSensorSummary(locationSid, refreshTrigger = 0) {
  const [summary, setSummary] = useState(null)

  useEffect(() => {
    if (!locationSid) {
      return
    }

    let isMounted = true

    fetch(`api/v1/sensors/summary?location_sid=${locationSid}`)
      .then(res => res.json())
      .then(data => {
        if (isMounted) {
          setSummary(data)
        }
      })
      .catch(err => {
        console.error('Error fetching summary:', err)
        if (isMounted) {
          setSummary(null)
        }
      })

    return () => {
      isMounted = false
    }
  }, [locationSid, refreshTrigger])

  return summary
}
