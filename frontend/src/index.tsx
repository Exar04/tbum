import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';

import { router } from './routes';
import { RouterProvider } from 'react-router-dom';
import {AuthProvider} from './context/authContext'
import WebSocketProvider, { useWebSocketSender } from './context/webSocketContext';


const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
)

root.render(
  <React.StrictMode>
    {/* <AuthProvider> */}
    <WebSocketProvider>
    <AuthProvider>
        <RouterProvider router={router}/>
    </AuthProvider>
    </WebSocketProvider>
    {/* </AuthProvider> */}
  </React.StrictMode>
)
