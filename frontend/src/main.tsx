import React from 'react'
import ReactDOM from 'react-dom/client'
import './css/index.css'
import {RouterProvider} from "react-router-dom";
import router from "./routers";
import { Toaster } from 'react-hot-toast';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <Toaster position='top-center' toastOptions={{
      duration: 3000,
    }} />
  </React.StrictMode>,
)
