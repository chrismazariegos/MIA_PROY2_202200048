import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import Consola from './Componentes/Consola'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <h1>Proyecto 2 - 202200048</h1>
      <Consola />
    </>
  )
}

export default App
