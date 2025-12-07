import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import Sensors from './pages/Sensors'

export default function Router() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Sensors />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
