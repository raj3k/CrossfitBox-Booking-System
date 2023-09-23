import { CssVarsProvider } from "@mui/joy";
import ReactDOM from 'react-dom/client'
import './css/index.css'
import {RouterProvider} from "react-router-dom";
import router from "./routers";
import { Toaster } from 'react-hot-toast';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <CssVarsProvider>
    <RouterProvider router={router} />
    <Toaster position='top-center' toastOptions={{
      duration: 3000,
    }} />
  </CssVarsProvider>
)
