import { useState, useEffect } from 'react'

export function useLocations() {
  const [locations, setLocations] = useState([])
  const [selectedLocation, setSelectedLocation] = useState('')

  useEffect(() => {
    let isMounted = true

    fetch('/api/v1/locations')
      .then(res => res.json())
      .then(data => {
        if (isMounted) {
          setLocations(data || [])
          if (data && data.length > 0) {
            setSelectedLocation(data[0].sid)
          }
        }
      })
      .catch(err => {
        console.error('Error fetching locations:', err)
        if (isMounted) {
          setLocations([])
        }
      })

    return () => {
      isMounted = false
    }
  }, [])

  return { locations, selectedLocation, setSelectedLocation }
}
